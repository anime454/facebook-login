package main

import (
	"log"
	"net/http"

	"github.com/anime454/facebook-login/handler"
	"github.com/anime454/facebook-login/service"
	"github.com/gin-gonic/gin"
)

func main() {

	facebookService := service.NewFacebookService()
	facebookHandler := handler.NewFacebookHandler(facebookService)

	server := gin.Default()
	server.GET("/facebook/login_page", facebookHandler.LoginPage())
	server.GET("/facebook/login", facebookHandler.Login())
	server.GET("/facebook/login/callback", facebookHandler.Callback())

	srv := &http.Server{
		Addr:    ":" + "10003",
		Handler: server,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
