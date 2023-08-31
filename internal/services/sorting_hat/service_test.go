package sorting_hat_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	sortinghat "github.com/frutonanny/slug-service/internal/services/sorting_hat"
)

func TestService_HitOrNot(t *testing.T) {
	percent := int64(30)
	s := sortinghat.New()

	t.Run("hit", func(t *testing.T) {
		// Arrange.
		userIDtHit := uuid.MustParse("030c48d3-12a3-4ea0-a334-d9a141c7a1b8") // Должен быть ID – 51136723.

		// Action.
		hit := s.Hit(userIDtHit, percent)

		// Assert.
		assert.True(t, hit)
	})

	t.Run("no hit", func(t *testing.T) {
		// Arrange.
		userIDNotHit := uuid.MustParse("9305ea2f-0c7c-4a08-a1fa-918a496e1701") // Должен быть ID – 2466638383.

		// Action.
		hit := s.Hit(userIDNotHit, percent)

		// Assert.
		assert.False(t, hit)
	})
}
