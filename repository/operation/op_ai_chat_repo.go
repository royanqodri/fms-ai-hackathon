package operation

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	dbConn "github.com/royanqodri/fms-ai-hackathon/database"
	"gorm.io/gorm"
)

type AiRepository interface {
	GetTruckContext(ctx *gin.Context, tx *gorm.DB, plateNumber string) (map[string]interface{}, error)
	GetRecentAlerts(ctx *gin.Context, tx *gorm.DB, limit int) ([]map[string]interface{}, error)
	GetTemperatureHistory(ctx *gin.Context, tx *gorm.DB, plateNumber string, limit int) ([]map[string]interface{}, error)
}

type AiRepositoryImpl struct{}

func NewAiRepository() AiRepository {
	return &AiRepositoryImpl{}
}

func (repo AiRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dbConn.DBConn
	}
	return tx
}

func (repo AiRepositoryImpl) GetTruckContext(ctx *gin.Context, tx *gorm.DB, plateNumber string) (map[string]interface{}, error) {

	upper := strings.ToUpper(plateNumber)
	re := regexp.MustCompile(`([A-Z]{1,2})[\s\.\-]*([0-9]{1,4})[\s\.\-]*([A-Z]{1,3})`)
	matches := re.FindStringSubmatch(upper)

	plateWithSpace := plateNumber
	plateWithoutSpace := plateNumber

	if len(matches) > 0 {
		plateWithSpace = fmt.Sprintf("%s %s %s", matches[1], matches[2], matches[3])
		plateWithoutSpace = matches[1] + matches[2] + matches[3]
	}

	query := `
		SELECT 
			t.plate_number, 
			t.cargo_type, 
			t.status, 
			c.name AS company_name,
			tel.temperature, 
			tel.weight_load, 
			tel.latitude, 
			tel.longitude, 
			tel.fuel_level,   
			tel.oil_pressure, 
			tel.engine_rpm,   
			CAST(tel.recorded_at AS VARCHAR) AS last_update
		FROM master.mst_trucks t
		LEFT JOIN master.mst_customers c ON t.customer_id = c.id
		LEFT JOIN operations.trx_telemetries tel ON t.id = tel.truck_id
		WHERE t.plate_number ILIKE ? OR t.plate_number ILIKE ?
		ORDER BY tel.recorded_at DESC
		LIMIT 1
	`

	type TruckContext struct {
		PlateNumber string          `json:"plate_number"`
		CargoType   string          `json:"cargo_type"`
		Status      string          `json:"status"`
		CompanyName string          `json:"company_name"`
		Temperature sql.NullFloat64 `json:"temperature"`
		WeightLoad  sql.NullFloat64 `json:"weight_load"`
		Latitude    sql.NullFloat64 `json:"latitude"`
		Longitude   sql.NullFloat64 `json:"longitude"`
		FuelLevel   sql.NullFloat64 `json:"fuel_level"`
		OilPressure sql.NullFloat64 `json:"oil_pressure"`
		EngineRpm   sql.NullFloat64 `json:"engine_rpm"`
		LastUpdate  sql.NullString  `json:"last_update"`
	}

	var data TruckContext

	err := repo.getTx(tx).Raw(query, "%"+plateWithSpace+"%", "%"+plateWithoutSpace+"%").Scan(&data).Error

	if err != nil {
		log.Println("❌ Error GetTruckContext:", err)
		return nil, err
	}

	if data.PlateNumber == "" {
		log.Println("⚠️ Truk not found in database for plat:", plateNumber)
		return nil, nil
	}

	log.Println("✅ Truk success found in database:", data.PlateNumber)

	result := map[string]interface{}{
		"plate_number": data.PlateNumber,
		"cargo_type":   data.CargoType,
		"status":       data.Status,
		"company_name": data.CompanyName,
		"temperature":  data.Temperature.Float64,
		"weight_load":  data.WeightLoad.Float64,
		"location":     map[string]float64{"lat": data.Latitude.Float64, "lng": data.Longitude.Float64},
		"fuel_level":   data.FuelLevel.Float64,
		"oil_pressure": data.OilPressure.Float64,
		"engine_rpm":   data.EngineRpm.Float64,
		"last_update":  data.LastUpdate.String,
	}

	return result, nil
}

func (repo AiRepositoryImpl) GetRecentAlerts(ctx *gin.Context, tx *gorm.DB, limit int) ([]map[string]interface{}, error) {
	type AlertContext struct {
		PlateNumber string `json:"plate_number"`
		AlertType   string `json:"alert_type"`
		AiAnalysis  string `json:"ai_analysis"`
		CreatedAt   string `json:"created_at"`
	}

	var alerts []AlertContext

	query := `
		SELECT 
			t.plate_number, 
			a.alert_type, 
			a.ai_analysis, 
			CAST(a.created_at AS VARCHAR) AS created_at
		FROM operations.log_alerts a
		JOIN master.mst_trucks t ON a.truck_id = t.id
		WHERE a.is_resolved = false
		ORDER BY a.created_at DESC
		LIMIT ?
	`

	err := repo.getTx(tx).Raw(query, limit).Scan(&alerts).Error
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, a := range alerts {
		results = append(results, map[string]interface{}{
			"plate_number": a.PlateNumber,
			"alert_type":   a.AlertType,
			"ai_analysis":  a.AiAnalysis,
			"created_at":   a.CreatedAt,
		})
	}

	return results, nil
}

func (repo AiRepositoryImpl) GetTemperatureHistory(ctx *gin.Context, tx *gorm.DB, plateNumber string, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			tel.temperature, 
			tel.fuel_level,
			tel.oil_pressure,
			tel.engine_rpm,
			CAST(tel.recorded_at AS VARCHAR) AS recorded_at
		FROM operations.trx_telemetries tel
		JOIN master.mst_trucks t ON tel.truck_id = t.id
		WHERE t.plate_number ILIKE ?
		ORDER BY tel.recorded_at DESC
		LIMIT ?
	`
	type TempHistory struct {
		Temperature float64 `json:"temperature"`
		FuelLevel   float64 `json:"fuel_level"`
		OilPressure float64 `json:"oil_pressure"`
		EngineRpm   float64 `json:"engine_rpm"`
		RecordedAt  string  `json:"recorded_at"`
	}

	var history []TempHistory
	err := repo.getTx(tx).Raw(query, "%"+plateNumber+"%", limit).Scan(&history).Error
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for i := len(history) - 1; i >= 0; i-- {
		h := history[i]
		results = append(results, map[string]interface{}{
			"temperature":  h.Temperature,
			"fuel_level":   h.FuelLevel,
			"oil_pressure": h.OilPressure,
			"engine_rpm":   h.EngineRpm,
			"recorded_at":  h.RecordedAt,
		})
	}
	return results, nil
}
