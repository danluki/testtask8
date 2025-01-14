package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/danluki/test-task-8/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HandlerTestSuite struct {
	suite.Suite
	db      *gorm.DB
	handler *Handler
	router  *gin.Engine
}

func (suite *HandlerTestSuite) SetupSuite() {
	// Initialize PostgreSQL database connection
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err)

	// Migrate the User schema
	err = db.AutoMigrate(&store.User{})
	suite.Require().NoError(err)

	// Create handler and router
	handler := &Handler{gorm: db}
	router := gin.Default()
	api := router.Group("/api")
	handler.initUsersRoutes(api)

	// Save components to the suite
	suite.db = db
	suite.handler = handler
	suite.router = router
}

func (suite *HandlerTestSuite) TearDownSuite() {
	// Drop the users table after the tests
	suite.db.Exec("DROP TABLE IF EXISTS users")
}

// Test addUser
func (suite *HandlerTestSuite) TestAddUser() {
	newUser := addUserInput{Name: "John Doe", Email: "john.doe@example.com"}
	body, _ := json.Marshal(newUser)

	req := httptest.NewRequest(http.MethodPost, "/api/users/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response store.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), newUser.Name, response.Name)
	assert.Equal(suite.T(), newUser.Email, response.Email)
}

// Test getUsers
func (suite *HandlerTestSuite) TestGetUsers() {
	// Add dummy users
	suite.db.Create(&store.User{Name: "User 1", Email: "user1@example.com"})
	suite.db.Create(&store.User{Name: "User 2", Email: "user2@example.com"})

	req := httptest.NewRequest(http.MethodGet, "/api/users/?page=1&limit=1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		Users []store.User `json:"users"`
		Meta  struct {
			Total    int `json:"total"`
			Page     int `json:"page"`
			PageSize int `json:"pageSize"`
		} `json:"meta"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	assert.Len(suite.T(), response.Users, 1)
	assert.Equal(suite.T(), 3, response.Meta.Total)
	assert.Equal(suite.T(), 1, response.Meta.Page)
	assert.Equal(suite.T(), 1, response.Meta.PageSize)
}

// Test updateUser
func (suite *HandlerTestSuite) TestUpdateUser() {
	user := store.User{Name: "John Doe", Email: "john.doe@example.com"}
	suite.db.Create(&user)

	update := map[string]string{"name": "Updated Name"}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/users/"+strconv.Itoa(int(user.ID)),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response store.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), "Updated Name", response.Name)
	assert.Equal(suite.T(), user.Email, response.Email)
}

// Test deleteUser
func (suite *HandlerTestSuite) TestDeleteUser() {
	user := store.User{Name: "John Doe", Email: "john.doe@example.com"}
	suite.db.Create(&user)

	req := httptest.NewRequest(http.MethodDelete, "/api/users/"+strconv.Itoa(int(user.ID)), nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), "user deleted", response["message"])

	// Ensure the user no longer exists
	var count int64
	suite.db.Model(&store.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

// Run the test suite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
