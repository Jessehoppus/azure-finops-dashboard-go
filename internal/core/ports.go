package core

import (
	"context"
	"time"
)

// Porta (interface) que define o serviço de custos
type CostsService interface {
	CostByDimension(ctx context.Context, scope string, from, to time.Time, dimension, granularity string) (rows [][]any, headers []string, err error)
}
