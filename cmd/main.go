package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

type ServiceImplementation struct{}

type Service interface {
	DaysLeft() int64
}

func (i ServiceImplementation) DaysLeft() int64 {
	days := time.Date(2025, time.January, 0, 0, 0, 0, 0, time.UTC)
	dur := time.Until(days)
	result := int64(dur.Hours() / 24)

	return result
}

type Endpoint struct {
	service Service
}

func NewEndpoint(svc Service) *Endpoint {
	return &Endpoint{
		service: svc,
	}
}

func setupRouter(e *Endpoint) *gin.Engine {
	router := gin.Default()
	router.GET("/start", e.StatusHandler, e.StatusMiddleware)

	return router
}

func (e *Endpoint) StatusHandler(ctx *gin.Context) {
	days := e.service.DaysLeft()
	message := fmt.Sprintf("Number of days: %d", days)

	ctx.String(http.StatusOK, message)
}

func (e *Endpoint) StatusMiddleware(ctx *gin.Context) {
	role := ctx.Request.Header.Get("User-Role")
	result := strings.ToLower(role)

	if strings.Contains("admin", result) {
		fmt.Println("Red button user detected")
	}

	ctx.Next()
}

func main() {
	router := setupRouter(NewEndpoint(ServiceImplementation{}))
	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
