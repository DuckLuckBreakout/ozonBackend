package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/DuckLuckBreakout/ozonBackend/configs"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/metrics"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/grpc_utils"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
	cart_repo "github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart/repository"
	cart_usecase "github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart/usecase"
	proto "github.com/DuckLuckBreakout/ozonBackend/services/cart/proto/cart"
	cart_server "github.com/DuckLuckBreakout/ozonBackend/services/cart/server"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func InitCartService() {
	// Load session_service environment
	err := godotenv.Load(configs.CartServiceMainEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Load session service redis environment
	err = godotenv.Load(configs.CartServiceRedisEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Load network environment
	err = godotenv.Load(configs.NetworkEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Init logger
	mainLogger := logger.Logger{}
	err = mainLogger.InitLogger(configs.CartServiceLog, os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitCartService()

	// Connect to redis db
	redisConn := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
	if redisConn == nil {
		log.Fatal(errors.ErrDBFailedConnection.Error())
	}
	defer redisConn.Close()

	cartRepo := cart_repo.NewSessionRedisRepository(redisConn)
	cartUCase := cart_usecase.NewUseCase(cartRepo)
	cartServer := cart_server.NewCartServer(cartUCase)

	lis, err := net.Listen(
		os.Getenv("CART_SERVICE_PROTOCOL"),
		fmt.Sprintf("%s:%s",
			os.Getenv("CART_SERVICE_HOST"),
			os.Getenv("CART_SERVICE_PORT")),
	)
	if err != nil {
		log.Fatalf("error start session service %v", err)
	}

	metric, err := metrics.CreateNewMetrics("cart_service")
	if err != nil {
		log.Fatal(err)
	}
	accessInterceptor := grpc_utils.AccessInterceptor(metric)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(accessInterceptor),
	)
	proto.RegisterCartServiceServer(server, cartServer)

	go metrics.CreateNewMetricsRouter(os.Getenv("CART_SERVICE_HOST"))

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
