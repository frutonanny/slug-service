package slug_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	cnf "github.com/frutonanny/slug-service/internal/config"
	"github.com/frutonanny/slug-service/internal/database"
	"github.com/frutonanny/slug-service/internal/repositories/slug"
)

const fileConfig = "../../../config/config.local.json"

var config = cnf.Must(fileConfig)

type RepositorySuite struct {
	database.DBSuite

	repository *slug.Repository
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &RepositorySuite{DBSuite: database.NewDBSuite(config.DB.DSN)})
}

func (s *RepositorySuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	s.repository = slug.New(s.DB)
}

func (s *RepositorySuite) TestRepository_Create() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	name := "123"
	percent := 20
	option := slug.Options{Percent: &percent}

	// Action.
	id, err := s.repository.Create(ctx, name, option)

	// Assert.
	s.Require().NoError(err)

	const query = `select "name", "options"
				from "slugs" 
				where id = $1;`
	row := s.DB.QueryRow(ctx, query, id)
	s.Require().NoError(row.Err())

	var (
		actualName   string
		actualOption []byte
	)

	err = row.Scan(&actualName, &actualOption)
	s.Require().NoError(err)

	s.Require().EqualValues(name, actualName)
	s.Require().NotEmpty(actualOption)

	var actualPercent slug.Options
	err = json.Unmarshal(actualOption, &actualPercent)
	s.Require().NoError(err)
	s.Assert().EqualValues(percent, *actualPercent.Percent)
}

func (s *RepositorySuite) TestRepository_Delete_Success() {
	s.Run("delete success", func() {
		// Arrange.
		ctx, cancel := s.Context()
		defer cancel()

		name := "123"

		id, err := s.repository.Create(ctx, name, slug.Options{})
		s.Require().NoError(err)

		// Action.
		err = s.repository.Delete(ctx, name)

		// Assert.
		s.Require().NoError(err)

		const query = `select "name", "deleted_at"
				from "slugs" 
				where id = $1;`
		row := s.DB.QueryRow(ctx, query, id)
		s.Require().NoError(row.Err())

		var (
			actualName string
			deletedAt  time.Time
		)

		err = row.Scan(&actualName, &deletedAt)
		s.Require().NoError(err)

		s.Assert().EqualValues(name, actualName)
		s.Assert().False(deletedAt.IsZero())
	})

	s.Run("delete with error: slug not found", func() {
		// Arrange.
		ctx, cancel := s.Context()
		defer cancel()

		name := "123"

		// Action.
		err := s.repository.Delete(ctx, name)

		// Assert.
		s.Require().Error(err)
		s.ErrorIs(err, slug.ErrRepoSlugNotFound)
	})
}

func (s *RepositorySuite) TestRepository_GetID() {

	s.Run("geiID success", func() {
		// Arrange.
		ctx, cancel := s.Context()
		defer cancel()

		name := "123"

		expectedID, err := s.repository.Create(ctx, name, slug.Options{})
		s.Require().NoError(err)

		// Action.
		actualID, err := s.repository.GetID(ctx, name)

		// Assert.
		s.Require().NoError(err)
		s.Assert().EqualValues(expectedID, actualID)
	})

	s.Run("getID with error: slug not found", func() {
		// Arrange.
		ctx, cancel := s.Context()
		defer cancel()

		name := "123"

		// Action.
		_, err := s.repository.GetID(ctx, name)

		// Assert.
		s.Require().Error(err)
		s.ErrorIs(err, slug.ErrRepoSlugNotFound)
	})
}
