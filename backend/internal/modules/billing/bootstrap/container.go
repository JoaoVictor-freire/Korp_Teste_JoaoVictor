package bootstrap

import (
	"fmt"
	"log"
	"time"

	"korp_backend/internal/modules/billing/application"
	billingdomain "korp_backend/internal/modules/billing/domain"
	stockclient "korp_backend/internal/modules/billing/infra/http"
	postgresrepo "korp_backend/internal/modules/billing/infra/postgres"
	billinghttp "korp_backend/internal/modules/billing/interfaces/http"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/config"
	"korp_backend/internal/platform/database"
	"korp_backend/internal/platform/httpx"
)

type App struct {
	Server *httpx.Server
}

func New() App {
	httpConfig := config.NewHTTPConfig("BILLING_SERVICE", "8082")
	stockConfig := config.NewHTTPConfig("STOCK_SERVICE", "8081")
	server := httpx.NewServer(httpConfig.Address())

	var invoiceRepository billingdomain.InvoiceRepository
	stockService := stockclient.NewStockClient(
		fmt.Sprintf("http://%s", stockConfig.Address()),
		stockclient.StockClientConfig{
			RequestTimeout:   time.Duration(config.GetEnvAsInt("BILLING_STOCK_REQUEST_TIMEOUT_MS", 5000)) * time.Millisecond,
			FailureThreshold: config.GetEnvAsInt("BILLING_STOCK_CIRCUIT_FAILURE_THRESHOLD", 3),
			ResetTimeout:     time.Duration(config.GetEnvAsInt("BILLING_STOCK_CIRCUIT_RESET_TIMEOUT_MS", 15000)) * time.Millisecond,
		},
	)

	db, err := database.OpenFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db, postgresrepo.Models()...); err != nil {
		log.Fatal(err)
	}

	invoiceRepository = postgresrepo.NewInvoiceRepository(db)

	handler := billinghttp.NewHandler(
		application.NewCreateInvoiceUseCase(invoiceRepository, stockService),
		application.NewListInvoicesUseCase(invoiceRepository),
		application.NewGetInvoiceUseCase(invoiceRepository),
		application.NewUpdateInvoiceUseCase(invoiceRepository),
		application.NewDeleteInvoiceUseCase(invoiceRepository),
		application.NewCloseInvoiceUseCase(invoiceRepository, stockService),
	)

	signer := auth.NewTokenSignerFromEnv()
	authMiddleware := auth.RequireAuth(signer)

	billinghttp.RegisterRoutes(server.Engine(), handler, authMiddleware)

	return App{
		Server: server,
	}
}
