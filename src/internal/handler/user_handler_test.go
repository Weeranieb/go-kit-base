package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/weeranieb/go-kit-base/src/internal/model"
)

// Test CreateUser handler
func (s *HandlerTestSuite) TestCreateUser_Success() {
	createReq := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResponse := &model.UserResponse{
		ID:        1,
		Username:  createReq.Username,
		Email:     createReq.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userService.On("CreateUser", createReq).Return(expectedResponse, nil)

	app := fiber.New()
	app.Post("/users", s.userHandler.CreateUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusCreated, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestCreateUser_InvalidBody() {
	app := fiber.New()
	app.Post("/users", s.userHandler.CreateUser)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestCreateUser_ValidationError() {
	req := &model.CreateUserRequest{
		Username: "ab", // Too short
		Email:    "invalid-email",
		Password: "123", // Too short
	}

	app := fiber.New()
	app.Post("/users", s.userHandler.CreateUser)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestCreateUser_ServiceError() {
	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	s.userService.On("CreateUser", req).Return(nil, errors.New("email already exists"))

	app := fiber.New()
	app.Post("/users", s.userHandler.CreateUser)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusConflict, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test GetUser handler
func (s *HandlerTestSuite) TestGetUser_Success() {
	userID := uint(1)
	expectedResponse := &model.UserResponse{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userService.On("GetUser", userID).Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/users/:id", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/users/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestGetUser_InvalidID() {
	app := fiber.New()
	app.Get("/users/:id", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/users/invalid", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestGetUser_NotFound() {
	userID := uint(999)
	s.userService.On("GetUser", userID).Return(nil, errors.New("user not found"))

	app := fiber.New()
	app.Get("/users/:id", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/users/999", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusNotFound, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test UpdateUser handler
func (s *HandlerTestSuite) TestUpdateUser_Success() {
	userID := uint(1)
	req := &model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	expectedResponse := &model.UserResponse{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userService.On("UpdateUser", userID, req).Return(expectedResponse, nil)

	app := fiber.New()
	app.Put("/users/:id", s.userHandler.UpdateUser)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestUpdateUser_InvalidID() {
	app := fiber.New()
	app.Put("/users/:id", s.userHandler.UpdateUser)

	req := httptest.NewRequest("PUT", "/users/invalid", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestUpdateUser_InvalidBody() {
	app := fiber.New()
	app.Put("/users/:id", s.userHandler.UpdateUser)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestUpdateUser_ValidationError() {
	req := &model.UpdateUserRequest{
		Email: "invalid-email",
	}

	app := fiber.New()
	app.Put("/users/:id", s.userHandler.UpdateUser)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestUpdateUser_ServiceError() {
	userID := uint(1)
	req := &model.UpdateUserRequest{
		Username: "updateduser",
	}

	s.userService.On("UpdateUser", userID, req).Return(nil, errors.New("username already exists"))

	app := fiber.New()
	app.Put("/users/:id", s.userHandler.UpdateUser)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusConflict, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test DeleteUser handler
func (s *HandlerTestSuite) TestDeleteUser_Success() {
	userID := uint(1)
	s.userService.On("DeleteUser", userID).Return(nil)

	app := fiber.New()
	app.Delete("/users/:id", s.userHandler.DeleteUser)

	req := httptest.NewRequest("DELETE", "/users/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusNoContent, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestDeleteUser_InvalidID() {
	app := fiber.New()
	app.Delete("/users/:id", s.userHandler.DeleteUser)

	req := httptest.NewRequest("DELETE", "/users/invalid", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestDeleteUser_NotFound() {
	userID := uint(999)
	s.userService.On("DeleteUser", userID).Return(errors.New("user not found"))

	app := fiber.New()
	app.Delete("/users/:id", s.userHandler.DeleteUser)

	req := httptest.NewRequest("DELETE", "/users/999", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusNotFound, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test ListUsers handler
func (s *HandlerTestSuite) TestListUsers_Success() {
	limit := 10
	offset := 0

	expectedUsers := []*model.UserResponse{
		{
			ID:        1,
			Username:  "user1",
			Email:     "user1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Username:  "user2",
			Email:     "user2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	s.userService.On("ListUsers", limit, offset).Return(expectedUsers, nil)

	app := fiber.New()
	app.Get("/users", s.userHandler.ListUsers)

	req := httptest.NewRequest("GET", "/users", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestListUsers_WithQueryParams() {
	limit := 5
	offset := 10

	expectedUsers := []*model.UserResponse{
		{
			ID:        1,
			Username:  "user1",
			Email:     "user1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	s.userService.On("ListUsers", limit, offset).Return(expectedUsers, nil)

	app := fiber.New()
	app.Get("/users", s.userHandler.ListUsers)

	req := httptest.NewRequest("GET", "/users?limit=5&offset=10", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestListUsers_InvalidLimit() {
	limit := 10 // Default limit when invalid
	offset := 0

	expectedUsers := []*model.UserResponse{}
	s.userService.On("ListUsers", limit, offset).Return(expectedUsers, nil)

	app := fiber.New()
	app.Get("/users", s.userHandler.ListUsers)

	req := httptest.NewRequest("GET", "/users?limit=invalid", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestListUsers_ServiceError() {
	limit := 10
	offset := 0

	s.userService.On("ListUsers", limit, offset).Return(nil, errors.New("database error"))

	app := fiber.New()
	app.Get("/users", s.userHandler.ListUsers)

	req := httptest.NewRequest("GET", "/users", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusInternalServerError, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test GetUserProfile handler
func (s *HandlerTestSuite) TestGetUserProfile_Success() {
	userID := uint(1)
	expectedResponse := &model.UserResponse{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userService.On("GetUser", userID).Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/users/:id/profile", s.userHandler.GetUserProfile)

	req := httptest.NewRequest("GET", "/users/1/profile", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestGetUserProfile_InvalidID() {
	app := fiber.New()
	app.Get("/users/:id/profile", s.userHandler.GetUserProfile)

	req := httptest.NewRequest("GET", "/users/invalid/profile", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestGetUserProfile_NotFound() {
	userID := uint(999)
	s.userService.On("GetUser", userID).Return(nil, errors.New("user not found"))

	app := fiber.New()
	app.Get("/users/:id/profile", s.userHandler.GetUserProfile)

	req := httptest.NewRequest("GET", "/users/999/profile", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusNotFound, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test UpdateUserProfile handler
func (s *HandlerTestSuite) TestUpdateUserProfile_Success() {
	userID := uint(1)
	req := &model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	expectedResponse := &model.UserResponse{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userService.On("UpdateUser", userID, req).Return(expectedResponse, nil)

	app := fiber.New()
	app.Put("/users/:id/profile", s.userHandler.UpdateUserProfile)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("PUT", "/users/1/profile", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestUpdateUserProfile_InvalidID() {
	app := fiber.New()
	app.Put("/users/:id/profile", s.userHandler.UpdateUserProfile)

	req := httptest.NewRequest("PUT", "/users/invalid/profile", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestUpdateUserProfile_InvalidBody() {
	app := fiber.New()
	app.Put("/users/:id/profile", s.userHandler.UpdateUserProfile)

	req := httptest.NewRequest("PUT", "/users/1/profile", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestUpdateUserProfile_ValidationError() {
	req := &model.UpdateUserRequest{
		Email: "invalid-email",
	}

	app := fiber.New()
	app.Put("/users/:id/profile", s.userHandler.UpdateUserProfile)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("PUT", "/users/1/profile", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *HandlerTestSuite) TestUpdateUserProfile_ServiceError() {
	userID := uint(1)
	req := &model.UpdateUserRequest{
		Username: "updateduser",
	}

	s.userService.On("UpdateUser", userID, req).Return(nil, errors.New("username already exists"))

	app := fiber.New()
	app.Put("/users/:id/profile", s.userHandler.UpdateUserProfile)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("PUT", "/users/1/profile", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusConflict, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}
