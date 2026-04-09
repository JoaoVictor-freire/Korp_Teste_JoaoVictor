package bootstrap

import (
	"log"
	"os"
	"strings"

	"korp_backend/internal/modules/stock/application"
	stockdomain "korp_backend/internal/modules/stock/domain"
	memoryrepo "korp_backend/internal/modules/stock/infra/memory"
	postgresrepo "korp_backend/internal/modules/stock/infra/postgres"
	stockhttp "korp_backend/internal/modules/stock/interfaces/http"
	usersapp "korp_backend/internal/modules/users/application"
	usersdomain "korp_backend/internal/modules/users/domain"
	usersmemory "korp_backend/internal/modules/users/infra/memory"
	userspostgres "korp_backend/internal/modules/users/infra/postgres"
	usershttp "korp_backend/internal/modules/users/interfaces/http"
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

	// Choose Postgres if DATABASE_URL is set, otherwise fallback to in-memory repositories.
	var stockRepository stockdomain.ProductRepository
	var userRepository usersdomain.UserRepository

	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn != "" {
		db, err := database.Open(dsn)
		if err != nil {
			log.Fatal(err)
		}

		models := make([]any, 0, 8)
		models = append(models, userspostgres.Models()...)
		models = append(models, postgresrepo.Models()...)
		if err := database.Migrate(db, models...); err != nil {
			log.Fatal(err)
		}

		stockRepository = postgresrepo.NewProductRepository(db)
		userRepository = userspostgres.NewUserRepository(db)
	} else {
		stockRepository = memoryrepo.NewProductRepository()
		userRepository = usersmemory.NewUserRepository()
	}

	handler := stockhttp.NewHandler(
		application.NewCreateProductUseCase(stockRepository),
		application.NewListProductsUseCase(stockRepository),
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
