package dto

import (
	"time"

	"github.com/google/uuid"
)

// APIError is the standard error envelope.
type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ErrorResponse wraps APIError.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// LoginRequest is the login body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is a safe user payload.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Phone     *string   `json:"phone,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}

// AuthResponse is returned on login.
type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// CreateUserRequest creates a user (admin).
type CreateUserRequest struct {
	Name     string  `json:"name" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Role     string  `json:"role" binding:"required"`
	Phone    *string `json:"phone"`
}

// UpdateUserRequest patches a user.
type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Role  *string `json:"role"`
	Phone *string `json:"phone"`
}

// LocationRequest creates/updates a location.
type LocationRequest struct {
	Name        string   `json:"name" binding:"required"`
	Type        string   `json:"type" binding:"required"`
	Address     string   `json:"address"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	Description *string  `json:"description"`
}

// AssetRequest creates/updates an asset.
type AssetRequest struct {
	AssetCode         string                 `json:"asset_code" binding:"required"`
	Name              string                 `json:"name" binding:"required"`
	Model             string                 `json:"model"`
	Type              string                 `json:"type" binding:"required"`
	Category          string                 `json:"category" binding:"required"`
	Specs             map[string]interface{} `json:"specs"`
	Status            string                 `json:"status"`
	CurrentLocationID *uuid.UUID             `json:"current_location_id"`
	PurchaseCost      *float64               `json:"purchase_cost"`
	HourlyRate        *float64               `json:"hourly_rate"`
	PhotoURL          *string                `json:"photo_url"`
}

// ProjectRequest creates/updates a project.
type ProjectRequest struct {
	Name             string     `json:"name" binding:"required"`
	Address          string     `json:"address"`
	LocationID       *uuid.UUID `json:"location_id"`
	StartDate        string     `json:"start_date" binding:"required"`
	EndDate          string     `json:"end_date" binding:"required"`
	Status           string     `json:"status"`
	ProjectManagerID *uuid.UUID `json:"project_manager_id"`
	Description      string     `json:"description"`
}

// AssignmentRequest creates/updates an assignment.
type AssignmentRequest struct {
	AssetID   uuid.UUID `json:"asset_id" binding:"required"`
	ProjectID uuid.UUID `json:"project_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Notes     string    `json:"notes"`
}

// MaintenanceRequest creates/updates maintenance.
type MaintenanceRequest struct {
	AssetID        uuid.UUID `json:"asset_id" binding:"required"`
	Type           string    `json:"type" binding:"required"`
	Description    string    `json:"description"`
	ScheduledStart time.Time `json:"scheduled_start" binding:"required"`
	ScheduledEnd   time.Time `json:"scheduled_end" binding:"required"`
	Status         string    `json:"status"`
	Cost           *float64  `json:"cost"`
	PerformedBy    *string   `json:"performed_by"`
	Notes          string    `json:"notes"`
}

// DispatchSearchRequest is the Smart Dispatch Finder input.
type DispatchSearchRequest struct {
	EquipmentType       string     `json:"equipment_type" binding:"required"`
	StartTime           time.Time  `json:"start_time" binding:"required"`
	EndTime             time.Time  `json:"end_time" binding:"required"`
	ProjectID           *uuid.UUID `json:"project_id"`
	JobLat              *float64   `json:"job_lat"`
	JobLng              *float64   `json:"job_lng"`
	MinWeightT          *float64   `json:"min_weight_t"`
	IncludeAlternatives bool       `json:"include_alternatives"`
}

// DispatchResultItem is one ranked dispatch candidate.
type DispatchResultItem struct {
	Asset             interface{} `json:"asset"`
	Available         bool        `json:"available"`
	AvailableFrom     time.Time   `json:"available_from"`
	DistanceKm        *float64    `json:"distance_km,omitempty"`
	TransportMinutes  *float64    `json:"transport_minutes,omitempty"`
	SuitabilityScore  float64     `json:"suitability_score"`
	RankScore         float64     `json:"rank_score"`
	IsAlternative     bool        `json:"is_alternative"`
	Warnings          []string    `json:"warnings"`
}

// DashboardSummary is the dashboard KPI payload.
type DashboardSummary struct {
	FleetByStatus       map[string]int64 `json:"fleet_by_status"`
	OverdueMaintenance  int64            `json:"overdue_maintenance"`
	ActiveAssignments   int64            `json:"active_assignments"`
	UtilizationPercent  float64          `json:"utilization_percent"`
	UpcomingAssignments int64            `json:"upcoming_assignments"`
}
