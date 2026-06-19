package persistence

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// JSONB is a custom type for JSONB columns
type JSONB map[string]interface{}

// Value implements driver.Valuer
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// StyleModel represents the database model for styles
type StyleModel struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)"`
	Name        string         `gorm:"type:varchar(100);not null;index"`
	Description string         `gorm:"type:varchar(500)"`
	Prompt      string         `gorm:"type:text;not null"`
	Category    string         `gorm:"type:varchar(50);index"`
	Tags        string         `gorm:"type:text"` // JSON array stored as text
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name
func (StyleModel) TableName() string {
	return "styles"
}

// UserModel represents the database model for users
type UserModel struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)"`
	Email        string         `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Name         string         `gorm:"type:varchar(100);not null"`
	Role         string         `gorm:"type:varchar(20);not null;default:'user'"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name
func (UserModel) TableName() string {
	return "users"
}

// WorkflowModel represents the database model for workflows
type WorkflowModel struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)"`
	UserID    string         `gorm:"type:varchar(36);not null;index"`
	StyleID   string         `gorm:"type:varchar(36);not null;index"`
	State     string         `gorm:"type:varchar(20);not null;index"`
	Config    JSONB          `gorm:"type:jsonb;not null"` // JSON config
	Result    JSONB          `gorm:"type:jsonb"`          // JSON result, nullable
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name
func (WorkflowModel) TableName() string {
	return "workflows"
}
