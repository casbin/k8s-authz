package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	"github.com/casbin/k8s-authz/internal/handler"
)

func tlsHandler(c *gin.Context) {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     "localhost:8080",
	})
	err := secureMiddleware.Process(c.Writer, c.Request)
	// If there was an error, do not continue.
	if err != nil {
		return
	}
	c.Next()
}

func main() {
	r := gin.Default()
	r.Any("/", handler.Handler)
	r.Use(tlsHandler)
	r.RunTLS(":8080", "config/certificate/server.crt", "config/certificate/server.key")
}
