package v1

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/danluki/test-task-8/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("/", h.addUser)
		users.GET("/", h.getUsers)
		users.PUT("/:id", h.updateUser)
		users.DELETE("/:id", h.deleteUser)
	}
}

type addUserInput struct {
	Name  string `json:"name"  binding:"required,min=2,max=64"`
	Email string `json:"email" binding:"required,email,max=64"`
}

func (h *Handler) addUser(c *gin.Context) {
	var inp addUserInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newUser store.User
	err := h.gorm.Where("email = ?", inp.Email).Find(&newUser).Limit(1).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	if newUser.ID != 0 {
		newResponse(c, http.StatusBadRequest, "user with this email already exists")
		return
	}

	err = h.gorm.Create(&store.User{Name: inp.Name, Email: inp.Email}).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	err = h.gorm.First(&newUser, "email = ?", inp.Email).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (h *Handler) getUsers(c *gin.Context) {
	var users []store.User
	var total int64

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	offset := (pageNum - 1) * limitNum

	h.gorm.Limit(limitNum).Offset(offset).Find(&users)
	h.gorm.Model(&store.User{}).Count(&total)

	// Ответ с метаинформацией
	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"meta": gin.H{
			"total":    total,
			"page":     pageNum,
			"pageSize": limitNum,
		},
	})
}

func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id")
	var user store.User

	if err := h.gorm.First(&user, id).Error; err != nil {
		newResponse(c, http.StatusNotFound, "user not found")
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.gorm.Save(&user).Error; err != nil {
		newResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) deleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.gorm.First(&store.User{}, id).Error; err != nil {
		newResponse(c, http.StatusNotFound, "user not found")
		return
	}

	err := h.gorm.Delete(&store.User{}, id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted",
	})
}

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	slog.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
