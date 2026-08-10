package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/clementscontractors/equipment/internal/authutil"
	"github.com/clementscontractors/equipment/internal/dto"
	"github.com/clementscontractors/equipment/internal/middleware"
	"github.com/clementscontractors/equipment/internal/models"
	"github.com/clementscontractors/equipment/internal/repository"
	"github.com/clementscontractors/equipment/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const refreshCookie = "refresh_token"

// API wires HTTP handlers.
type API struct {
	S *service.Services
	R *repository.Repos
}

// Register mounts all routes on r.
func (a *API) Register(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", a.Login)
		auth.POST("/refresh", a.Refresh)
		auth.POST("/logout", a.Logout)
		auth.GET("/me", middleware.Auth(a.S.JWTSecret), a.Me)
	}

	secured := v1.Group("")
	secured.Use(middleware.Auth(a.S.JWTSecret))
	{
		secured.GET("/users", middleware.RequireAdmin(), a.ListUsers)
		secured.POST("/users", middleware.RequireAdmin(), a.CreateUser)
		secured.GET("/users/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleManager), a.GetUser)
		secured.PATCH("/users/:id", middleware.RequireAdmin(), a.UpdateUser)
		secured.DELETE("/users/:id", middleware.RequireAdmin(), a.DeleteUser)

		secured.GET("/locations", a.ListLocations)
		secured.POST("/locations", middleware.RequireManagerPlus(), a.CreateLocation)
		secured.GET("/locations/:id", a.GetLocation)
		secured.PATCH("/locations/:id", middleware.RequireManagerPlus(), a.UpdateLocation)
		secured.DELETE("/locations/:id", middleware.RequireManagerPlus(), a.DeleteLocation)

		secured.GET("/assets", a.ListAssets)
		secured.POST("/assets", middleware.RequireManagerPlus(), a.CreateAsset)
		secured.GET("/assets/:id", a.GetAsset)
		secured.PATCH("/assets/:id", middleware.RequireManagerPlus(), a.UpdateAsset)
		secured.DELETE("/assets/:id", middleware.RequireManagerPlus(), a.DeleteAsset)

		secured.GET("/projects", a.ListProjects)
		secured.POST("/projects", middleware.RequireManagerPlus(), a.CreateProject)
		secured.GET("/projects/:id", a.GetProject)
		secured.PATCH("/projects/:id", middleware.RequireManagerPlus(), a.UpdateProject)
		secured.DELETE("/projects/:id", middleware.RequireManagerPlus(), a.DeleteProject)

		secured.GET("/assignments", a.ListAssignments)
		secured.POST("/assignments", middleware.RequireOperatorPlus(), a.CreateAssignment)
		secured.GET("/assignments/:id", a.GetAssignment)
		secured.PATCH("/assignments/:id", middleware.RequireOperatorPlus(), a.UpdateAssignment)
		secured.DELETE("/assignments/:id", middleware.RequireOperatorPlus(), a.CancelAssignment)

		secured.GET("/maintenance", a.ListMaintenance)
		secured.POST("/maintenance", middleware.RequireManagerPlus(), a.CreateMaintenance)
		secured.GET("/maintenance/:id", a.GetMaintenance)
		secured.PATCH("/maintenance/:id", middleware.RequireManagerPlus(), a.UpdateMaintenance)
		secured.DELETE("/maintenance/:id", middleware.RequireManagerPlus(), a.DeleteMaintenance)

		secured.GET("/dashboard/summary", a.DashboardSummary)
		secured.POST("/dispatch/search", a.DispatchSearch)
	}
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: dto.APIError{Code: "UNAUTHORIZED", Message: err.Error()}})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: dto.APIError{Code: "FORBIDDEN", Message: err.Error()}})
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: dto.APIError{Code: "NOT_FOUND", Message: err.Error()}})
	case errors.Is(err, service.ErrConflict):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: dto.APIError{Code: "CONFLICT", Message: "assignment or maintenance conflict"}})
	case errors.Is(err, service.ErrValidation):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: dto.APIError{Code: "INTERNAL", Message: err.Error()}})
	}
}

func (a *API) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	access, refresh, user, err := a.S.Login(req.Email, req.Password)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.SetCookie(refreshCookie, refresh, int(a.S.RefreshTTL.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, dto.AuthResponse{AccessToken: access, User: service.ToUserResponse(user)})
}

func (a *API) Refresh(c *gin.Context) {
	refresh, err := c.Cookie(refreshCookie)
	if err != nil || refresh == "" {
		writeErr(c, service.ErrUnauthorized)
		return
	}
	access, err := a.S.Refresh(refresh)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

func (a *API) Logout(c *gin.Context) {
	refresh, _ := c.Cookie(refreshCookie)
	_ = a.S.Logout(refresh)
	c.SetCookie(refreshCookie, "", -1, "/", "", false, true)
	c.Status(http.StatusNoContent)
}

func (a *API) Me(c *gin.Context) {
	uid := c.MustGet(middleware.ContextUserID).(uuid.UUID)
	u, err := a.R.FindUserByID(uid)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, service.ToUserResponse(u))
}

func (a *API) ListUsers(c *gin.Context) {
	list, err := a.R.ListUsers()
	if err != nil {
		writeErr(c, err)
		return
	}
	out := make([]dto.UserResponse, 0, len(list))
	for i := range list {
		out = append(out, service.ToUserResponse(&list[i]))
	}
	c.JSON(http.StatusOK, out)
}

func (a *API) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	hash, err := authutil.HashPassword(req.Password)
	if err != nil {
		writeErr(c, err)
		return
	}
	u := &models.User{ID: uuid.New(), Name: req.Name, Email: req.Email, PasswordHash: hash, Role: req.Role, Phone: req.Phone}
	if err := a.R.CreateUser(u); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, service.ToUserResponse(u))
}

func (a *API) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	u, err := a.R.FindUserByID(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, service.ToUserResponse(u))
}

func (a *API) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	u, err := a.R.FindUserByID(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Role != nil {
		u.Role = *req.Role
	}
	if req.Phone != nil {
		u.Phone = req.Phone
	}
	if err := a.R.UpdateUser(u); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, service.ToUserResponse(u))
}

func (a *API) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	if err := a.R.DeleteUser(id); err != nil {
		writeErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *API) ListLocations(c *gin.Context) {
	list, err := a.R.ListLocations()
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *API) CreateLocation(c *gin.Context) {
	var req dto.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	l := &models.Location{ID: uuid.New(), Name: req.Name, Type: req.Type, Address: req.Address, Lat: req.Lat, Lng: req.Lng, Description: req.Description}
	if err := a.R.CreateLocation(l); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, l)
}

func (a *API) GetLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	l, err := a.R.FindLocation(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, l)
}

func (a *API) UpdateLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	l, err := a.R.FindLocation(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	var req dto.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	l.Name, l.Type, l.Address, l.Lat, l.Lng, l.Description = req.Name, req.Type, req.Address, req.Lat, req.Lng, req.Description
	if err := a.R.UpdateLocation(l); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, l)
}

func (a *API) DeleteLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	if err := a.R.DeleteLocation(id); err != nil {
		writeErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *API) ListAssets(c *gin.Context) {
	f := repository.AssetFilter{Type: c.Query("type"), Status: c.Query("status"), Q: c.Query("q")}
	if loc := c.Query("location_id"); loc != "" {
		id, err := uuid.Parse(loc)
		if err == nil {
			f.LocationID = &id
		}
	}
	list, err := a.R.ListAssets(f)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func specsJSON(m map[string]interface{}) datatypes.JSON {
	if m == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, _ := json.Marshal(m)
	return datatypes.JSON(b)
}

func (a *API) CreateAsset(c *gin.Context) {
	var req dto.AssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	status := req.Status
	if status == "" {
		status = models.AssetAvailable
	}
	asset := &models.Asset{
		ID: uuid.New(), AssetCode: req.AssetCode, Name: req.Name, Model: req.Model, Type: req.Type,
		Category: req.Category, Specs: specsJSON(req.Specs), Status: status, CurrentLocationID: req.CurrentLocationID,
		PurchaseCost: req.PurchaseCost, HourlyRate: req.HourlyRate, PhotoURL: req.PhotoURL,
	}
	if err := a.R.CreateAsset(asset); err != nil {
		writeErr(c, err)
		return
	}
	out, _ := a.R.FindAsset(asset.ID)
	c.JSON(http.StatusCreated, out)
}

func (a *API) GetAsset(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	asset, err := a.R.FindAsset(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	next, _ := a.S.NextAvailableAt(id)
	c.JSON(http.StatusOK, gin.H{"asset": asset, "next_available_at": next})
}

func (a *API) UpdateAsset(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	asset, err := a.R.FindAsset(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	var req dto.AssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	asset.AssetCode, asset.Name, asset.Model, asset.Type = req.AssetCode, req.Name, req.Model, req.Type
	asset.Category, asset.Specs = req.Category, specsJSON(req.Specs)
	if req.Status != "" {
		asset.Status = req.Status
	}
	asset.CurrentLocationID, asset.PurchaseCost, asset.HourlyRate, asset.PhotoURL = req.CurrentLocationID, req.PurchaseCost, req.HourlyRate, req.PhotoURL
	if err := a.R.UpdateAsset(asset); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, asset)
}

func (a *API) DeleteAsset(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	asset, err := a.R.FindAsset(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	asset.Status = models.AssetRetired
	if err := a.R.UpdateAsset(asset); err != nil {
		writeErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseDate(s string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}

func (a *API) ListProjects(c *gin.Context) {
	list, err := a.R.ListProjects(c.Query("status"))
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *API) CreateProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	start, err := parseDate(req.StartDate)
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	end, err := parseDate(req.EndDate)
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	status := req.Status
	if status == "" {
		status = "planning"
	}
	p := &models.Project{
		ID: uuid.New(), Name: req.Name, Address: req.Address, LocationID: req.LocationID,
		StartDate: start, EndDate: end, Status: status, ProjectManagerID: req.ProjectManagerID, Description: req.Description,
	}
	if err := a.R.CreateProject(p); err != nil {
		writeErr(c, err)
		return
	}
	out, _ := a.R.FindProject(p.ID)
	c.JSON(http.StatusCreated, out)
}

func (a *API) GetProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	p, err := a.R.FindProject(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (a *API) UpdateProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	p, err := a.R.FindProject(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	start, err := parseDate(req.StartDate)
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	end, err := parseDate(req.EndDate)
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	p.Name, p.Address, p.LocationID = req.Name, req.Address, req.LocationID
	p.StartDate, p.EndDate = start, end
	if req.Status != "" {
		p.Status = req.Status
	}
	p.ProjectManagerID, p.Description = req.ProjectManagerID, req.Description
	if err := a.R.UpdateProject(p); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (a *API) DeleteProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	if err := a.R.DeleteProject(id); err != nil {
		writeErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *API) ListAssignments(c *gin.Context) {
	var from, to *time.Time
	if v := c.Query("from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			from = &t
		}
	}
	if v := c.Query("to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			to = &t
		}
	}
	var assetID, projectID *uuid.UUID
	if v := c.Query("asset_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			assetID = &id
		}
	}
	if v := c.Query("project_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			projectID = &id
		}
	}
	list, err := a.R.ListAssignments(from, to, assetID, projectID)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *API) CreateAssignment(c *gin.Context) {
	var req dto.AssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	uid := c.MustGet(middleware.ContextUserID).(uuid.UUID)
	aOut, err := a.S.CreateAssignment(req, uid)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, aOut)
}

func (a *API) GetAssignment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	as, err := a.R.FindAssignment(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, as)
}

func (a *API) UpdateAssignment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	var req dto.AssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	as, err := a.S.UpdateAssignment(id, req)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, as)
}

func (a *API) CancelAssignment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	if err := a.S.CancelAssignment(id); err != nil {
		writeErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *API) ListMaintenance(c *gin.Context) {
	var assetID *uuid.UUID
	if v := c.Query("asset_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			assetID = &id
		}
	}
	list, err := a.R.ListMaintenance(assetID, c.Query("status"))
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *API) CreateMaintenance(c *gin.Context) {
	var req dto.MaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	m, err := a.S.CreateMaintenance(req)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (a *API) GetMaintenance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	m, err := a.R.FindMaintenance(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (a *API) UpdateMaintenance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	m, err := a.R.FindMaintenance(id)
	if err != nil {
		writeErr(c, service.ErrNotFound)
		return
	}
	var req dto.MaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	m.Type, m.Description = req.Type, req.Description
	m.ScheduledStart, m.ScheduledEnd = req.ScheduledStart, req.ScheduledEnd
	if req.Status != "" {
		m.Status = req.Status
	}
	m.Cost, m.PerformedBy, m.Notes = req.Cost, req.PerformedBy, req.Notes
	if err := a.R.UpdateMaintenance(m); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (a *API) DeleteMaintenance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, service.ErrValidation)
		return
	}
	if err := a.R.DeleteMaintenance(id); err != nil {
		writeErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *API) DashboardSummary(c *gin.Context) {
	sum, err := a.S.DashboardSummary()
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, sum)
}

func (a *API) DispatchSearch(c *gin.Context) {
	var req dto.DispatchSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.APIError{Code: "VALIDATION", Message: err.Error()}})
		return
	}
	results, err := a.S.DispatchSearch(req)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

// Ensure gorm import used for ErrRecordNotFound mapping elsewhere.
var _ = gorm.ErrRecordNotFound
