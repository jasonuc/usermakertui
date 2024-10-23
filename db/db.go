package db

import (
	"errors"
	"fmt"
)

// Create a mock in-memory database to simulate user storage
var mockUserStorage = make(map[string]User)

func InitMockDB() {
	// Initialize the mock database with some initial users
	mockUserStorage["tac@hi.com"] = User{
		ID:       1,
		Email:    "tac@hi.com",
		Password: "thankY*Utac0",
	}

	mockUserStorage["jay@hi.com"] = User{
		ID:       2,
		Email:    "jay@hi.com",
		Password: "+h3ll0J_y$",
	}

	mockUserStorage["john@hi.com"] = User{
		ID:       3,
		Email:    "john@hi.com",
		Password: "superDOE1$",
	}

	mockUserStorage["alice@hi.com"] = User{
		ID:       4,
		Email:    "gopher@hi.com",
		Password: "go4£v£rRustn3v3r",
	}
}

// CreateUserParams represents the parameters required to create a user
type CreateUserParams struct {
	Email    string
	Password string
}

// User represents a user record in the mock database
type User struct {
	ID       int64
	Email    string
	Password string
}

// Q is a struct that acts as a placeholder for query methods, simulating a database connection
type Queries struct{}

// Q is a globally accessible instance of Queries
var Q = &Queries{}

// Auto-increment ID for new records, starting at 5
var autoIncrementID int64 = 5

func (q *Queries) CreateUser(params CreateUserParams) (User, error) {
	if _, exists := mockUserStorage[params.Email]; exists {
		return User{}, fmt.Errorf("user with email %s already exists", params.Email)
	}

	user := User{
		ID:       autoIncrementID,
		Email:    params.Email,
		Password: params.Password,
	}

	mockUserStorage[params.Email] = user
	autoIncrementID++

	return user, nil
}

func (q *Queries) SearchUser(email string) (User, error) {
	user, exists := mockUserStorage[email]
	if !exists {
		return User{}, errors.New("user not found")
	}

	return user, nil
}
