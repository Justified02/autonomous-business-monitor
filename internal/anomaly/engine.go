package anomaly

import (
	"context"
	"fmt"
	"math"

	"github.com/Justified02/abm/internal/storage/db"
)

type AnomalyResult struct {
	Source       string
	Metric       string
	CurrentValue float64
	Average      float64
	DeltaPct     float64
	IsAnomaly    bool
}

type Engine struct {
	db *db.Queries
}

func NewEngine(db *db.Queries) *Engine {
	newEngine := &Engine{
		db: db,
	}

	return newEngine
}

func (e *Engine) Analyze(ctx context.Context, source string, todayRevenue float64) (AnomalyResult, error) {
	data, err := e.db.GetLastSevenDays(ctx, source)
	if err != nil {
		return AnomalyResult{}, fmt.Errorf("failed to get data: %w", err)
	}

	if len(data) < 2 {
		return AnomalyResult{IsAnomaly: false}, nil
	}

	var sum float64

	for _, row := range data {
		revenue, _ := row.Revenue.Float64Value()
		sum += revenue.Float64
	}

	// calculate average
	aveRevenue := sum / float64(len(data))

	// calculate delta percentage
	deltaPc := ((aveRevenue - todayRevenue) / aveRevenue) * 100

	// build and return anomalyResult
	isAnomaly := math.Abs(deltaPc) > 20

	return AnomalyResult{
		Source:       source,
		Metric:       "revenue",
		CurrentValue: todayRevenue,
		Average:      aveRevenue,
		DeltaPct:     deltaPc,
		IsAnomaly:    isAnomaly,
	}, nil
}
