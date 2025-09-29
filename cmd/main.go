package main

import (
	"DIA/config"
	"DIA/internal/app/dsn"
	"DIA/internal/app/handler"
	"DIA/internal/app/repository"
	"DIA/internal/pkg"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.NewRepository(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}