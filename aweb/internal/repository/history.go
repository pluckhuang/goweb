package repository

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/internal/domain"
)

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
