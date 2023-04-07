package helpers

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	mongodbTarget = "mongo-exporter-prometheus-mongodb-exporter"
)

func QueryPrometheus(api v1.API, metric string) (float64, error) {
	query := fmt.Sprintf(`avg by (instance) (
        rate(%s{job="%s"}[5m])
    )`, metric, mongodbTarget)

	value, _, err := api.Query(context.Background(), query, time.Now())
	if err != nil {
		return 0, err
	}

	// fmt.Print("value is",value)

	switch value.Type() {
	case model.ValVector:
		samples := value.(model.Vector)
		if len(samples) > 0 {
			return float64(samples[0].Value), nil
		} else {
			return 0, fmt.Errorf("vector has no samples")
		}
	case model.ValScalar:
		return float64(value.(*model.Scalar).Value), nil
	case model.ValString:
		return 0, fmt.Errorf("unexpected metric value type: %s", value.Type())
	default:
		return 0, fmt.Errorf("unknown metric value type: %s", value.Type())
	}
}
