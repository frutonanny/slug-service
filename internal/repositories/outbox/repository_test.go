package outbox_test

import (
	"testing"
	"time"

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
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	var (
		name = "123"
		data = "data"
	)

	err := s.repository.CreateJob(ctx, name, data)
	s.Require().NoError(err)

	const queryInsert = `insert into "outbox" (name, data) values($1, $2) returning id;`
	rowInsert := s.DB.QueryRow(ctx, queryInsert, name, data)
	s.Require().NoError(rowInsert.Err())

	var id int64
	err = rowInsert.Scan(&id)
	s.Require().NoError(rowInsert.Err())

	// Action.
	actualJob, err := s.repository.FindJob(ctx)

	// Assert.
	s.Require().NoError(err)
	s.Assert().EqualValues(name, actualJob.Name)
	s.Assert().EqualValues(data, actualJob.Data)
}

func (s *RepositorySuite) TestRepository_ReserveJob() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	var (
		name = "123"
		data = "data"
	)

	err := s.repository.CreateJob(ctx, name, data)
	s.Require().NoError(err)

	const queryInsert = `insert into "outbox" (name, data) values($1, $2) returning id;`
	rowInsert := s.DB.QueryRow(ctx, queryInsert, name, data)
	s.Require().NoError(rowInsert.Err())

	var id int64
	err = rowInsert.Scan(&id)
	s.Require().NoError(rowInsert.Err())

	until := time.Now().Add(2 * time.Hour).UTC()

	// Action.
	err = s.repository.ReserveJob(ctx, id, until)

	// Assert.
	s.Require().NoError(err)

	const queryCheck = `select "reserved_until" from "outbox" where id = $1;`
	row := s.DB.QueryRow(ctx, queryCheck, id)
	s.Require().NoError(row.Err())

	var reservedUntil time.Time
	err = row.Scan(&reservedUntil)
	s.Assert().Equal(until.Unix(), reservedUntil.Unix())
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

	const queryInsert = `insert into "outbox" (name, data) values($1, $2) returning id;`
	rowInsert := s.DB.QueryRow(ctx, queryInsert, name, data)
	s.Require().NoError(rowInsert.Err())

	var id int64
	err := rowInsert.Scan(&id)
	s.Require().NoError(rowInsert.Err())

	// Action.
	err = s.repository.DeleteJob(ctx, id)

	// Assert.
	s.Require().NoError(err)

	const query = `select exists(select 1 from "outbox" where "name" = $1 and "data" = $2);`

	row := s.DB.QueryRow(ctx, query, name, data)
	s.Require().NoError(row.Err())

	var exist bool

	err = row.Scan(&exist)
	s.Require().NoError(err)
	s.Assert().False(exist)
}
