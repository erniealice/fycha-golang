package fycha

import (
	"context"
	"time"

	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
)

// DataSource provides data access for fycha report views.
type DataSource interface {
	GetGrossProfitReport(ctx context.Context, req *reportpb.GrossProfitReportRequest) (*reportpb.GrossProfitReportResponse, error)
	ListRevenue(ctx context.Context, start, end *time.Time) ([]map[string]any, error)
	ListExpenses(ctx context.Context, start, end *time.Time) ([]map[string]any, error)
}
