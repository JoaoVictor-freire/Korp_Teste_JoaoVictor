package application

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	billingdomain "korp_backend/internal/modules/billing/domain"
	stockdomain "korp_backend/internal/modules/stock/domain"
)

type InvoiceReader interface {
	ListByOwner(ctx context.Context, ownerID string) ([]billingdomain.Invoice, error)
}

type AIInsightsGenerator interface {
	GenerateOperationalInsights(ctx context.Context, prompt string) (string, string, error)
}

type GenerateAIInsightsUseCase struct {
	productRepository stockdomain.ProductRepository
	invoiceReader     InvoiceReader
	generator         AIInsightsGenerator
	lowStockThreshold int
}

type AIInsightsOutput struct {
	GeneratedAt      time.Time `json:"generated_at"`
	Model            string    `json:"model"`
	Content          string    `json:"content"`
	ProductCount     int       `json:"product_count"`
	InvoiceCount     int       `json:"invoice_count"`
	OpenInvoiceCount int       `json:"open_invoice_count"`
	LowStockCount    int       `json:"low_stock_count"`
	OutOfStockCount  int       `json:"out_of_stock_count"`
}

func NewGenerateAIInsightsUseCase(
	productRepository stockdomain.ProductRepository,
	invoiceReader InvoiceReader,
	generator AIInsightsGenerator,
	lowStockThreshold int,
) GenerateAIInsightsUseCase {
	return GenerateAIInsightsUseCase{
		productRepository: productRepository,
		invoiceReader:     invoiceReader,
		generator:         generator,
		lowStockThreshold: lowStockThreshold,
	}
}

func (uc GenerateAIInsightsUseCase) Execute(ctx context.Context, ownerID string) (AIInsightsOutput, error) {
	products, err := uc.productRepository.ListByOwner(ctx, ownerID)
	if err != nil {
		return AIInsightsOutput{}, err
	}

	invoices, err := uc.invoiceReader.ListByOwner(ctx, ownerID)
	if err != nil {
		return AIInsightsOutput{}, err
	}

	snapshot := buildInsightSnapshot(products, invoices, uc.lowStockThreshold)
	content, model, err := uc.generator.GenerateOperationalInsights(ctx, snapshot.prompt())
	if err != nil {
		return AIInsightsOutput{}, err
	}

	return AIInsightsOutput{
		GeneratedAt:      time.Now().UTC(),
		Model:            model,
		Content:          content,
		ProductCount:     len(products),
		InvoiceCount:     len(invoices),
		OpenInvoiceCount: snapshot.openInvoices,
		LowStockCount:    snapshot.lowStockCount,
		OutOfStockCount:  snapshot.outOfStockCount,
	}, nil
}

type insightSnapshot struct {
	productLines      []string
	invoiceLines      []string
	totalStockUnits   int
	openInvoices      int
	closedInvoices    int
	lowStockCount     int
	outOfStockCount   int
	totalInvoiceItems int
	lowStockThreshold int
}

func buildInsightSnapshot(products []stockdomain.Product, invoices []billingdomain.Invoice, lowStockThreshold int) insightSnapshot {
	snapshot := insightSnapshot{
		productLines:      make([]string, 0, len(products)),
		invoiceLines:      make([]string, 0, len(invoices)),
		lowStockThreshold: lowStockThreshold,
	}

	sortedProducts := append([]stockdomain.Product(nil), products...)
	sort.Slice(sortedProducts, func(i, j int) bool {
		return sortedProducts[i].Code < sortedProducts[j].Code
	})

	for _, product := range sortedProducts {
		snapshot.totalStockUnits += product.Stock
		if product.Stock == 0 {
			snapshot.outOfStockCount++
		}
		if product.Stock <= lowStockThreshold {
			snapshot.lowStockCount++
		}

		snapshot.productLines = append(snapshot.productLines, fmt.Sprintf(
			"- codigo=%s | descricao=%s | estoque=%d | criado_em=%s",
			product.Code,
			sanitizePromptText(product.Description),
			product.Stock,
			formatTimestamp(product.CreatedAt),
		))
	}

	sortedInvoices := append([]billingdomain.Invoice(nil), invoices...)
	sort.Slice(sortedInvoices, func(i, j int) bool {
		return sortedInvoices[i].Number < sortedInvoices[j].Number
	})

	for _, invoice := range sortedInvoices {
		if invoice.Status == billingdomain.StatusOpen {
			snapshot.openInvoices++
		} else {
			snapshot.closedInvoices++
		}

		itemDescriptions := make([]string, 0, len(invoice.Items))
		for _, item := range invoice.Items {
			snapshot.totalInvoiceItems += item.Quantity
			itemDescriptions = append(itemDescriptions, fmt.Sprintf("%s x%d", item.ProductCode, item.Quantity))
		}

		snapshot.invoiceLines = append(snapshot.invoiceLines, fmt.Sprintf(
			"- numero=%d | status=%s | itens=%s | criada_em=%s",
			invoice.Number,
			invoice.Status,
			strings.Join(itemDescriptions, ", "),
			formatTimestamp(invoice.CreatedAt),
		))
	}

	return snapshot
}

func (s insightSnapshot) prompt() string {
	productSection := "Nenhum produto cadastrado."
	if len(s.productLines) > 0 {
		productSection = strings.Join(s.productLines, "\n")
	}

	invoiceSection := "Nenhuma nota fiscal encontrada."
	if len(s.invoiceLines) > 0 {
		invoiceSection = strings.Join(s.invoiceLines, "\n")
	}

	return fmt.Sprintf(
		`Analise os dados operacionais abaixo e gere um texto curto com estes blocos:
1. Resumo geral
2. Alertas imediatos
3. Oportunidades de acao
4. Observacoes sobre faturamento

Regras:
- Nao invente dados que nao estao no contexto.
- Se algo for inferencia, diga explicitamente "inferencia".
- Seja objetivo, pratico e use bullets curtos.
- Considere estoque baixo quando o estoque for menor ou igual a %d.

METRICAS:
- total_produtos=%d
- total_unidades_em_estoque=%d
- total_notas=%d
- notas_abertas=%d
- notas_fechadas=%d
- produtos_com_estoque_baixo=%d
- produtos_sem_estoque=%d
- total_unidades_faturadas_nas_notas=%d

PRODUTOS:
%s

NOTAS_FISCAIS:
%s`,
		s.lowStockThreshold,
		len(s.productLines),
		s.totalStockUnits,
		len(s.invoiceLines),
		s.openInvoices,
		s.closedInvoices,
		s.lowStockCount,
		s.outOfStockCount,
		s.totalInvoiceItems,
		productSection,
		invoiceSection,
	)
}

func sanitizePromptText(value string) string {
	replacer := strings.NewReplacer("\n", " ", "\r", " ", "\t", " ")
	return replacer.Replace(strings.TrimSpace(value))
}

func formatTimestamp(value time.Time) string {
	if value.IsZero() {
		return "nao informado"
	}

	return value.Format(time.RFC3339)
}
