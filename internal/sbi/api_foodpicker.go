package sbi

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	foodList = []string{"Beef Noodle", "Burger", "Bento", "McDonalds", "Luway"}
	mu       sync.Mutex // To ensure thread-safe operations on the food list
)

func (s *Server) getFoodPickerRoutes() []Route {
	return []Route{
		{
			Name:    "FoodPicker GET",
			Method:  http.MethodGet,
			Pattern: "",
			APIFunc: func(c *gin.Context) {
				mu.Lock()
				defer mu.Unlock()

				rand.Seed(time.Now().UnixNano())
				randomFood := foodList[rand.Intn(len(foodList))]
				c.JSON(http.StatusOK, gin.H{"lunch/dinner pick": randomFood})
			},
			// curl -X GET http://127.0.0.1:8000/foodpicker -w "\n"
		},
		{
			Name:    "FoodPicker POST",
			Method:  http.MethodPost,
			Pattern: "",
			APIFunc: func(c *gin.Context) {
				var newFood struct {
					Name string `json:"name"`
				}
				if err := c.ShouldBindJSON(&newFood); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				mu.Lock()
				defer mu.Unlock()

				foodList = append(foodList, newFood.Name)
				c.JSON(http.StatusOK, gin.H{"message": "Food added successfully", "foodList": foodList})
			},
			// curl -X POST http://127.0.0.1:8000/foodpicker -H "Content-Type: application/json" -d '{"name":"McDonalds"}' -w "\n"
		},
	}
}

// foodpicker.go
func (s *Server) GetFoodPickerRoutes() []Route {
	return s.getFoodPickerRoutes()
}
