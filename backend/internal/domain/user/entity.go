package user

import (
	"errors"
	"regexp"
	"time"
)

// Entity represents a user in the domain
type Entity struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Role represents user role
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Validation errors
var (
	ErrEmptyEmail       = errors.New("email cannot be empty")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrInvalidRole      = errors.New("invalid user role")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// New creates a new User entity with validation
func New(email, passwordHash, name string, role Role) (*Entity, error) {
	user := &Entity{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate checks if the user entity is valid
func (u *Entity) Validate() error {
	if u.Email == "" {
		return ErrEmptyEmail
	}
	if !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	if u.PasswordHash == "" {
		return ErrEmptyPassword
	}
	if u.Name == "" {
		return ErrEmptyName
	}
	if u.Role != RoleUser && u.Role != RoleAdmin {
		return ErrInvalidRole
	}
	return nil
}

// Update modifies the user entity
func (u *Entity) Update(email, name string, role Role) error {
	u.Email = email
	u.Name = name
	u.Role = role
	u.UpdatedAt = time.Now()

	return u.Validate()
}

// IsAdmin checks if user has admin role
func (u *Entity) IsAdmin() bool {
	return u.Role == RoleAdmin
}
