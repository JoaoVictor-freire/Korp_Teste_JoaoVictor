package bootstrap

import (
	"log"

	"korp_backend/internal/modules/billing/application"
	billingdomain "korp_backend/internal/modules/billing/domain"
	postgresrepo "korp_backend/internal/modules/billing/infra/postgres"
	billinghttp "korp_backend/internal/modules/billing/interfaces/http"
	stockdomain "korp_backend/internal/modules/stock/domain"
	stockpostgres "korp_backend/internal/modules/stock/infra/postgres"
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
	server := httpx.NewServer(httpConfig.Address())

	var invoiceRepository billingdomain.InvoiceRepository
	var productRepository stockdomain.ProductRepository

	db, err := database.OpenFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db, postgresrepo.Models()...); err != nil {
		log.Fatal(err)
	}

	invoiceRepository = postgresrepo.NewInvoiceRepository(db)
	productRepository = stockpostgres.NewProductRepository(db)

	handler := billinghttp.NewHandler(
		application.NewCreateInvoiceUseCase(invoiceRepository, productRepository),
		application.NewListInvoicesUseCase(invoiceRepository),
		application.NewGetInvoiceUseCase(invoiceRepository),
		application.NewUpdateInvoiceUseCase(invoiceRepository),
		application.NewDeleteInvoiceUseCase(invoiceRepository),
		application.NewCloseInvoiceUseCase(invoiceRepository, productRepository),
	)

	signer := auth.NewTokenSignerFromEnv()
	authMiddleware := auth.RequireAuth(signer)

	billinghttp.RegisterRoutes(server.Engine(), handler, authMiddleware)

	return App{
		Server: server,
	}
}
