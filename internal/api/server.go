package api

import (
	"DIA/internal/app/handler"
	"DIA/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


var lines = []string{"first line", "second line", "third line", "fourth line"}


func StartServer(){
	log.Println("server start up")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}
	handler := handler.NewHandler(repo)
	
	r := gin.Default()
	
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	
	r.GET("/materials", handler.GetMaterials)
	r.GET("/material/:id", handler.GetMaterial)
	r.GET("/material/cart/:id", handler.GetCart)

	r.Run()
	log.Println("server down")
}