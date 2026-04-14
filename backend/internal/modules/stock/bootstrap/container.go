package bootstrap

import (
	"log"

	billingpostgres "korp_backend/internal/modules/billing/infra/postgres"
	"korp_backend/internal/modules/stock/application"
	stockdomain "korp_backend/internal/modules/stock/domain"
	postgresrepo "korp_backend/internal/modules/stock/infra/postgres"
	stockhttp "korp_backend/internal/modules/stock/interfaces/http"
	usersapp "korp_backend/internal/modules/users/application"
	usersdomain "korp_backend/internal/modules/users/domain"
	userspostgres "korp_backend/internal/modules/users/infra/postgres"
	usershttp "korp_backend/internal/modules/users/interfaces/http"
	platformai "korp_backend/internal/platform/ai"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/config"
	"korp_backend/internal/platform/database"
	"korp_backend/internal/platform/httpx"
)

type App struct {
	Server *httpx.Server
}

func New() App {
	httpConfig := config.NewHTTPConfig("STOCK_SERVICE", "8081")
	server := httpx.NewServer(httpConfig.Address())

	var stockRepository stockdomain.ProductRepository
	var userRepository usersdomain.UserRepository

	db, err := database.OpenFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	models := make([]any, 0, 8)
	models = append(models, userspostgres.Models()...)
	models = append(models, postgresrepo.Models()...)
	models = append(models, billingpostgres.Models()...)
	if err := database.Migrate(db, models...); err != nil {
		log.Fatal(err)
	}

	stockRepository = postgresrepo.NewProductRepository(db)
	userRepository = userspostgres.NewUserRepository(db)
	invoiceRepository := billingpostgres.NewInvoiceRepository(db)
	geminiClient := platformai.NewGeminiClientFromEnv()
	aiInsights := application.NewGenerateAIInsightsUseCase(
		stockRepository,
		invoiceRepository,
		geminiClient,
		config.GetEnvAsInt("AI_LOW_STOCK_THRESHOLD", 5),
	)

	handler := stockhttp.NewHandler(
		application.NewCreateProductUseCase(stockRepository),
		application.NewListProductsUseCase(stockRepository),
		application.NewGetProductUseCase(stockRepository),
		application.NewUpdateProductUseCase(stockRepository),
		application.NewDeleteProductUseCase(stockRepository),
		application.NewDecreaseStockUseCase(stockRepository),
		aiInsights,
	)

	signer := auth.NewTokenSignerFromEnv()
	authMiddleware := auth.RequireAuth(signer)

	stockhttp.RegisterRoutes(server.Engine(), handler, authMiddleware)

	// User/Auth endpoints live in this service for the skeleton.
	usersHandler := usershttp.NewHandler(
		usersapp.NewRegisterUseCase(userRepository),
		usersapp.NewLoginUseCase(userRepository),
		signer,
	)
	usershttp.RegisterRoutes(server.Engine(), usersHandler)

	return App{
		Server: server,
	}
}
