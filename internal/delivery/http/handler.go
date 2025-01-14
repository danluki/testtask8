package http

import (
	"net/http"

	"github.com/danluki/test-task-8/internal/config"
	v1 "github.com/danluki/test-task-8/internal/delivery/http/v1"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	gorm *gorm.DB
}

func NewHandler(gorm *gorm.DB) *Handler {
	return &Handler{
		gorm: gorm,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.gorm)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}

func corsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
