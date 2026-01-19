package rest

import (
	"log/slog"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/adapter/transport/rest/handler"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	router *gin.Engine
	uc     *usecase.MonitorInteractor
	logger *slog.Logger
}

func NewGinServer(uc *usecase.MonitorInteractor) *GinServer {
	s := &GinServer{
		router: gin.New(),
		uc:     uc,
		logger: slog.With(slog.String("component", "gin_server")),
	}
	s.router.Use(gin.Recovery())

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
	}
}

func (s *GinServer) GetRouter() *gin.Engine {
	return s.router
}
