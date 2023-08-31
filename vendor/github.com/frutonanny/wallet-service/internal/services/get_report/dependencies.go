package get_report

import (
	"github.com/frutonanny/wallet-service/internal/postgres"
	repoReport "github.com/frutonanny/wallet-service/internal/repositories/report"
)

type dependenciesImpl struct{}

func (b *dependenciesImpl) NewRepository(db postgres.Database) Repository {
	return repoReport.New(db)
}
