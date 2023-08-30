package database

import (
	"context"

	"github.com/stretchr/testify/suite"
)

type DBSuite struct {
	suite.Suite

	DB  *DB
	dsn string
}

func NewDBSuite(dsn string) DBSuite {
	return DBSuite{
		dsn: dsn,
	}
}

func (s *DBSuite) SetupSuite() {
	var err error

	s.DB, err = New(s.dsn)
	s.Require().NoError(err)
}

func (s *DBSuite) TearDownSuite() {
	s.Assert().NoError(s.DB.Close())
}

func (s *DBSuite) Context() (context.Context, func()) {
	ctx := context.Background()

	tx, err := s.DB.DB.BeginTx(ctx, nil)
	s.Require().NoError(err)

	cancel := func() { s.Assert().NoError(tx.Rollback()) }

	return NewTxContext(ctx, tx), cancel
}
