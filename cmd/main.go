package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/fms-ai-hackathon/config"
	"github.com/royanqodri/fms-ai-hackathon/database"
	"github.com/royanqodri/fms-ai-hackathon/util"
	"github.com/royanqodri/fms-ai-hackathon/util/logging"
)

var (
// repository

// service

// controller

)

func main() {
	// init config
	initConfig()

	// init location
	util.InitLocation()

	// init database
	initDatabase()

	// init log
	logging.InitLog()
	defer logging.CloseLog()

	// init repository
	initRepository()

	// init service
	initService()

	// init controller
	initController()

	// init router
	initRouter()
}

func initConfig() {
	err := config.Init()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}
}

func initDatabase() {
	database.ConnectDatabaseMaster()
	database.ConnectDatabaseOperation()
}

func initRepository() {
	// truckRepo = repository.NewTruckRepository()
}

func initService() {
	// truckService = service.NewTruckService(truckRepo)
}

func initController() {
	// truckController = controller.NewTruckController(truckService)
}

func initRouter() {
	// set service mode
	if strings.ToLower(config.Get().ServiceMode) == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if strings.ToLower(config.Get().ServiceMode) == "test" {
		gin.SetMode(gin.TestMode)
	} else if strings.ToLower(config.Get().ServiceMode) == "release" || strings.ToLower(config.Get().ServiceMode) == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// init router
	router := gin.New()
	router.SetTrustedProxies([]string{"0.0.0.0/0"})

	// public route - require no authentication
	router.Use(gin.Recovery(), util.CORSMiddleware())
	// Endpoint
	// router.GET("/api/trucks", truckController.GetAll)

	// Endpoint untuk Health Check / Test
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "FMS AI Backend Server OK 🚀"})
	})

	// graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	// init timeout
	readTimeout := time.Duration(config.Get().Service.ReadTimeout) * time.Second
	writeTimeout := time.Duration(config.Get().Service.WriteTimeout) * time.Second
	idleTimeout := time.Duration(config.Get().Service.IdleTimeout) * time.Second

	// init server
	server := &http.Server{
		Addr:         config.Get().Service.HttpAddress,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	log.Println("🚀 Service is Running on", config.Get().Service.HttpAddress)

	// listen server
	err := server.ListenAndServe()
	util.PanicIfError(err)
}
