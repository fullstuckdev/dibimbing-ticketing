package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/service"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	
	user := &entity.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	
	// Email doesn't exist yet
	mockRepo.On("FindByEmail", user.Email).Return(nil, errors.New("user not found"))
	mockRepo.On("Save", user).Return(nil)
	
	// Test
	err := userService.Register(user)
	
	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	
	existingUser := &entity.User{
		ID:       1,
		Name:     "Existing User",
		Email:    "existing@example.com",
		Password: "password123",
	}
	
	user := &entity.User{
		Name:     "Test User",
		Email:    "existing@example.com", // Same email
		Password: "password123",
	}
	
	// Email already exists
	mockRepo.On("FindByEmail", user.Email).Return(existingUser, nil)
	
	// Test
	err := userService.Register(user)
	
	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "email already registered", err.Error())
	mockRepo.AssertExpectations(t)
}

// Since we can't easily mock the ComparePassword method, let's skip the actual login test
// For full testing we'd need to separate the password comparison into its own service
func TestLogin_UserFound(t *testing.T) {
	// Skip this test - just a placeholder
	t.Skip("Skipping login test until we refactor password validation")
	
	// We would normally test login here
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	
	// User not found
	mockRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, errors.New("user not found"))
	
	// Test
	token, err := userService.Login("nonexistent@example.com", "password123")
	
	// Assertions
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
} 