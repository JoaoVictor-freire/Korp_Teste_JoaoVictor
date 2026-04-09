package bootstrap

import (
	"korp_backend/internal/modules/stock/application"
	"korp_backend/internal/modules/stock/infra/memory"
	stockhttp "korp_backend/internal/modules/stock/interfaces/http"
	usersapp "korp_backend/internal/modules/users/application"
	usersmemory "korp_backend/internal/modules/users/infra/memory"
	usershttp "korp_backend/internal/modules/users/interfaces/http"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/config"
	"korp_backend/internal/platform/httpx"
)

type App struct {
	Server *httpx.Server
}

func New() App {
	httpConfig := config.NewHTTPConfig("STOCK_SERVICE", "8081")
	server := httpx.NewServer(httpConfig.Address())

	productRepository := memory.NewProductRepository()

	handler := stockhttp.NewHandler(
		application.NewCreateProductUseCase(productRepository),
		application.NewListProductsUseCase(productRepository),
	)

	signer := auth.NewTokenSignerFromEnv()
	authMiddleware := auth.RequireAuth(signer)

	stockhttp.RegisterRoutes(server.Engine(), handler, authMiddleware)

	// User/Auth endpoints live in this service for the skeleton.
	userRepo := usersmemory.NewUserRepository()
	usersHandler := usershttp.NewHandler(
		usersapp.NewRegisterUseCase(userRepo),
		usersapp.NewLoginUseCase(userRepo),
		signer,
	)
	usershttp.RegisterRoutes(server.Engine(), usersHandler)

	return App{
		Server: server,
	}
}
