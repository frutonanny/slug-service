package report

import (
	"context"
	"fmt"
	"time"

	"github.com/frutonanny/wallet-service/internal/postgres"
)

const (
	PeriodLayout = "2006-01"
)

type Repository struct {
	db postgres.Database
}

func New(db postgres.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) AddRecord(ctx context.Context, serviceID, amount int64, period time.Time) error {
	query := `insert into report("period", service_id, total_revenue)
values ($1, $2, $3)
on conflict (period, service_id) do update set total_revenue = report.total_revenue + $3;`

	_, err := r.db.ExecContext(ctx, query, period.Format(PeriodLayout), serviceID, amount)
	if err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}

func (r *Repository) GetReport(ctx context.Context, period string) ([]Service, error) {
	query := `select service_id, total_revenue
from report
where period = $1;`

	rows, err := r.db.QueryContext(ctx, query, period)
	if err != nil {
		return nil, fmt.Errorf("exec query: %v", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var result []Service

	for rows.Next() {
		s := Service{}

		if err := rows.Scan(&s.ServiceID, &s.TotalRevenue); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, s)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}
