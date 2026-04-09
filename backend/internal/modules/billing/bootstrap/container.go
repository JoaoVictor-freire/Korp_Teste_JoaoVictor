package bootstrap

import (
	"korp_backend/internal/modules/billing/application"
	"korp_backend/internal/modules/billing/infra/memory"
	billinghttp "korp_backend/internal/modules/billing/interfaces/http"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/config"
	"korp_backend/internal/platform/httpx"
)

type App struct {
	Server *httpx.Server
}

func New() App {
	httpConfig := config.NewHTTPConfig("BILLING_SERVICE", "8082")
	server := httpx.NewServer(httpConfig.Address())

	invoiceRepository := memory.NewInvoiceRepository()

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
