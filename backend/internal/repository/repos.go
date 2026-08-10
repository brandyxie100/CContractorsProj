package repository

import (
	"time"

	"github.com/clementscontractors/equipment/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repos aggregates data access.
type Repos struct {
	DB *gorm.DB
}

// New creates a Repos wrapper.
func New(db *gorm.DB) *Repos {
	return &Repos{DB: db}
}

func (r *Repos) CreateUser(u *models.User) error {
	return r.DB.Create(u).Error
}

func (r *Repos) FindUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.DB.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repos) FindUserByID(id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.DB.First(&u, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repos) ListUsers() ([]models.User, error) {
	var list []models.User
	err := r.DB.Order("name").Find(&list).Error
	return list, err
}

func (r *Repos) UpdateUser(u *models.User) error {
	return r.DB.Save(u).Error
}

func (r *Repos) DeleteUser(id uuid.UUID) error {
	return r.DB.Delete(&models.User{}, "id = ?", id).Error
}

func (r *Repos) SaveRefreshToken(t *models.RefreshToken) error {
	return r.DB.Create(t).Error
}

func (r *Repos) FindRefreshToken(hash string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.DB.Where("token_hash = ? AND expires_at > ?", hash, time.Now()).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repos) DeleteRefreshToken(hash string) error {
	return r.DB.Where("token_hash = ?", hash).Delete(&models.RefreshToken{}).Error
}

func (r *Repos) DeleteRefreshTokensForUser(userID uuid.UUID) error {
	return r.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

// --- Locations ---

func (r *Repos) ListLocations() ([]models.Location, error) {
	var list []models.Location
	err := r.DB.Order("name").Find(&list).Error
	return list, err
}

func (r *Repos) CreateLocation(l *models.Location) error {
	return r.DB.Create(l).Error
}

func (r *Repos) FindLocation(id uuid.UUID) (*models.Location, error) {
	var l models.Location
	err := r.DB.First(&l, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *Repos) UpdateLocation(l *models.Location) error {
	return r.DB.Save(l).Error
}

func (r *Repos) DeleteLocation(id uuid.UUID) error {
	return r.DB.Delete(&models.Location{}, "id = ?", id).Error
}

// --- Assets ---

type AssetFilter struct {
	Type       string
	Status     string
	LocationID *uuid.UUID
	Q          string
}

func (r *Repos) ListAssets(f AssetFilter) ([]models.Asset, error) {
	q := r.DB.Preload("CurrentLocation").Model(&models.Asset{})
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.LocationID != nil {
		q = q.Where("current_location_id = ?", *f.LocationID)
	}
	if f.Q != "" {
		like := "%" + f.Q + "%"
		q = q.Where("asset_code ILIKE ? OR name ILIKE ? OR model ILIKE ?", like, like, like)
	}
	var list []models.Asset
	err := q.Order("asset_code").Find(&list).Error
	return list, err
}

func (r *Repos) CreateAsset(a *models.Asset) error {
	return r.DB.Create(a).Error
}

func (r *Repos) FindAsset(id uuid.UUID) (*models.Asset, error) {
	var a models.Asset
	err := r.DB.Preload("CurrentLocation").First(&a, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repos) UpdateAsset(a *models.Asset) error {
	return r.DB.Save(a).Error
}

func (r *Repos) DeleteAsset(id uuid.UUID) error {
	return r.DB.Delete(&models.Asset{}, "id = ?", id).Error
}

func (r *Repos) ListAssetsByTypes(types []string) ([]models.Asset, error) {
	var list []models.Asset
	err := r.DB.Preload("CurrentLocation").Where("type IN ? AND status NOT IN ?", types, []string{models.AssetRetired, models.AssetSold}).Find(&list).Error
	return list, err
}

// --- Projects ---

func (r *Repos) ListProjects(status string) ([]models.Project, error) {
	q := r.DB.Preload("Location").Model(&models.Project{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var list []models.Project
	err := q.Order("start_date desc").Find(&list).Error
	return list, err
}

func (r *Repos) CreateProject(p *models.Project) error {
	return r.DB.Create(p).Error
}

func (r *Repos) FindProject(id uuid.UUID) (*models.Project, error) {
	var p models.Project
	err := r.DB.Preload("Location").First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repos) UpdateProject(p *models.Project) error {
	return r.DB.Save(p).Error
}

func (r *Repos) DeleteProject(id uuid.UUID) error {
	return r.DB.Delete(&models.Project{}, "id = ?", id).Error
}

// --- Assignments ---

func (r *Repos) ListAssignments(from, to *time.Time, assetID, projectID *uuid.UUID) ([]models.Assignment, error) {
	q := r.DB.Preload("Asset").Preload("Project").Where("status <> ?", "cancelled")
	if from != nil {
		q = q.Where("end_time > ?", *from)
	}
	if to != nil {
		q = q.Where("start_time < ?", *to)
	}
	if assetID != nil {
		q = q.Where("asset_id = ?", *assetID)
	}
	if projectID != nil {
		q = q.Where("project_id = ?", *projectID)
	}
	var list []models.Assignment
	err := q.Order("start_time").Find(&list).Error
	return list, err
}

func (r *Repos) CreateAssignment(a *models.Assignment) error {
	return r.DB.Create(a).Error
}

func (r *Repos) FindAssignment(id uuid.UUID) (*models.Assignment, error) {
	var a models.Assignment
	err := r.DB.Preload("Asset").Preload("Project").First(&a, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repos) UpdateAssignment(a *models.Assignment) error {
	return r.DB.Save(a).Error
}

func (r *Repos) FindOverlappingAssignments(assetID uuid.UUID, start, end time.Time, excludeID *uuid.UUID) ([]models.Assignment, error) {
	q := r.DB.Where("asset_id = ? AND status NOT IN ? AND start_time < ? AND end_time > ?",
		assetID, []string{"cancelled", "completed"}, end, start)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var list []models.Assignment
	err := q.Find(&list).Error
	return list, err
}

// --- Maintenance ---

func (r *Repos) ListMaintenance(assetID *uuid.UUID, status string) ([]models.MaintenanceRecord, error) {
	q := r.DB.Preload("Asset").Model(&models.MaintenanceRecord{})
	if assetID != nil {
		q = q.Where("asset_id = ?", *assetID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var list []models.MaintenanceRecord
	err := q.Order("scheduled_start desc").Find(&list).Error
	return list, err
}

func (r *Repos) CreateMaintenance(m *models.MaintenanceRecord) error {
	return r.DB.Create(m).Error
}

func (r *Repos) FindMaintenance(id uuid.UUID) (*models.MaintenanceRecord, error) {
	var m models.MaintenanceRecord
	err := r.DB.Preload("Asset").First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repos) UpdateMaintenance(m *models.MaintenanceRecord) error {
	return r.DB.Save(m).Error
}

func (r *Repos) DeleteMaintenance(id uuid.UUID) error {
	return r.DB.Delete(&models.MaintenanceRecord{}, "id = ?", id).Error
}

func (r *Repos) FindOverlappingMaintenance(assetID uuid.UUID, start, end time.Time) ([]models.MaintenanceRecord, error) {
	var list []models.MaintenanceRecord
	err := r.DB.Where(
		"asset_id = ? AND status IN ? AND scheduled_start < ? AND scheduled_end > ?",
		assetID, []string{"scheduled", "in_progress", "overdue"}, end, start,
	).Find(&list).Error
	return list, err
}

// --- Dashboard aggregates ---

func (r *Repos) CountAssetsByStatus() (map[string]int64, error) {
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	err := r.DB.Model(&models.Asset{}).Select("status, count(*) as count").Group("status").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := map[string]int64{}
	for _, rw := range rows {
		out[rw.Status] = rw.Count
	}
	return out, nil
}

func (r *Repos) CountOverdueMaintenance() (int64, error) {
	var n int64
	err := r.DB.Model(&models.MaintenanceRecord{}).Where("status = ? OR (status IN ? AND scheduled_end < ?)",
		"overdue", []string{"scheduled", "in_progress"}, time.Now()).Count(&n).Error
	return n, err
}

func (r *Repos) CountActiveAssignments() (int64, error) {
	now := time.Now()
	var n int64
	err := r.DB.Model(&models.Assignment{}).
		Where("status IN ? AND start_time <= ? AND end_time > ?", []string{"scheduled", "active"}, now, now).
		Count(&n).Error
	return n, err
}

func (r *Repos) CountUpcomingAssignments(within time.Duration) (int64, error) {
	now := time.Now()
	var n int64
	err := r.DB.Model(&models.Assignment{}).
		Where("status = ? AND start_time > ? AND start_time <= ?", "scheduled", now, now.Add(within)).
		Count(&n).Error
	return n, err
}

func (r *Repos) CountAssets() (int64, error) {
	var n int64
	err := r.DB.Model(&models.Asset{}).Where("status NOT IN ?", []string{models.AssetRetired, models.AssetSold}).Count(&n).Error
	return n, err
}
