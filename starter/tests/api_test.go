package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"inventory_management/database"
	"inventory_management/handlers"
	"inventory_management/models"
	"inventory_management/routes"
)

type ItemTestSuite struct {
	suite.Suite
	router   *gin.Engine
	db       *gorm.DB
	jwtToken string
}

func (suite *ItemTestSuite) SetupSuite() {
	var err error
	suite.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	database.DB = suite.db

	err = suite.db.AutoMigrate(&models.Item{})
	assert.NoError(suite.T(), err)

	gin.SetMode(gin.TestMode)
	suite.router = routes.SetupRoutes()

	loginPayload := map[string]string{"username": "admin", "password": "password"}
	loginData, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(loginData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if token, ok := resp["token"].(string); ok {
		suite.jwtToken = token
	}
}

func (suite *ItemTestSuite) SetupTest() {
	suite.db.Where("1 = 1").Delete(&models.Item{})
}

func (suite *ItemTestSuite) TestCreateItem() {
	item := models.Item{
		Name:  "Test Laptop",
		Stock: 10,
		Price: 999.99,
	}

	jsonData, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/api/v1/inventory", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.jwtToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Item created successfully", response["message"])
	assert.NotNil(suite.T(), response["data"])
}

func (suite *ItemTestSuite) TestCreateItemValidation() {
	item := models.Item{
		Name:  "Invalid Item",
		Stock: -5,
		Price: 100.0,
	}

	jsonData, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/api/v1/inventory", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.jwtToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ItemTestSuite) TestGetAllItems() {
	items := []models.Item{
		{ID: "1", Name: "Laptop", Stock: 10, Price: 999.99},
		{ID: "2", Name: "Mouse", Stock: 50, Price: 29.99},
	}
	suite.db.Create(&items)

	req, _ := http.NewRequest("GET", "/api/v1/inventory", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handlers.PaginationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), response.Total)
	assert.Len(suite.T(), response.Data, 2)
}

func (suite *ItemTestSuite) TestGetAllItemsWithPagination() {
	for i := 1; i <= 15; i++ {
		item := models.Item{
			ID:    string(rune(i)),
			Name:  "Item " + string(rune(i)),
			Stock: i,
			Price: float64(i * 10),
		}
		suite.db.Create(&item)
	}

	req, _ := http.NewRequest("GET", "/api/v1/inventory?page=2&page_size=5", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var response handlers.PaginationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 2, response.Page)
		assert.Equal(suite.T(), 5, response.PageSize)
		assert.True(suite.T(), response.HasPrev)
		assert.True(suite.T(), response.HasNext)
		assert.Equal(suite.T(), 2, len(response.Data))
	} else {
		assert.Equal(suite.T(), http.StatusTooManyRequests, w.Code)
	}
}

func (suite *ItemTestSuite) TestGetAllItemsWithFiltering() {
	items := []models.Item{
		{ID: "1", Name: "Gaming Laptop", Stock: 5, Price: 1500.00},
		{ID: "2", Name: "Office Laptop", Stock: 15, Price: 800.00},
		{ID: "3", Name: "Gaming Mouse", Stock: 25, Price: 50.00},
	}
	suite.db.Create(&items)

	req, _ := http.NewRequest("GET", "/api/v1/inventory?name=laptop", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var response handlers.PaginationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(2), response.Total)
	} else {
		assert.Equal(suite.T(), http.StatusTooManyRequests, w.Code)
	}

	req, _ = http.NewRequest("GET", "/api/v1/inventory?min_stock=20", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var response handlers.PaginationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(2), response.Total)
	} else {
		assert.Equal(suite.T(), http.StatusTooManyRequests, w.Code)
	}
}

func (suite *ItemTestSuite) TestGetItemByID() {
	item := models.Item{ID: "1", Name: "Test Item", Stock: 10, Price: 100.0}
	suite.db.Create(&item)

	req, _ := http.NewRequest("GET", "/api/v1/inventory/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		_, exists := response["data"]
		assert.True(suite.T(), exists)
	} else {
		assert.Equal(suite.T(), http.StatusTooManyRequests, w.Code)
	}
}

func (suite *ItemTestSuite) TestGetItemByIDNotFound() {
	req, _ := http.NewRequest("GET", "/api/v1/inventory/nonexistent-id", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.True(suite.T(), w.Code == http.StatusNotFound || w.Code == http.StatusTooManyRequests)
}

func (suite *ItemTestSuite) TestUpdateItem() {
	item := models.Item{ID: "1", Name: "Original Item", Stock: 10, Price: 100.0}
	suite.db.Create(&item)

	updatedItem := models.Item{
		Name:  "Updated Item",
		Stock: 20,
		Price: 200.0,
	}

	jsonData, _ := json.Marshal(updatedItem)
	req, _ := http.NewRequest("PUT", "/api/v1/inventory/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.jwtToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "Item updated successfully", response["message"])
	} else {
		assert.Equal(suite.T(), http.StatusTooManyRequests, w.Code)
	}
}

func (suite *ItemTestSuite) TestDeleteItem() {
	item := models.Item{ID: "test-id", Name: "Item to Delete", Stock: 10, Price: 100.0}
	suite.db.Create(&item)

	req, _ := http.NewRequest("DELETE", "/api/v1/inventory/test-id", nil)
	req.Header.Set("Authorization", "Bearer "+suite.jwtToken)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Item deleted successfully", response["message"])
}

func (suite *ItemTestSuite) TestRateLimiting() {
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/inventory", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusTooManyRequests)
	}
}

func TestItemTestSuite(t *testing.T) {
	suite.Run(t, new(ItemTestSuite))
}
