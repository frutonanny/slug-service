package events_test

import (
	"testing"

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

func (s *RepositorySuite) TestRepository_AddUsersSlugsEvent() {
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
