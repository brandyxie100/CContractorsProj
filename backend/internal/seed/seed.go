package seed

import (
	"encoding/json"
	"log"
	"time"

	"github.com/clementscontractors/equipment/internal/authutil"
	"github.com/clementscontractors/equipment/internal/config"
	"github.com/clementscontractors/equipment/internal/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Run seeds admin user, sample locations, assets, and a project when empty.
func Run(db *gorm.DB, cfg config.Config) error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	log.Println("seeding database…")

	hash, err := authutil.HashPassword(cfg.SeedAdminPassword)
	if err != nil {
		return err
	}
	admin := models.User{
		ID: uuid.New(), Name: "Admin", Email: cfg.SeedAdminEmail,
		PasswordHash: hash, Role: models.RoleAdmin,
	}
	viewerHash, _ := authutil.HashPassword("viewer123")
	viewer := models.User{
		ID: uuid.New(), Name: "Viewer", Email: "viewer@clements.local",
		PasswordHash: viewerHash, Role: models.RoleViewer,
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	if err := db.Create(&viewer).Error; err != nil {
		return err
	}

	lat1, lng1 := -36.8485, 174.7633
	lat2, lng2 := -36.9100, 174.6900
	depot := models.Location{
		ID: uuid.New(), Name: "Main Depot", Type: "depot",
		Address: "12 Quarry Rd, Auckland", Lat: &lat1, Lng: &lng1,
	}
	yard := models.Location{
		ID: uuid.New(), Name: "West Yard", Type: "depot",
		Address: "88 Industrial Ave", Lat: &lat2, Lng: &lng2,
	}
	if err := db.Create(&depot).Error; err != nil {
		return err
	}
	if err := db.Create(&yard).Error; err != nil {
		return err
	}

	specs := func(weight float64) datatypes.JSON {
		b, _ := json.Marshal(map[string]interface{}{"weight_t": weight})
		return datatypes.JSON(b)
	}
	assets := []models.Asset{
		{ID: uuid.New(), AssetCode: "EXC-007", Name: "Excavator 20t", Model: "CAT 320", Type: "excavator", Category: "heavy_equipment", Specs: specs(20), Status: models.AssetAvailable, CurrentLocationID: &depot.ID},
		{ID: uuid.New(), AssetCode: "EXC-003", Name: "Excavator 25t", Model: "Komatsu PC210", Type: "excavator", Category: "heavy_equipment", Specs: specs(25), Status: models.AssetAvailable, CurrentLocationID: &yard.ID},
		{ID: uuid.New(), AssetCode: "EXC-005", Name: "Excavator 15t", Model: "Hitachi ZX135", Type: "excavator", Category: "heavy_equipment", Specs: specs(15), Status: models.AssetAvailable, CurrentLocationID: &depot.ID},
		{ID: uuid.New(), AssetCode: "TRK-012", Name: "Tipper Truck", Model: "Isuzu FVR", Type: "truck", Category: "trucks", Specs: specs(0), Status: models.AssetAvailable, CurrentLocationID: &depot.ID},
	}
	for i := range assets {
		if err := db.Create(&assets[i]).Error; err != nil {
			return err
		}
	}

	project := models.Project{
		ID: uuid.New(), Name: "Smith Street Earthworks", Address: "Smith Street",
		LocationID: &depot.ID, StartDate: time.Now(), EndDate: time.Now().AddDate(0, 1, 0),
		Status: "active", ProjectManagerID: &admin.ID, Description: "Driveway and excavation",
	}
	return db.Create(&project).Error
}
