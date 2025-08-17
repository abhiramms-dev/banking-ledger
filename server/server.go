package server

import (
	"banking-ledger/adapter/storage/database"
	"banking-ledger/config"
	"banking-ledger/consts"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HTTPServer interface {
	Run() error
	ShutDown() error
}

type GinServer struct {
	Cfg    *config.EnvConfig
	Router *gin.Engine
}

func New() *GinServer {
	cfg, err := config.Load(consts.AppName)
	if err != nil {
		panic(err)
	}

	// here initializing the router
	router := initRouter()
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	psqlConfig := database.NewPsqlConnectivity(&cfg.PsqlDB)
	psqlConn, err := psqlConfig.Connection()
	if err != nil {
		log.Printf("unable to connect the postgresql database : %v", err)
		return nil
	}
	mongoConfig := database.NewMongoConnectivity(&cfg.MongoDB)
	mongoConn, err := mongoConfig.Connection()
	if err != nil {
		log.Printf("unable to connect the mongo database : %v", err)
		return nil
	}
	api := router.Group("/api/:version")

	fmt.Println(api, psqlConn, mongoConn)

	return &GinServer{
		Cfg:    cfg,
		Router: router,
	}
}

func initRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return router
}

func (gs GinServer) Run() error {
	log.Print("Starting gin server in...", gs.Cfg.Port)
	return gs.launch()
}

func (gs GinServer) ShutDown() error {
	return nil
}

func (gs GinServer) launch() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", gs.Cfg.Port),
		Handler:           gs.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe()
}
