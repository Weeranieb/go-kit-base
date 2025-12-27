package service

import (
	"errors"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/weeranieb/go-kit-base/src/internal/model"
)

func (s *ServiceTestSuite) TestCreateUser_Success() {
	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock repository calls
	s.userRepo.On("GetByEmail", req.Email).Return(nil, errors.New("not found"))
	s.userRepo.On("GetByUsername", req.Username).Return(nil, errors.New("not found"))

	expectedUser := &model.User{
		ID:        1,
		Username:  req.Username,
		Email:     req.Email,
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*model.User)
		user.ID = expectedUser.ID
		user.CreatedAt = expectedUser.CreatedAt
		user.UpdatedAt = expectedUser.UpdatedAt
	})

	// Execute
	result, err := s.userService.CreateUser(req)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Username, result.Username)
	assert.Equal(s.T(), req.Email, result.Email)
	assert.Equal(s.T(), expectedUser.ID, result.ID)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestCreateUser_EmailExists() {
	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	existingUser := &model.User{
		ID:       1,
		Username: "existing",
		Email:    req.Email,
	}

	s.userRepo.On("GetByEmail", req.Email).Return(existingUser, nil)

	// Execute
	result, err := s.userService.CreateUser(req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), "email already exists", err.Error())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestCreateUser_UsernameExists() {
	req := &model.CreateUserRequest{
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "password123",
	}

	existingUser := &model.User{
		ID:       1,
		Username: req.Username,
		Email:    "other@example.com",
	}

	s.userRepo.On("GetByEmail", req.Email).Return(nil, errors.New("not found"))
	s.userRepo.On("GetByUsername", req.Username).Return(existingUser, nil)

	// Execute
	result, err := s.userService.CreateUser(req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), "username already exists", err.Error())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestGetUser_Success() {
	userID := uint(1)
	expectedUser := &model.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userRepo.On("GetByID", userID).Return(expectedUser, nil)

	// Execute
	result, err := s.userService.GetUser(userID)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.ID, result.ID)
	assert.Equal(s.T(), expectedUser.Username, result.Username)
	assert.Equal(s.T(), expectedUser.Email, result.Email)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestGetUser_NotFound() {
	userID := uint(999)

	s.userRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	// Execute
	result, err := s.userService.GetUser(userID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestUpdateUser_Success() {
	userID := uint(1)
	req := &model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	existingUser := &model.User{
		ID:        userID,
		Username:  "olduser",
		Email:     "old@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userRepo.On("GetByID", userID).Return(existingUser, nil)
	s.userRepo.On("GetByUsername", req.Username).Return(nil, errors.New("not found"))
	s.userRepo.On("GetByEmail", req.Email).Return(nil, errors.New("not found"))
	s.userRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	// Execute
	result, err := s.userService.UpdateUser(userID, req)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Username, result.Username)
	assert.Equal(s.T(), req.Email, result.Email)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestUpdateUser_NotFound() {
	userID := uint(999)
	req := &model.UpdateUserRequest{
		Username: "updateduser",
	}

	s.userRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	// Execute
	result, err := s.userService.UpdateUser(userID, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestUpdateUser_UsernameExists() {
	userID := uint(1)
	req := &model.UpdateUserRequest{
		Username: "existinguser",
	}

	existingUser := &model.User{
		ID:       userID,
		Username: "olduser",
		Email:    "old@example.com",
	}

	conflictingUser := &model.User{
		ID:       2,
		Username: req.Username,
		Email:    "other@example.com",
	}

	s.userRepo.On("GetByID", userID).Return(existingUser, nil)
	s.userRepo.On("GetByUsername", req.Username).Return(conflictingUser, nil)

	// Execute
	result, err := s.userService.UpdateUser(userID, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), "username already exists", err.Error())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestDeleteUser_Success() {
	userID := uint(1)

	s.userRepo.On("Delete", userID).Return(nil)

	// Execute
	err := s.userService.DeleteUser(userID)

	// Assert
	assert.NoError(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestDeleteUser_Error() {
	userID := uint(999)

	s.userRepo.On("Delete", userID).Return(errors.New("delete failed"))

	// Execute
	err := s.userService.DeleteUser(userID)

	// Assert
	assert.Error(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestListUsers_Success() {
	limit := 10
	offset := 0

	expectedUsers := []*model.User{
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

	s.userRepo.On("List", limit, offset).Return(expectedUsers, nil)

	// Execute
	result, err := s.userService.ListUsers(limit, offset)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedUsers[0].ID, result[0].ID)
	assert.Equal(s.T(), expectedUsers[1].ID, result[1].ID)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestListUsers_Empty() {
	limit := 10
	offset := 0

	s.userRepo.On("List", limit, offset).Return([]*model.User{}, nil)

	// Execute
	result, err := s.userService.ListUsers(limit, offset)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 0)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestListUsers_Error() {
	limit := 10
	offset := 0

	s.userRepo.On("List", limit, offset).Return(nil, errors.New("database error"))

	// Execute
	result, err := s.userService.ListUsers(limit, offset)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.userRepo.AssertExpectations(s.T())
}
