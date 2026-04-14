package application

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	billingdomain "korp_backend/internal/modules/billing/domain"
	stockdomain "korp_backend/internal/modules/stock/domain"
	platformai "korp_backend/internal/platform/ai"
)

type InvoiceReader interface {
	ListByOwner(ctx context.Context, ownerID string) ([]billingdomain.Invoice, error)
}

type AIInsightsGenerator interface {
	GenerateOperationalInsights(ctx context.Context, prompt string) (platformai.GeneratedOperationalInsights, string, error)
}

type GenerateAIInsightsUseCase struct {
	productRepository stockdomain.ProductRepository
	invoiceReader     InvoiceReader
	generator         AIInsightsGenerator
	lowStockThreshold int
}

type AIInsightsOutput struct {
	GeneratedAt        time.Time                      `json:"generated_at"`
	Model              string                         `json:"model"`
	Overview           string                         `json:"overview"`
	Alerts             []string                       `json:"alerts"`
	Actions            []string                       `json:"actions"`
	BillingNotes       []string                       `json:"billing_notes"`
	BuyRecommendations []platformai.BuyRecommendation `json:"buy_recommendations"`
	SearchQueries      []string                       `json:"search_queries"`
	Sources            []platformai.GroundingSource   `json:"sources"`
	ProductCount       int                            `json:"product_count"`
	InvoiceCount       int                            `json:"invoice_count"`
	OpenInvoiceCount   int                            `json:"open_invoice_count"`
	LowStockCount      int                            `json:"low_stock_count"`
	OutOfStockCount    int                            `json:"out_of_stock_count"`
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
	log.Printf("ai insights usecase: loading products owner_id=%s", ownerID)
	products, err := uc.productRepository.ListByOwner(ctx, ownerID)
	if err != nil {
		log.Printf("ai insights usecase: failed loading products owner_id=%s err=%v", ownerID, err)
		return AIInsightsOutput{}, err
	}

	log.Printf("ai insights usecase: loading invoices owner_id=%s", ownerID)
	invoices, err := uc.invoiceReader.ListByOwner(ctx, ownerID)
	if err != nil {
		log.Printf("ai insights usecase: failed loading invoices owner_id=%s err=%v", ownerID, err)
		return AIInsightsOutput{}, err
	}

	snapshot := buildInsightSnapshot(products, invoices, uc.lowStockThreshold)
	log.Printf(
		"ai insights usecase: snapshot owner_id=%s products=%d invoices=%d open_invoices=%d low_stock=%d out_of_stock=%d prompt_chars=%d",
		ownerID,
		len(products),
		len(invoices),
		snapshot.openInvoices,
		snapshot.lowStockCount,
		snapshot.outOfStockCount,
		len(snapshot.prompt()),
	)
	insights, model, err := uc.generator.GenerateOperationalInsights(ctx, snapshot.prompt())
	if err != nil {
		log.Printf("ai insights usecase: generator failed owner_id=%s err=%v", ownerID, err)
		return AIInsightsOutput{}, err
	}

	log.Printf(
		"ai insights usecase: generator succeeded owner_id=%s model=%s alerts=%d actions=%d billing_notes=%d recommendations=%d sources=%d",
		ownerID,
		model,
		len(insights.Alerts),
		len(insights.Actions),
		len(insights.BillingNotes),
		len(insights.BuyRecommendations),
		len(insights.Sources),
	)

	return AIInsightsOutput{
		GeneratedAt:        time.Now().UTC(),
		Model:              model,
		Overview:           insights.Overview,
		Alerts:             insights.Alerts,
		Actions:            insights.Actions,
		BillingNotes:       insights.BillingNotes,
		BuyRecommendations: insights.BuyRecommendations,
		SearchQueries:      insights.SearchQueries,
		Sources:            insights.Sources,
		ProductCount:       len(products),
		InvoiceCount:       len(invoices),
		OpenInvoiceCount:   snapshot.openInvoices,
		LowStockCount:      snapshot.lowStockCount,
		OutOfStockCount:    snapshot.outOfStockCount,
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
		`Analise os dados operacionais abaixo e responda apenas com JSON valido, sem markdown, sem comentarios e sem texto antes ou depois.

		Formato esperado:
		{
		"overview": "string",
		"alerts": ["string"],
		"actions": ["string"],
		"billing_notes": ["string"],
		"buy_recommendations": [
			{
			"name": "string",
			"category": "string",
			"reason": "string",
			"market_signal": "string",
			"stock_relation": "string"
			}
		]
		}

		Regras:
		- Seja extremamente objetivo.
		- Nao use quebras de linha dentro de strings JSON.
		- Cada item de lista deve ter no maximo 12 palavras.
		- "overview" deve ter no maximo 1 frase curta.
		- Em "alerts", retorne no maximo 3 itens.
		- Em "actions", retorne no maximo 3 itens.
		- Em "billing_notes", retorne no maximo 2 itens.
		- Em "buy_recommendations", retorne exatamente 5 itens.
		- Em cada recomendacao, "reason", "market_signal" e "stock_relation" devem ter no maximo 14 palavras cada.
		- Nao invente dados que nao estao no contexto interno.
		- Para buy_recommendations, pesquise na web com foco no mercado brasileiro atual e recomende exatamente 5 produtos em alta.
		- Cada recomendacao deve considerar o que ja existe no estoque e explicar a relacao com o estoque atual.
		- Se alguma observacao depender de interpretacao, deixe isso claro no texto.
		- Seja objetivo e util para uma pequena operacao comercial.
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
