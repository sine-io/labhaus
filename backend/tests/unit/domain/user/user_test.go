package user_test

import (
	"testing"

	"github.com/labhaus/backend/internal/domain/user"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		passwordHash string
		userName     string
		role         user.Role
		wantErr      error
	}{
		{
			name:         "valid user",
			email:        "test@example.com",
			passwordHash: "hashedpassword123",
			userName:     "Test User",
			role:         user.RoleUser,
			wantErr:      nil,
		},
		{
			name:         "valid admin",
			email:        "admin@example.com",
			passwordHash: "hashedpassword123",
			userName:     "Admin User",
			role:         user.RoleAdmin,
			wantErr:      nil,
		},
		{
			name:         "empty email",
			email:        "",
			passwordHash: "hashedpassword123",
			userName:     "Test",
			role:         user.RoleUser,
			wantErr:      user.ErrEmptyEmail,
		},
		{
			name:         "invalid email format",
			email:        "invalid-email",
			passwordHash: "hashedpassword123",
			userName:     "Test",
			role:         user.RoleUser,
			wantErr:      user.ErrInvalidEmail,
		},
		{
			name:         "empty password hash",
			email:        "test@example.com",
			passwordHash: "",
			userName:     "Test",
			role:         user.RoleUser,
			wantErr:      user.ErrEmptyPassword,
		},
		{
			name:         "empty name",
			email:        "test@example.com",
			passwordHash: "hashedpassword123",
			userName:     "",
			role:         user.RoleUser,
			wantErr:      user.ErrEmptyName,
		},
		{
			name:         "invalid role",
			email:        "test@example.com",
			passwordHash: "hashedpassword123",
			userName:     "Test",
			role:         "invalid",
			wantErr:      user.ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := user.New(tt.email, tt.passwordHash, tt.userName, tt.role)
			if err != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && u == nil {
				t.Error("New() returned nil user with no error")
			}
			if err == nil {
				if u.Email != tt.email {
					t.Errorf("Email = %v, want %v", u.Email, tt.email)
				}
				if u.CreatedAt.IsZero() {
					t.Error("CreatedAt should not be zero")
				}
				if u.UpdatedAt.IsZero() {
					t.Error("UpdatedAt should not be zero")
				}
			}
		})
	}
}

func TestEntity_Update(t *testing.T) {
	u, err := user.New("test@example.com", "hash123", "Test User", user.RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	originalUpdatedAt := u.UpdatedAt

	err = u.Update("updated@example.com", "Updated Name", user.RoleAdmin)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if u.Email != "updated@example.com" {
		t.Errorf("Email = %v, want updated@example.com", u.Email)
	}
	if u.Name != "Updated Name" {
		t.Errorf("Name = %v, want Updated Name", u.Name)
	}
	if u.Role != user.RoleAdmin {
		t.Errorf("Role = %v, want RoleAdmin", u.Role)
	}
	if u.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt should be updated")
	}

	// Test invalid update
	err = u.Update("", "Name", user.RoleUser)
	if err != user.ErrEmptyEmail {
		t.Errorf("Update() with empty email should return ErrEmptyEmail, got %v", err)
	}
}

func TestEntity_IsAdmin(t *testing.T) {
	adminUser, _ := user.New("admin@example.com", "hash", "Admin", user.RoleAdmin)
	normalUser, _ := user.New("user@example.com", "hash", "User", user.RoleUser)

	if !adminUser.IsAdmin() {
		t.Error("Admin user should return true for IsAdmin()")
	}

	if normalUser.IsAdmin() {
		t.Error("Normal user should return false for IsAdmin()")
	}
}

func TestEntity_Validate(t *testing.T) {
	validUser, _ := user.New("test@example.com", "hash123", "Test", user.RoleUser)

	if err := validUser.Validate(); err != nil {
		t.Errorf("Validate() should pass for valid user, got %v", err)
	}

	// Test various invalid cases
	invalidCases := []struct {
		name string
		user *user.Entity
		want error
	}{
		{
			name: "empty email",
			user: &user.Entity{Email: "", PasswordHash: "hash", Name: "Test", Role: user.RoleUser},
			want: user.ErrEmptyEmail,
		},
		{
			name: "invalid email",
			user: &user.Entity{Email: "invalid", PasswordHash: "hash", Name: "Test", Role: user.RoleUser},
			want: user.ErrInvalidEmail,
		},
		{
			name: "empty password",
			user: &user.Entity{Email: "test@example.com", PasswordHash: "", Name: "Test", Role: user.RoleUser},
			want: user.ErrEmptyPassword,
		},
		{
			name: "empty name",
			user: &user.Entity{Email: "test@example.com", PasswordHash: "hash", Name: "", Role: user.RoleUser},
			want: user.ErrEmptyName,
		},
		{
			name: "invalid role",
			user: &user.Entity{Email: "test@example.com", PasswordHash: "hash", Name: "Test", Role: "superuser"},
			want: user.ErrInvalidRole,
		},
	}

	for _, tt := range invalidCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.user.Validate(); err != tt.want {
				t.Errorf("Validate() error = %v, want %v", err, tt.want)
			}
		})
	}
}
