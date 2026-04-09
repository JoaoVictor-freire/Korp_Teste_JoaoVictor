package bootstrap

import (
	"log"
	"os"
	"strings"

	"korp_backend/internal/modules/billing/application"
	billingdomain "korp_backend/internal/modules/billing/domain"
	memoryrepo "korp_backend/internal/modules/billing/infra/memory"
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
	server := httpx.NewServer(httpConfig.Address())

	// Choose Postgres if DATABASE_URL is set, otherwise fallback to in-memory repositories.
	var invoiceRepository billingdomain.InvoiceRepository

	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn != "" {
		db, err := database.Open(dsn)
		if err != nil {
			log.Fatal(err)
		}

		if err := database.Migrate(db, postgresrepo.Models()...); err != nil {
			log.Fatal(err)
		}

		invoiceRepository = postgresrepo.NewInvoiceRepository(db)
	} else {
		invoiceRepository = memoryrepo.NewInvoiceRepository()
	}

	handler := billinghttp.NewHandler(
		application.NewCreateInvoiceUseCase(invoiceRepository),
		application.NewListInvoicesUseCase(invoiceRepository),
	)

	signer := auth.NewTokenSignerFromEnv()
	authMiddleware := auth.RequireAuth(signer)

	billinghttp.RegisterRoutes(server.Engine(), handler, authMiddleware)

	return App{
		Server: server,
	}
}
