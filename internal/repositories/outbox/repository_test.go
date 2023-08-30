package outbox_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	cnf "github.com/frutonanny/slug-service/internal/config"
	"github.com/frutonanny/slug-service/internal/database"
	"github.com/frutonanny/slug-service/internal/repositories/outbox"
)

const fileConfig = "../../../config/config.local.json"

var config = cnf.Must(fileConfig)

type RepositorySuite struct {
	database.DBSuite

	repository *outbox.Repository
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &RepositorySuite{DBSuite: database.NewDBSuite(config.DB.DSN)})
}

func (s *RepositorySuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	s.repository = outbox.New(s.DB)
}

func (s *RepositorySuite) TestRepository_FindJob() {
	// TODO
}

func (s *RepositorySuite) TestRepository_ReserveJob() {
	// TODO
}

func (s *RepositorySuite) TestRepository_CreateJob() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	var (
		name = "123"
		data = "data"
	)

	// Action.
	err := s.repository.CreateJob(ctx, name, data)

	// Assert.
	s.Require().NoError(err)

	const query = `select exists(select 1 from "outbox" where "name" = $1 and "data" = $2);`

	row := s.DB.QueryRow(ctx, query, name, data)
	s.Require().NoError(row.Err())

	var exist bool

	err = row.Scan(&exist)
	s.Require().NoError(err)
	s.Assert().True(exist)
}

func (s *RepositorySuite) TestRepository_DeleteJob() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	var (
		name = "123"
		data = "data"
	)

	// Action.
	err := s.repository.CreateJob(ctx, name, data)

	// Assert.
	s.Require().NoError(err)

	const query = `select exists(select 1 from "outbox" where "name" = $1 and "data" = $2);`

	row := s.DB.QueryRow(ctx, query, name, data)
	s.Require().NoError(row.Err())

	var exist bool

	err = row.Scan(&exist)
	s.Require().NoError(err)
	s.Assert().True(exist)
}
