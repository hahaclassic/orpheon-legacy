package auth_ctrl_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	auth_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/auth"
// 	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
// 	"github.com/hahaclassic/orpheon/backend/mocks"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestRegister(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	mockService := mocks.NewAuthService(t)
// 	controller := auth_ctrl.NewAuthController(mockService)

// 	router := gin.Default()
// 	controller.RegisterRoutes(router.Group("/auth"))

// 	creds := entity.UserCredentials{
// 		Login:    "test@example.com",
// 		Password: "securepass123",
// 	}
// 	expectedTokens := &entity.AuthTokens{
// 		Access:  "access-token",
// 		Refresh: "refresh-token",
// 	}

// 	mockService.EXPECT().
// 		RegisterUser(mock.Anything, &creds).
// 		Return(expectedTokens, nil)

// 	body, _ := json.Marshal(creds)
// 	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var resp map[string]string
// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedTokens.Access, resp["access_token"])
// }
