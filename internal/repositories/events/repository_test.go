package events_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	cnf "github.com/frutonanny/slug-service/internal/config"
	"github.com/frutonanny/slug-service/internal/database"
	"github.com/frutonanny/slug-service/internal/repositories/events"
)

const fileConfig = "../../../config/config.local.json"

var config = cnf.Must(fileConfig)

type RepositorySuite struct {
	database.DBSuite

	repository *events.Repository
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &RepositorySuite{DBSuite: database.NewDBSuite(config.DB.DSN)})
}

func (s *RepositorySuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	s.repository = events.New(s.DB)
}

func (s *RepositorySuite) TestRepository_AddEvent() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	userID := uuid.New()
	slugID := int64(1)
	event := "something happened"

	// Action.
	id, err := s.repository.AddEvent(ctx, userID, slugID, event)

	// Assert.
	s.Require().NoError(err)

	const query = `select user_id, slug_id, event 
				from "events" 
				where id = $1;`
	row := s.DB.QueryRow(ctx, query, id)
	s.Require().NoError(row.Err())

	var (
		actualUserID uuid.UUID
		actualSlugID int64
		actualEvent  string
	)

	err = row.Scan(&actualUserID, &actualSlugID, &actualEvent)
	s.Require().NoError(err)

	s.Assert().EqualValues(userID, actualUserID)
	s.Assert().EqualValues(slugID, actualSlugID)
	s.Assert().EqualValues(event, actualEvent)
}

func (s *RepositorySuite) TestRepository_AddEventWithCreatedAt() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	userID := uuid.New()
	slugID := int64(1)
	event := "something happened"
	createdAt := time.Now().Add(1 * time.Hour)

	// Action.
	id, err := s.repository.AddEventWithCreatedAt(ctx, userID, slugID, event, createdAt)

	// Assert.
	s.Require().NoError(err)

	const query = `select user_id, slug_id, event, created_at 
				from "events" 
				where id = $1;`
	row := s.DB.QueryRow(ctx, query, id)
	s.Require().NoError(row.Err())

	var (
		actualUserID    uuid.UUID
		actualSlugID    int64
		actualEvent     string
		actualCreatedAt time.Time
	)

	err = row.Scan(&actualUserID, &actualSlugID, &actualEvent, &actualCreatedAt)
	s.Require().NoError(err)

	s.Assert().EqualValues(userID, actualUserID)
	s.Assert().EqualValues(slugID, actualSlugID)
	s.Assert().EqualValues(event, actualEvent)
	s.Assert().EqualValues(createdAt.Minute(), actualCreatedAt.Minute())
}

func (s *RepositorySuite) TestRepository_GetReport() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	userID := uuid.New()
	name := "123"
	event := "something happened"

	period := time.Now().Format("2006-01")
	from, to := s.parsePeriod(period)

	const queryAdd = `insert into "slugs" (name) values ($1) returning id;`
	row := s.DB.QueryRow(ctx, queryAdd, name)
	s.Require().NoError(row.Err())

	var slugID int64
	err := row.Scan(&slugID)
	s.Require().NoError(err)

	_, err = s.repository.AddEvent(ctx, userID, slugID, event)

	// Action.
	result, err := s.repository.GetReport(ctx, userID, from, to)

	// Assert.
	s.Require().NoError(err)

	s.Assert().Len(result, 1)
	s.Assert().EqualValues(result[0].UserID, userID)
	s.Assert().EqualValues(result[0].EventName, event)
}

func (s *RepositorySuite) parsePeriod(period string) (from, to time.Time) {
	from, err := time.Parse("2006-01", period)
	s.Assert().NoError(err)

	to = from.AddDate(0, 1, 0)

	return from, to
}
