package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"splitFiles/controller"
)

func StartServer(controller controller.Controller) {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers := Handlers{controller: controller}

	router.GET("/savePage", handlers.GetSavePage)
	router.POST("/saveFile", handlers.SaveFile)
	router.GET("/loadFile/:filename", handlers.LoadFile)

	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:8080/savePage")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
