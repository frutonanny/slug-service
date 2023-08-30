package users_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	cnf "github.com/frutonanny/slug-service/internal/config"
	"github.com/frutonanny/slug-service/internal/database"
	"github.com/frutonanny/slug-service/internal/repositories/slug"
	"github.com/frutonanny/slug-service/internal/repositories/users"
)

const fileConfig = "../../../config/config.local.json"

var config = cnf.Must(fileConfig)

type RepositorySuite struct {
	database.DBSuite

	usersRepo *users.Repository
	slugRepo  *slug.Repository
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &RepositorySuite{DBSuite: database.NewDBSuite(config.DB.DSN)})
}

func (s *RepositorySuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	s.usersRepo = users.New(s.DB)
	s.slugRepo = slug.New(s.DB)
}

func (s *RepositorySuite) TestRepository_CreateUserIfNotExist() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	userID := uuid.New()

	// Action.
	err := s.usersRepo.CreateUserIfNotExist(ctx, userID)

	// Assert.
	s.Require().NoError(err)

	const query = `select exists(select 1 from "users" where id = $1);`

	row := s.DB.QueryRow(ctx, query, userID)
	s.Require().NoError(row.Err())

	var exist bool

	err = row.Scan(&exist)
	s.Require().NoError(err)
	s.Assert().True(exist)
}

func (s *RepositorySuite) TestRepository_AddUserSlug() {
	s.Run("add user_slug success", func() {
		// Arrange.
		ctx, cancel := s.Context()
		defer cancel()

		userID := uuid.New()
		name := "123"

		err := s.usersRepo.CreateUserIfNotExist(ctx, userID)
		s.Require().NoError(err)

		slugID, err := s.slugRepo.Create(ctx, name, slug.Options{})
		s.Require().NoError(err)

		// Action.
		id, err := s.usersRepo.AddUserSlug(ctx, userID, slugID, name)

		// Assert.
		s.Require().NoError(err)

		const query = `select "user_id", "slug_id"
				from "users_slugs" 
				where id = $1;`
		row := s.DB.QueryRow(ctx, query, id)
		s.Require().NoError(row.Err())

		var (
			actualUserID uuid.UUID
			actualSlugID int64
		)

		err = row.Scan(&actualUserID, &actualSlugID)
		s.Require().NoError(err)

		s.Assert().EqualValues(userID, actualUserID)
		s.Assert().EqualValues(slugID, actualSlugID)
	})

	s.Run("add user_slug error: nameID and slugID not exist", func() {
		// Arrange.
		ctx, cancel := s.Context()
		defer cancel()

		userID := uuid.New()
		slugID := int64(1)
		name := "123"

		// Action.
		_, err := s.usersRepo.AddUserSlug(ctx, userID, slugID, name)

		// Assert.
		s.Assert().Error(err)
	})

}

func (s *RepositorySuite) TestRepository_AddUserSlugWithTtl() {
	s.Run("add user_slug with ttl success", func() {
		// Arrange.
		ctx, cancel := s.Context()
		defer cancel()

		userID := uuid.New()
		name := "123"
		ttl := time.Now().Add(1 * time.Hour)

		err := s.usersRepo.CreateUserIfNotExist(ctx, userID)
		s.Require().NoError(err)

		slugID, err := s.slugRepo.Create(ctx, name, slug.Options{})
		s.Require().NoError(err)

		// Action.
		id, err := s.usersRepo.AddUserSlugWithTtl(ctx, userID, slugID, name, ttl)

		// Assert.
		s.Require().NoError(err)

		const query = `select "user_id", "slug_id", "slug_ttl"
				from "users_slugs" 
				where id = $1;`
		row := s.DB.QueryRow(ctx, query, id)
		s.Require().NoError(row.Err())

		var (
			actualUserID uuid.UUID
			actualSlugID int64
			actualTtl    time.Time
		)

		err = row.Scan(&actualUserID, &actualSlugID, &actualTtl)
		s.Require().NoError(err)

		s.Assert().EqualValues(userID, actualUserID)
		s.Assert().EqualValues(slugID, actualSlugID)
		s.Assert().False(actualTtl.IsZero())
	})

	s.Run("add user_slug with ttl error: nameID and slugID not exist", func() {
		// Arrange.
		ctx, cancel := s.Context()
		defer cancel()

		userID := uuid.New()
		slugID := int64(1)
		name := "123"
		ttl := time.Now().Add(1 * time.Hour)

		// Action.
		_, err := s.usersRepo.AddUserSlugWithTtl(ctx, userID, slugID, name, ttl)

		// Assert.
		s.Assert().Error(err)
	})
}

func (s *RepositorySuite) TestRepository_DeleteUserSlug() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	userID := uuid.New()
	name := "123"

	err := s.usersRepo.CreateUserIfNotExist(ctx, userID)
	s.Require().NoError(err)

	slugID, err := s.slugRepo.Create(ctx, name, slug.Options{})
	s.Require().NoError(err)

	// Action.
	err = s.usersRepo.DeleteUserSlug(ctx, userID, slugID)

	// Assert.
	s.Require().NoError(err)

	const query = `select "user_id"
				from "users_slugs" 
				where slug_id = $1;`

	row := s.DB.QueryRow(ctx, query, slugID)
	s.Assert().Error(row.Scan())
}

func (s *RepositorySuite) TestRepository_GetUserSlugs() {
	// Arrange.
	ctx, cancel := s.Context()
	defer cancel()

	userID := uuid.New()
	name := "123"

	err := s.usersRepo.CreateUserIfNotExist(ctx, userID)
	s.Require().NoError(err)

	slugID, err := s.slugRepo.Create(ctx, name, slug.Options{})
	s.Require().NoError(err)

	_, err = s.usersRepo.AddUserSlug(ctx, userID, slugID, name)
	s.Require().NoError(err)

	// Action.
	slugs, err := s.usersRepo.GetUserSlugs(ctx, userID)

	// Assert.
	s.Require().NoError(err)

	const query = `select "slug_id", "slug_name"
				from "users_slugs" 
				where "user_id" = $1;`
	row := s.DB.QueryRow(ctx, query, userID)
	s.Require().NoError(row.Err())

	var (
		actualID   int64
		actualName string
	)

	err = row.Scan(&actualID, &actualName)
	s.Require().NoError(err)

	s.Require().NotEmpty(slugs)
	s.Require().Len(slugs, 1)
	s.Assert().EqualValues(slugs[0].ID, actualID)
	s.Assert().EqualValues(slugs[0].Name, actualName)
}
