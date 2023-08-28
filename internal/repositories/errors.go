package repositories

import "errors"

// Ошибки, о которых необходимо сообщить другому сервису / пользователю.

var (
	ErrRepoSlugNotFound = errors.New("slug not found")
)
