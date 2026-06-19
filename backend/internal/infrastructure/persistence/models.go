package persistence

import (
	"time"

	"gorm.io/gorm"
)

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
	ID         string         `gorm:"primaryKey;type:varchar(36)"`
	UserID     string         `gorm:"type:varchar(36);not null;index"`
	StyleID    string         `gorm:"type:varchar(36);not null;index"`
	State      string         `gorm:"type:varchar(20);not null;index"`
	Config     string         `gorm:"type:jsonb;not null"` // JSON config
	Result     string         `gorm:"type:jsonb"`          // JSON result, nullable
	CreatedAt  time.Time      `gorm:"not null"`
	UpdatedAt  time.Time      `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name
func (WorkflowModel) TableName() string {
	return "workflows"
}
