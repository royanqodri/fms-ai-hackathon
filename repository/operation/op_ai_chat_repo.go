package operation

import (
	"errors"

	"github.com/gin-gonic/gin"
	dbConn "github.com/royanqodri/fms-ai-hackathon/database"
	"gorm.io/gorm"
)

type AiRepository interface {
	GetTruckContext(ctx *gin.Context, tx *gorm.DB, plateNumber string) (map[string]interface{}, error)
	GetRecentAlerts(ctx *gin.Context, tx *gorm.DB, limit int) ([]map[string]interface{}, error)
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
	var result map[string]interface{}

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
			CAST(tel.recorded_at AS VARCHAR) AS last_update
		FROM master.mst_trucks t
		LEFT JOIN master.mst_customers c ON t.customer_id = c.id
		LEFT JOIN operations.trx_telemetries tel ON t.id = tel.truck_id
		WHERE t.plate_number ILIKE ?
		ORDER BY tel.recorded_at DESC
		LIMIT 1
	`

	row := repo.getTx(tx).Raw(query, "%"+plateNumber+"%").Row()

	err := row.Scan(
		new(string), new(string), new(string), new(string),
		new(float64), new(float64), new(float64), new(float64), new(string),
	)

	type TruckContext struct {
		PlateNumber string  `json:"plate_number"`
		CargoType   string  `json:"cargo_type"`
		Status      string  `json:"status"`
		CompanyName string  `json:"company_name"`
		Temperature float64 `json:"temperature"`
		WeightLoad  float64 `json:"weight_load"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		LastUpdate  string  `json:"last_update"`
	}

	var data TruckContext
	err = repo.getTx(tx).Raw(query, "%"+plateNumber+"%").Scan(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	result = map[string]interface{}{
		"plate_number": data.PlateNumber,
		"cargo_type":   data.CargoType,
		"status":       data.Status,
		"company_name": data.CompanyName,
		"temperature":  data.Temperature,
		"weight_load":  data.WeightLoad,
		"location":     map[string]float64{"lat": data.Latitude, "lng": data.Longitude},
		"last_update":  data.LastUpdate,
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
