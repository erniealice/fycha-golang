package fycha

import (
	"context"
	"time"

	expreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/expenditure_report"
	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
	agingpb    "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/receivables_aging"
	payagingpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/payables_aging"
	revreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/revenue_report"
	collsumpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/collection_summary"
	disbreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/disbursement_report"
	suppstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
)

// DataSource provides data access for fycha report views.
type DataSource interface {
	GetGrossProfitReport(ctx context.Context, req *reportpb.GrossProfitReportRequest) (*reportpb.GrossProfitReportResponse, error)
	GetRevenueReport(ctx context.Context, req *revreportpb.RevenueReportRequest) (*revreportpb.RevenueReportResponse, error)
	GetExpenditureReport(ctx context.Context, req *expreportpb.ExpenditureReportRequest) (*expreportpb.ExpenditureReportResponse, error)
	GetDisbursementReport(ctx context.Context, req *disbreportpb.DisbursementReportRequest) (*disbreportpb.DisbursementReportResponse, error)
	GetReceivablesAgingReport(ctx context.Context, req *agingpb.ReceivablesAgingRequest) (*agingpb.ReceivablesAgingResponse, error)
	GetPayablesAgingReport(ctx context.Context, req *payagingpb.PayablesAgingRequest) (*payagingpb.PayablesAgingResponse, error)
	GetCollectionSummaryReport(ctx context.Context, req *collsumpb.CollectionSummaryRequest) (*collsumpb.CollectionSummaryResponse, error)
	GetSupplierStatement(ctx context.Context, req *suppstmtpb.SupplierStatementRequest) (*suppstmtpb.SupplierStatementResponse, error)
	GetSupplierBalances(ctx context.Context) (map[string]int64, error)
	ListRevenue(ctx context.Context, start, end *time.Time) ([]map[string]any, error)
	ListExpenses(ctx context.Context, start, end *time.Time) ([]map[string]any, error)
}
