package rest

import (
	"log/slog"
	"time"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/adapter/transport/rest/handler"
	"github.com/gin-contrib/cors"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	router *gin.Engine
	uc     *usecase.MonitorInteractor
	logger *slog.Logger
}

func NewGinServer(uc *usecase.MonitorInteractor) *GinServer {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))

	router.Use(gin.Recovery())

	router.Use(gin.Logger())

	s := &GinServer{
		router: router,
		uc:     uc,
		logger: slog.With(slog.String("component", "gin_server")),
	}

	s.mapRoutes()
	return s
}

func (s *GinServer) Run(port string) error {
	return s.router.Run(port)
}

func (s *GinServer) mapRoutes() {
	h := handler.NewHandler(s.uc)
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/getUserSites/:userID", h.GetUserSites)
		v1.POST("/addSite", h.AddSite)
		v1.DELETE("/deleteSite", h.DeleteSite)
	}
}

func (s *GinServer) GetRouter() *gin.Engine {
	return s.router
}
