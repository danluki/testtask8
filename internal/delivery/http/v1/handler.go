package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	// There is better to use usecases here to handle business logic and provide better testability
	gorm *gorm.DB
}

func NewHandler(gorm *gorm.DB) *Handler {
	return &Handler{
		gorm: gorm,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
	}
}
