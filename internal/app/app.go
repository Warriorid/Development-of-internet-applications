package app

import (
	"DIA/internal/handler"
	"DIA/internal/repository"
	"DIA/internal/service"
	"DIA/pkg"
	"DIA/pkg/config"
	"DIA/pkg/minio"
	redisClient "DIA/pkg/redis"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
	Redis   *redis.Client
}

func NewApp() *Application {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresCofigString := config.PostgresConfDB()
	db, err := pkg.NewPostgresDB(postgresCofigString)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	minioClient, err := minio.NewMinioClient()
    if err != nil {
        logrus.Fatalf("failed to create minio client: %v", err)
    }

	redisClient, err := redisClient.NewClient()
	if err != nil {
		logrus.Fatalf("failed to create redis client: %v", err)
	}
	logrus.Info("Redis client successfully connected")

	repository := repository.NewRepository(db)
	service := service.NewService(repository, minioClient)
	handler := handler.NewHandler(service, redisClient)
	router := gin.Default()

	return &Application{
		Config:  conf,
		Router:  router,
		Handler: handler,
		Redis:   redisClient,
	}
}

func (a *Application) RunApp() {

	a.Handler.RegisterHandler(a.Router)
	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	server := &http.Server{
		Addr:    serverAddress,
		Handler: a.Router,
	}

	go func() {
		logrus.Println("Service started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown error: %v", err)
	}
	logrus.Info("Server down")
}