package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Role values for RBAC.
const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

// Asset status values.
const (
	AssetAvailable   = "available"
	AssetAssigned    = "assigned"
	AssetInTransit   = "in_transit"
	AssetMaintenance = "maintenance"
	AssetReserved    = "reserved"
	AssetRetired     = "retired"
	AssetSold        = "sold"
)

// User is an application account.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string    `gorm:"size:200;not null" json:"name"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:32;not null;index" json:"role"`
	Phone        *string   `gorm:"size:64" json:"phone,omitempty"`
	AvatarURL    *string   `gorm:"size:512" json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Location is a depot, job site, or workshop.
type Location struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Type        string    `gorm:"size:32;not null" json:"type"`
	Address     string    `gorm:"size:500" json:"address"`
	Lat         *float64  `json:"lat,omitempty"`
	Lng         *float64  `json:"lng,omitempty"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Asset is a piece of equipment.
type Asset struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AssetCode         string         `gorm:"size:64;uniqueIndex;not null" json:"asset_code"`
	Name              string         `gorm:"size:200;not null" json:"name"`
	Model             string         `gorm:"size:200" json:"model"`
	Type              string         `gorm:"size:64;not null;index" json:"type"`
	Category          string         `gorm:"size:64;not null" json:"category"`
	Specs             datatypes.JSON `gorm:"type:jsonb" json:"specs"`
	Status            string         `gorm:"size:32;not null;index" json:"status"`
	CurrentLocationID *uuid.UUID     `gorm:"type:uuid;index" json:"current_location_id,omitempty"`
	PurchaseDate      *time.Time     `gorm:"type:date" json:"purchase_date,omitempty"`
	PurchaseCost      *float64       `json:"purchase_cost,omitempty"`
	HourlyRate        *float64       `json:"hourly_rate,omitempty"`
	PhotoURL          *string        `gorm:"size:512" json:"photo_url,omitempty"`
	QRCode            *string        `gorm:"size:128" json:"qr_code,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`

	CurrentLocation *Location `gorm:"foreignKey:CurrentLocationID" json:"current_location,omitempty"`
}

// Project is a job site / work order container.
type Project struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string     `gorm:"size:200;not null" json:"name"`
	Address          string     `gorm:"size:500" json:"address"`
	LocationID       *uuid.UUID `gorm:"type:uuid" json:"location_id,omitempty"`
	StartDate        time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate          time.Time  `gorm:"type:date;not null" json:"end_date"`
	Status           string     `gorm:"size:32;not null;index" json:"status"`
	ProjectManagerID *uuid.UUID `gorm:"type:uuid" json:"project_manager_id,omitempty"`
	Description      string     `gorm:"type:text" json:"description"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Location *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
}

// Assignment links an asset to a project for a time window.
type Assignment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AssetID    uuid.UUID `gorm:"type:uuid;not null;index" json:"asset_id"`
	ProjectID  uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	StartTime  time.Time `gorm:"not null;index" json:"start_time"`
	EndTime    time.Time `gorm:"not null;index" json:"end_time"`
	Status     string    `gorm:"size:32;not null;index" json:"status"`
	AssignedBy uuid.UUID `gorm:"type:uuid;not null" json:"assigned_by"`
	Notes      string    `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Asset   *Asset   `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// MaintenanceRecord blocks availability for a service window.
type MaintenanceRecord struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	AssetID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"asset_id"`
	Type           string     `gorm:"size:32;not null" json:"type"`
	Description    string     `gorm:"type:text" json:"description"`
	ScheduledStart time.Time  `gorm:"not null;index" json:"scheduled_start"`
	ScheduledEnd   time.Time  `gorm:"not null;index" json:"scheduled_end"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	Status         string     `gorm:"size:32;not null;index" json:"status"`
	Cost           *float64   `json:"cost,omitempty"`
	PerformedBy    *string    `gorm:"size:200" json:"performed_by,omitempty"`
	Notes          string     `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Asset *Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
}

// RefreshToken stores opaque refresh tokens hashed.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"size:128;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
