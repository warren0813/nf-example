package sbi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func setupFoodPickerServer(t *testing.T) *sbi.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000,
			},
		},
	}).AnyTimes()

	return sbi.NewServer(nfApp, "")
}

func Test_FoodPickerGET(t *testing.T) {
	server := setupFoodPickerServer(t)

	httpRecorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(httpRecorder)

	req, err := http.NewRequest(http.MethodGet, "/foodpicker", nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}
	ginCtx.Request = req

	// Call the GET handler directly
	for _, route := range server.GetFoodPickerRoutes() {
		if route.Method == http.MethodGet {
			route.APIFunc(ginCtx)
		}
	}

	if httpRecorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", httpRecorder.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(httpRecorder.Body.Bytes(), &body); err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if _, ok := body["lunch/dinner pick"]; !ok {
		t.Errorf("Expected key 'lunch/dinner pick' in response body")
	}
}

func Test_FoodPickerPOST(t *testing.T) {
	server := setupFoodPickerServer(t)

	httpRecorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(httpRecorder)

	payload := map[string]string{"name": "Sushi"}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "/foodpicker", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	ginCtx.Request = req

	// Call the POST handler directly
	for _, route := range server.GetFoodPickerRoutes() {
		if route.Method == http.MethodPost {
			route.APIFunc(ginCtx)
		}
	}

	if httpRecorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", httpRecorder.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(httpRecorder.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if resp["message"] != "Food added successfully" {
		t.Errorf("Expected success message, got %v", resp["message"])
	}

	// Optional: Check if "Sushi" is present in returned list
	list, ok := resp["foodList"].([]interface{})
	if !ok {
		t.Errorf("Expected foodList array in response")
	} else {
		found := false
		for _, f := range list {
			if f == "Sushi" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected 'Sushi' to be in foodList")
		}
	}
}
