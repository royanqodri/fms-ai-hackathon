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

	"github.com/royanqodri/fms-ai-hackathon/controller"

	repoMst "github.com/royanqodri/fms-ai-hackathon/repository/master"
	repoOp "github.com/royanqodri/fms-ai-hackathon/repository/operation"

	serviceMst "github.com/royanqodri/fms-ai-hackathon/service/master"
	serviceOp "github.com/royanqodri/fms-ai-hackathon/service/operation"
)

var (
	// REPOSITORY
	mstCustomerRepo  repoMst.MstCustomerRepository
	mstUserRepo      repoMst.MstUserRepository
	mstTruckRepo     repoMst.MstTruckRepository
	trxTelemetryRepo repoOp.TrxTelemetryRepository
	logAlertRepo     repoOp.LogAlertRepository
	aiRepo           repoOp.AiRepository

	// SERVICE
	mstCustomerService  serviceMst.MstCustomerService
	mstUserService      serviceMst.MstUserService
	mstTruckService     serviceMst.MstTruckService
	trxTelemetryService serviceOp.TrxTelemetryService
	logAlertService     serviceOp.LogAlertService
	chatService         serviceOp.ChatService

	// CONTROLLER
	mstCustomerController  controller.MstCustomerController
	mstUserController      controller.MstUserController
	mstTruckController     controller.MstTruckController
	trxTelemetryController controller.TrxTelemetryController
	logAlertController     controller.LogAlertController
	chatController         controller.ChatController
)

func main() {
	// init config
	initConfig()

	// init location
	util.InitLocation()

	// init database
	initDatabase()

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
}

func initRepository() {
	mstCustomerRepo = repoMst.NewMstCustomerRepository()
	mstUserRepo = repoMst.NewMstUserRepository()
	mstTruckRepo = repoMst.NewMstTruckRepository()
	trxTelemetryRepo = repoOp.NewTrxTelemetryRepository()
	logAlertRepo = repoOp.NewLogAlertRepository()
	aiRepo = repoOp.NewAiRepository()
}

func initService() {
	mstCustomerService = serviceMst.NewMstCustomerService(mstCustomerRepo)
	mstUserService = serviceMst.NewMstUserService(mstUserRepo)
	mstTruckService = serviceMst.NewMstTruckService(mstTruckRepo)
	trxTelemetryService = serviceOp.NewTrxTelemetryService(trxTelemetryRepo, logAlertRepo)
	logAlertService = serviceOp.NewLogAlertService(logAlertRepo)
	chatService = serviceOp.NewChatService(aiRepo)
}

func initController() {
	mstCustomerController = controller.NewMstCustomerController(mstCustomerService)
	mstUserController = controller.NewMstUserController(mstUserService)
	mstTruckController = controller.NewMstTruckController(mstTruckService)
	trxTelemetryController = controller.NewTrxTelemetryController(trxTelemetryService)
	logAlertController = controller.NewLogAlertController(logAlertService)
	chatController = controller.NewChatController(chatService, aiRepo)
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
	router.Use(gin.Recovery(), util.CORSMiddleware())

	router.Static("/public", "./public")

	//ENDPOINT
	api := router.Group("/api")
	{
		// Master
		api.GET("/customers", mstCustomerController.GetByParams)
		api.POST("/customers", mstCustomerController.Post)

		api.GET("/users", mstUserController.GetByParams)
		api.POST("/users", mstUserController.Post)

		api.GET("/trucks", mstTruckController.GetByParams)
		api.POST("/trucks", mstTruckController.Post)

		// Operation
		api.GET("/telemetries", trxTelemetryController.GetByParams)
		api.POST("/telemetries", trxTelemetryController.Post)

		api.GET("/alerts", logAlertController.GetByParams)
		api.POST("/alerts", logAlertController.Post)

		api.POST("/chat", chatController.Chat)
		api.POST("/telegram/webhook", chatController.TelegramWebhook)

	}

	// Endpoint for Health Check / Test
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
