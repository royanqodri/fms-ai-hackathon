package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ChartConfig struct {
	Type    string       `json:"type"`
	Data    ChartData    `json:"data"`
	Options ChartOptions `json:"options"`
}

type ChartData struct {
	Labels   []string       `json:"labels"`
	Datasets []ChartDataset `json:"datasets"`
}

type ChartDataset struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
	Color string    `json:"backgroundColor"`
}

type ChartOptions struct {
	Title ChartTitle `json:"title"`
}

type ChartTitle struct {
	Display bool   `json:"display"`
	Text    string `json:"text"`
}

func GenerateChartImage(title string, labels []string, tempData, fuelData, oilData, rpmData []float64) ([]byte, error) {
	config := ChartConfig{
		Type: "bar",
		Data: ChartData{
			Labels: labels,
			Datasets: []ChartDataset{
				{Label: "Temperature (°C)", Data: tempData, Color: "rgba(255, 99, 132, 0.6)"},
				{Label: "Fuel Level (%)", Data: fuelData, Color: "rgba(54, 162, 235, 0.6)"},
				{Label: "Oil Pressure (PSI)", Data: oilData, Color: "rgba(255, 206, 86, 0.6)"},
				{Label: "Engine RPM (/10)", Data: rpmData, Color: "rgba(75, 192, 192, 0.6)"},
			},
		},
		Options: ChartOptions{
			Title: ChartTitle{
				Display: true,
				Text:    title,
			},
		},
	}

	jsonConfig, _ := json.Marshal(config)
	payload := map[string]string{"c": string(jsonConfig)}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post("https://quickchart.io/chart", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ QuickChart Error: Status %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("quickchart error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}
