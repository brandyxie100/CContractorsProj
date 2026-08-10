package service

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/clementscontractors/equipment/internal/authutil"
	"github.com/clementscontractors/equipment/internal/dto"
	"github.com/clementscontractors/equipment/internal/models"
	"github.com/clementscontractors/equipment/internal/repository"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Services holds application services.
type Services struct {
	Repos       *repository.Repos
	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// New creates Services.
func New(repos *repository.Repos, jwtSecret string, accessTTL, refreshTTL time.Duration) *Services {
	return &Services{Repos: repos, JWTSecret: jwtSecret, AccessTTL: accessTTL, RefreshTTL: refreshTTL}
}

// --- Auth ---

// Login authenticates and returns access token + refresh plain + user.
func (s *Services) Login(email, password string) (access string, refreshPlain string, user *models.User, err error) {
	u, err := s.Repos.FindUserByEmail(email)
	if err != nil {
		return "", "", nil, ErrUnauthorized
	}
	if !authutil.CheckPassword(u.PasswordHash, password) {
		return "", "", nil, ErrUnauthorized
	}
	access, err = authutil.IssueAccessToken(s.JWTSecret, s.AccessTTL, u.ID, u.Email, u.Role)
	if err != nil {
		return "", "", nil, err
	}
	plain, hash, err := authutil.NewRefreshToken()
	if err != nil {
		return "", "", nil, err
	}
	rt := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    u.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.RefreshTTL),
	}
	if err := s.Repos.SaveRefreshToken(rt); err != nil {
		return "", "", nil, err
	}
	return access, plain, u, nil
}

// Refresh exchanges a refresh token for a new access token.
func (s *Services) Refresh(refreshPlain string) (string, error) {
	hash := authutil.HashRefreshToken(refreshPlain)
	rt, err := s.Repos.FindRefreshToken(hash)
	if err != nil {
		return "", ErrUnauthorized
	}
	u, err := s.Repos.FindUserByID(rt.UserID)
	if err != nil {
		return "", ErrUnauthorized
	}
	return authutil.IssueAccessToken(s.JWTSecret, s.AccessTTL, u.ID, u.Email, u.Role)
}

// Logout invalidates a refresh token.
func (s *Services) Logout(refreshPlain string) error {
	if refreshPlain == "" {
		return nil
	}
	return s.Repos.DeleteRefreshToken(authutil.HashRefreshToken(refreshPlain))
}

// --- Availability ---

// IsAvailable checks assignment and maintenance overlaps.
func (s *Services) IsAvailable(assetID uuid.UUID, start, end time.Time) (bool, error) {
	asset, err := s.Repos.FindAsset(assetID)
	if err != nil {
		return false, err
	}
	if asset.Status == models.AssetRetired || asset.Status == models.AssetSold {
		return false, nil
	}
	as, err := s.Repos.FindOverlappingAssignments(assetID, start, end, nil)
	if err != nil {
		return false, err
	}
	if len(as) > 0 {
		return false, nil
	}
	ms, err := s.Repos.FindOverlappingMaintenance(assetID, start, end)
	if err != nil {
		return false, err
	}
	return len(ms) == 0, nil
}

// NextAvailableAt returns when the asset becomes free.
func (s *Services) NextAvailableAt(assetID uuid.UUID) (time.Time, error) {
	now := time.Now().UTC()
	as, err := s.Repos.FindOverlappingAssignments(assetID, now, now.Add(365*24*time.Hour), nil)
	if err != nil {
		return time.Time{}, err
	}
	ms, err := s.Repos.FindOverlappingMaintenance(assetID, now, now.Add(365*24*time.Hour))
	if err != nil {
		return time.Time{}, err
	}

	var currentEnd *time.Time
	for _, a := range as {
		if !a.StartTime.After(now) && a.EndTime.After(now) {
			end := a.EndTime
			currentEnd = &end
			break
		}
	}
	for _, m := range ms {
		if !m.ScheduledStart.After(now) && m.ScheduledEnd.After(now) {
			end := m.ScheduledEnd
			if currentEnd == nil || end.After(*currentEnd) {
				currentEnd = &end
			}
		}
	}
	if currentEnd == nil {
		return now, nil
	}
	after := *currentEnd
	for _, m := range ms {
		if !m.ScheduledStart.Before(after) && m.ScheduledStart.Before(after.Add(time.Hour)) {
			if m.ScheduledEnd.After(after) {
				after = m.ScheduledEnd
			}
		}
	}
	return after, nil
}

// CheckAssignmentConflicts returns ErrConflict if overlaps exist.
func (s *Services) CheckAssignmentConflicts(assetID uuid.UUID, start, end time.Time, excludeID *uuid.UUID) error {
	asset, err := s.Repos.FindAsset(assetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if asset.Status == models.AssetRetired || asset.Status == models.AssetSold {
		return ErrConflict
	}
	as, err := s.Repos.FindOverlappingAssignments(assetID, start, end, excludeID)
	if err != nil {
		return err
	}
	if len(as) > 0 {
		return ErrConflict
	}
	ms, err := s.Repos.FindOverlappingMaintenance(assetID, start, end)
	if err != nil {
		return err
	}
	if len(ms) > 0 {
		return ErrConflict
	}
	return nil
}

// CreateAssignment validates conflicts then persists.
func (s *Services) CreateAssignment(req dto.AssignmentRequest, assignedBy uuid.UUID) (*models.Assignment, error) {
	if !req.EndTime.After(req.StartTime) {
		return nil, ErrValidation
	}
	if err := s.CheckAssignmentConflicts(req.AssetID, req.StartTime, req.EndTime, nil); err != nil {
		return nil, err
	}
	if _, err := s.Repos.FindProject(req.ProjectID); err != nil {
		return nil, ErrNotFound
	}
	a := &models.Assignment{
		ID:         uuid.New(),
		AssetID:    req.AssetID,
		ProjectID:  req.ProjectID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Status:     "scheduled",
		AssignedBy: assignedBy,
		Notes:      req.Notes,
	}
	if err := s.Repos.CreateAssignment(a); err != nil {
		return nil, err
	}
	_ = s.Repos.DB.Model(&models.Asset{}).Where("id = ?", req.AssetID).Update("status", models.AssetReserved)
	return s.Repos.FindAssignment(a.ID)
}

// UpdateAssignment updates times/notes with conflict check.
func (s *Services) UpdateAssignment(id uuid.UUID, req dto.AssignmentRequest) (*models.Assignment, error) {
	a, err := s.Repos.FindAssignment(id)
	if err != nil {
		return nil, ErrNotFound
	}
	if !req.EndTime.After(req.StartTime) {
		return nil, ErrValidation
	}
	if err := s.CheckAssignmentConflicts(req.AssetID, req.StartTime, req.EndTime, &id); err != nil {
		return nil, err
	}
	a.AssetID = req.AssetID
	a.ProjectID = req.ProjectID
	a.StartTime = req.StartTime
	a.EndTime = req.EndTime
	a.Notes = req.Notes
	if err := s.Repos.UpdateAssignment(a); err != nil {
		return nil, err
	}
	return s.Repos.FindAssignment(id)
}

// CancelAssignment marks assignment cancelled.
func (s *Services) CancelAssignment(id uuid.UUID) error {
	a, err := s.Repos.FindAssignment(id)
	if err != nil {
		return ErrNotFound
	}
	a.Status = "cancelled"
	return s.Repos.UpdateAssignment(a)
}

// --- Dispatch ---

// DispatchSearch ranks equipment for a request.
func (s *Services) DispatchSearch(req dto.DispatchSearchRequest) ([]dto.DispatchResultItem, error) {
	types := []string{req.EquipmentType}
	if req.IncludeAlternatives {
		// broaden: any non-retired assets of same category heuristics — keep same type family only for MVP
		types = []string{req.EquipmentType}
	}
	assets, err := s.Repos.ListAssetsByTypes(types)
	if err != nil {
		return nil, err
	}
	// If alternatives and no exact, still search same type; also include all excavators when type excavator etc.
	if req.IncludeAlternatives {
		all, err := s.Repos.ListAssets(repository.AssetFilter{})
		if err != nil {
			return nil, err
		}
		seen := map[uuid.UUID]struct{}{}
		for _, a := range assets {
			seen[a.ID] = struct{}{}
		}
		for _, a := range all {
			if _, ok := seen[a.ID]; ok {
				continue
			}
			if a.Status == models.AssetRetired || a.Status == models.AssetSold {
				continue
			}
			// alternative: same category
			if a.Category == "heavy_equipment" && req.EquipmentType == "excavator" && a.Type != "excavator" {
				assets = append(assets, a)
			}
			if a.Type != req.EquipmentType && a.Type == "excavator" && req.EquipmentType == "excavator" {
				continue
			}
			if a.Type != req.EquipmentType && req.EquipmentType == "excavator" && a.Type == "excavator" {
				assets = append(assets, a)
			}
			// include larger/same family: any excavator already in list; add trucks if searching truck
			if a.Type != req.EquipmentType && (a.Type == "excavator" || a.Type == "truck" || a.Type == "roller") {
				// only add if include alternatives and type differs but category heavy
				if a.Category == categoryForType(req.EquipmentType) {
					assets = append(assets, a)
				}
			}
		}
	}

	jobLat, jobLng, err := s.resolveJobCoords(req)
	if err != nil {
		return nil, err
	}

	results := make([]dto.DispatchResultItem, 0, len(assets))
	seen := map[uuid.UUID]struct{}{}
	for _, asset := range assets {
		if _, ok := seen[asset.ID]; ok {
			continue
		}
		seen[asset.ID] = struct{}{}

		available, err := s.IsAvailable(asset.ID, req.StartTime, req.EndTime)
		if err != nil {
			return nil, err
		}
		availableFrom := req.StartTime
		if !available {
			availableFrom, err = s.NextAvailableAt(asset.ID)
			if err != nil {
				return nil, err
			}
		}
		availScore := AvailabilityScore(available, availableFrom, req.StartTime)

		var dist *float64
		var transport *float64
		prox := 0.5
		if jobLat != nil && jobLng != nil && asset.CurrentLocation != nil && asset.CurrentLocation.Lat != nil && asset.CurrentLocation.Lng != nil {
			d := HaversineKm(*jobLat, *jobLng, *asset.CurrentLocation.Lat, *asset.CurrentLocation.Lng)
			dist = &d
			tm := TransportMinutes(d)
			transport = &tm
			prox = ProximityScore(d)
		}

		weight := weightFromSpecs(asset.Specs)
		suit, under := SpecsScore(weight, req.MinWeightT)
		warnings := []string{}
		if under {
			warnings = append(warnings, "smaller than requested capacity")
		}
		isAlt := asset.Type != req.EquipmentType
		rank := RankScore(availScore, prox, suit)
		if isAlt {
			rank *= 0.9
		}
		results = append(results, dto.DispatchResultItem{
			Asset:            asset,
			Available:        available,
			AvailableFrom:    availableFrom,
			DistanceKm:       dist,
			TransportMinutes: transport,
			SuitabilityScore: suit,
			RankScore:        rank,
			IsAlternative:    isAlt,
			Warnings:         warnings,
		})
	}

	sort.SliceStable(results, func(i, j int) bool {
		// exact type before alternatives when scores close
		if results[i].IsAlternative != results[j].IsAlternative {
			return !results[i].IsAlternative
		}
		return results[i].RankScore > results[j].RankScore
	})
	return results, nil
}

func categoryForType(t string) string {
	switch t {
	case "truck":
		return "trucks"
	case "excavator", "roller":
		return "heavy_equipment"
	default:
		return "heavy_equipment"
	}
}

func (s *Services) resolveJobCoords(req dto.DispatchSearchRequest) (*float64, *float64, error) {
	if req.JobLat != nil && req.JobLng != nil {
		return req.JobLat, req.JobLng, nil
	}
	if req.ProjectID != nil {
		p, err := s.Repos.FindProject(*req.ProjectID)
		if err != nil {
			return nil, nil, ErrNotFound
		}
		if p.Location != nil && p.Location.Lat != nil && p.Location.Lng != nil {
			return p.Location.Lat, p.Location.Lng, nil
		}
	}
	return nil, nil, nil
}

func weightFromSpecs(specs datatypes.JSON) *float64 {
	if len(specs) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(specs, &m); err != nil {
		return nil
	}
	if v, ok := m["weight_t"]; ok {
		switch n := v.(type) {
		case float64:
			return &n
		case int:
			f := float64(n)
			return &f
		}
	}
	return nil
}

// DashboardSummary builds KPI payload.
func (s *Services) DashboardSummary() (*dto.DashboardSummary, error) {
	byStatus, err := s.Repos.CountAssetsByStatus()
	if err != nil {
		return nil, err
	}
	overdue, err := s.Repos.CountOverdueMaintenance()
	if err != nil {
		return nil, err
	}
	active, err := s.Repos.CountActiveAssignments()
	if err != nil {
		return nil, err
	}
	upcoming, err := s.Repos.CountUpcomingAssignments(7 * 24 * time.Hour)
	if err != nil {
		return nil, err
	}
	total, err := s.Repos.CountAssets()
	if err != nil {
		return nil, err
	}
	util := 0.0
	if total > 0 {
		assigned := byStatus[models.AssetAssigned] + byStatus[models.AssetReserved]
		util = float64(assigned) / float64(total) * 100
	}
	return &dto.DashboardSummary{
		FleetByStatus:       byStatus,
		OverdueMaintenance:  overdue,
		ActiveAssignments:   active,
		UtilizationPercent:  util,
		UpcomingAssignments: upcoming,
	}, nil
}

// CreateMaintenance creates a maintenance window and sets asset status when in progress / now.
func (s *Services) CreateMaintenance(req dto.MaintenanceRequest) (*models.MaintenanceRecord, error) {
	if !req.ScheduledEnd.After(req.ScheduledStart) {
		return nil, ErrValidation
	}
	if _, err := s.Repos.FindAsset(req.AssetID); err != nil {
		return nil, ErrNotFound
	}
	status := req.Status
	if status == "" {
		status = "scheduled"
	}
	m := &models.MaintenanceRecord{
		ID:             uuid.New(),
		AssetID:        req.AssetID,
		Type:           req.Type,
		Description:    req.Description,
		ScheduledStart: req.ScheduledStart,
		ScheduledEnd:   req.ScheduledEnd,
		Status:         status,
		Cost:           req.Cost,
		PerformedBy:    req.PerformedBy,
		Notes:          req.Notes,
	}
	if err := s.Repos.CreateMaintenance(m); err != nil {
		return nil, err
	}
	now := time.Now()
	if status == "in_progress" || (status == "scheduled" && !m.ScheduledStart.After(now) && m.ScheduledEnd.After(now)) {
		_ = s.Repos.DB.Model(&models.Asset{}).Where("id = ?", req.AssetID).Update("status", models.AssetMaintenance)
	}
	return s.Repos.FindMaintenance(m.ID)
}

// ToUserResponse maps user model to DTO.
func ToUserResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role, Phone: u.Phone, AvatarURL: u.AvatarURL,
	}
}
