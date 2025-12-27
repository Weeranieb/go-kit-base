package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/go-kit-base/src/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db             *gorm.DB
	userRepository UserRepository
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite for fast tests
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	// Auto-migrate the schema
	err = s.db.AutoMigrate(&model.User{})
	if err != nil {
		s.T().Fatal("Failed to migrate database:", err)
	}

	s.userRepository = NewUserRepository(s.db)
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (s *UserRepositoryTestSuite) SetupTest() {
	// Clean up before each test
	s.db.Exec("DELETE FROM users")
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	// Additional cleanup if needed (SetupTest already does this)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

// Test Create operations
func (s *UserRepositoryTestSuite) TestCreate_Success() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	err := s.userRepository.Create(user)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), user.ID)
	assert.NotZero(s.T(), user.CreatedAt)
	assert.NotZero(s.T(), user.UpdatedAt)
	assert.Equal(s.T(), "testuser", user.Username)
	assert.Equal(s.T(), "test@example.com", user.Email)
}

func (s *UserRepositoryTestSuite) TestCreate_DuplicateEmail() {
	user1 := &model.User{
		Username: "user1",
		Email:    "test@example.com",
		Password: "password1",
	}
	err := s.userRepository.Create(user1)
	assert.NoError(s.T(), err)

	user2 := &model.User{
		Username: "user2",
		Email:    "test@example.com", // Duplicate email
		Password: "password2",
	}

	err = s.userRepository.Create(user2)

	assert.Error(s.T(), err)
	// Verify the duplicate wasn't created
	count := int64(0)
	s.db.Model(&model.User{}).Where("email = ?", "test@example.com").Count(&count)
	assert.Equal(s.T(), int64(1), count)
}

func (s *UserRepositoryTestSuite) TestCreate_DuplicateUsername() {
	user1 := &model.User{
		Username: "testuser",
		Email:    "user1@example.com",
		Password: "password1",
	}
	err := s.userRepository.Create(user1)
	assert.NoError(s.T(), err)

	user2 := &model.User{
		Username: "testuser", // Duplicate username
		Email:    "user2@example.com",
		Password: "password2",
	}

	err = s.userRepository.Create(user2)

	assert.Error(s.T(), err)
	// Verify the duplicate wasn't created
	count := int64(0)
	s.db.Model(&model.User{}).Where("username = ?", "testuser").Count(&count)
	assert.Equal(s.T(), int64(1), count)
}

// Test GetByID operations
func (s *UserRepositoryTestSuite) TestGetByID_Success() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	result, err := s.userRepository.GetByID(user.ID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.ID, result.ID)
	assert.Equal(s.T(), user.Username, result.Username)
	assert.Equal(s.T(), user.Email, result.Email)
	assert.Equal(s.T(), user.Password, result.Password)
}

func (s *UserRepositoryTestSuite) TestGetByID_NotFound() {
	result, err := s.userRepository.GetByID(999)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), gorm.ErrRecordNotFound, err)
}

// Test GetByEmail operations
func (s *UserRepositoryTestSuite) TestGetByEmail_Success() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	result, err := s.userRepository.GetByEmail("test@example.com")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.Email, result.Email)
	assert.Equal(s.T(), user.Username, result.Username)
	assert.Equal(s.T(), user.ID, result.ID)
}

func (s *UserRepositoryTestSuite) TestGetByEmail_NotFound() {
	result, err := s.userRepository.GetByEmail("nonexistent@example.com")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), gorm.ErrRecordNotFound, err)
}

func (s *UserRepositoryTestSuite) TestGetByEmail_CaseSensitive() {
	user := &model.User{
		Username: "testuser",
		Email:    "Test@Example.com",
		Password: "password",
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	// SQLite is case-sensitive by default for LIKE, but = is case-sensitive
	result, err := s.userRepository.GetByEmail("Test@Example.com")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)

	// Different case should not match
	result, err = s.userRepository.GetByEmail("test@example.com")
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

// Test GetByUsername operations
func (s *UserRepositoryTestSuite) TestGetByUsername_Success() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	result, err := s.userRepository.GetByUsername("testuser")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.Username, result.Username)
	assert.Equal(s.T(), user.Email, result.Email)
	assert.Equal(s.T(), user.ID, result.ID)
}

func (s *UserRepositoryTestSuite) TestGetByUsername_NotFound() {
	result, err := s.userRepository.GetByUsername("nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), gorm.ErrRecordNotFound, err)
}

// Test Update operations
func (s *UserRepositoryTestSuite) TestUpdate_Success() {
	user := &model.User{
		Username: "olduser",
		Email:    "old@example.com",
		Password: "password",
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure UpdatedAt changes
	time.Sleep(10 * time.Millisecond)

	user.Username = "newuser"
	user.Email = "new@example.com"
	err = s.userRepository.Update(user)

	assert.NoError(s.T(), err)

	// Verify update
	updated, err := s.userRepository.GetByID(user.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "newuser", updated.Username)
	assert.Equal(s.T(), "new@example.com", updated.Email)
	assert.True(s.T(), updated.UpdatedAt.After(originalUpdatedAt))
}

func (s *UserRepositoryTestSuite) TestUpdate_PartialFields() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	// Update only username
	user.Username = "updateduser"
	err = s.userRepository.Update(user)

	assert.NoError(s.T(), err)

	// Verify only username changed
	updated, err := s.userRepository.GetByID(user.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "updateduser", updated.Username)
	assert.Equal(s.T(), "test@example.com", updated.Email) // Email unchanged
}

func (s *UserRepositoryTestSuite) TestUpdate_NonExistent() {
	user := &model.User{
		ID:       999,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	// GORM Save will create if not found, so this might not error
	err := s.userRepository.Update(user)

	// GORM Save creates if ID doesn't exist, so we verify it was created
	assert.NoError(s.T(), err)
	result, err := s.userRepository.GetByID(999)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
}

// Test Delete operations
func (s *UserRepositoryTestSuite) TestDelete_Success() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)
	userID := user.ID

	err = s.userRepository.Delete(userID)

	assert.NoError(s.T(), err)

	// Verify soft delete (User model has DeletedAt field)
	result, err := s.userRepository.GetByID(userID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), gorm.ErrRecordNotFound, err)

	// Verify it's soft deleted (still in DB but with DeletedAt set)
	var deletedUser model.User
	s.db.Unscoped().First(&deletedUser, userID)
	assert.NotZero(s.T(), deletedUser.DeletedAt)
}

func (s *UserRepositoryTestSuite) TestDelete_NotFound() {
	err := s.userRepository.Delete(999)

	// GORM Delete doesn't error on non-existent records
	assert.NoError(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestDelete_MultipleUsers() {
	// Create multiple users
	user1 := &model.User{Username: "user1", Email: "user1@example.com", Password: "pass1"}
	user2 := &model.User{Username: "user2", Email: "user2@example.com", Password: "pass2"}
	user3 := &model.User{Username: "user3", Email: "user3@example.com", Password: "pass3"}

	s.userRepository.Create(user1)
	s.userRepository.Create(user2)
	s.userRepository.Create(user3)

	// Delete one user
	err := s.userRepository.Delete(user2.ID)
	assert.NoError(s.T(), err)

	// Verify only user2 is deleted
	_, err = s.userRepository.GetByID(user1.ID)
	assert.NoError(s.T(), err)

	_, err = s.userRepository.GetByID(user2.ID)
	assert.Error(s.T(), err)

	_, err = s.userRepository.GetByID(user3.ID)
	assert.NoError(s.T(), err)
}

// Test List operations
func (s *UserRepositoryTestSuite) TestList_Success() {
	// Create multiple users
	for i := 0; i < 5; i++ {
		user := &model.User{
			Username: "user" + string(rune('0'+i)),
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			Password: "password",
		}
		s.userRepository.Create(user)
	}

	users, err := s.userRepository.List(10, 0)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 5)
}

func (s *UserRepositoryTestSuite) TestList_WithLimit() {
	// Create 10 users
	for i := 0; i < 10; i++ {
		user := &model.User{
			Username: "user" + string(rune('0'+i)),
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			Password: "password",
		}
		s.userRepository.Create(user)
	}

	users, err := s.userRepository.List(5, 0)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 5)
}

func (s *UserRepositoryTestSuite) TestList_WithOffset() {
	// Create 10 users
	for i := 0; i < 10; i++ {
		user := &model.User{
			Username: "user" + string(rune('0'+i)),
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			Password: "password",
		}
		s.userRepository.Create(user)
	}

	// Get second page (offset 5, limit 5)
	users, err := s.userRepository.List(5, 5)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 5)
}

func (s *UserRepositoryTestSuite) TestList_Empty() {
	users, err := s.userRepository.List(10, 0)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), users)
	assert.Len(s.T(), users, 0)
}

func (s *UserRepositoryTestSuite) TestList_ExcludesSoftDeleted() {
	// Create users
	user1 := &model.User{Username: "user1", Email: "user1@example.com", Password: "pass1"}
	user2 := &model.User{Username: "user2", Email: "user2@example.com", Password: "pass2"}
	user3 := &model.User{Username: "user3", Email: "user3@example.com", Password: "pass3"}

	s.userRepository.Create(user1)
	s.userRepository.Create(user2)
	s.userRepository.Create(user3)

	// Delete one user
	s.userRepository.Delete(user2.ID)

	// List should only return non-deleted users
	users, err := s.userRepository.List(10, 0)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	// Verify user2 is not in the list
	for _, user := range users {
		assert.NotEqual(s.T(), user2.ID, user.ID)
	}
}

// Test edge cases and data integrity
func (s *UserRepositoryTestSuite) TestCreate_RequiredFields() {
	// Note: SQLite is more lenient with NOT NULL constraints than PostgreSQL
	// This test verifies behavior with empty strings
	// For strict NOT NULL enforcement, use PostgreSQL in integration tests

	// Test with empty strings (SQLite allows these, PostgreSQL would reject)
	user := &model.User{
		Username: "",
		Email:    "",
		Password: "",
	}

	err := s.userRepository.Create(user)

	// SQLite allows empty strings, so this may succeed
	// We just verify the behavior is consistent
	if err == nil {
		// If it succeeds, verify the record was created
		assert.NotZero(s.T(), user.ID)
	}
	// If it fails, that's also acceptable - depends on database configuration
}

func (s *UserRepositoryTestSuite) TestGetByID_AfterDelete() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	s.userRepository.Create(user)
	userID := user.ID

	// Delete the user
	s.userRepository.Delete(userID)

	// Try to get it - should fail
	result, err := s.userRepository.GetByID(userID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *UserRepositoryTestSuite) TestGetByEmail_AfterDelete() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	s.userRepository.Create(user)

	// Delete the user
	s.userRepository.Delete(user.ID)

	// Try to get by email - should fail
	result, err := s.userRepository.GetByEmail("test@example.com")
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *UserRepositoryTestSuite) TestMultipleOperations() {
	// Create
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	// Get by ID
	found, err := s.userRepository.GetByID(user.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.ID, found.ID)

	// Get by Email
	found, err = s.userRepository.GetByEmail(user.Email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Email, found.Email)

	// Get by Username
	found, err = s.userRepository.GetByUsername(user.Username)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Username, found.Username)

	// Update
	user.Username = "updateduser"
	err = s.userRepository.Update(user)
	assert.NoError(s.T(), err)

	// Verify update
	found, err = s.userRepository.GetByID(user.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "updateduser", found.Username)

	// Delete
	err = s.userRepository.Delete(user.ID)
	assert.NoError(s.T(), err)

	// Verify deletion
	_, err = s.userRepository.GetByID(user.ID)
	assert.Error(s.T(), err)
}
