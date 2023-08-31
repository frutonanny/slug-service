// Функция "попадания" пользователя в процент. Исходим из того, что идентификатор пользователя – UUID, который можем
// привести к целому числу. Глядя на остаток деления числа на 100 понимаем, попали ли в этот процент (реализация в
// пакете sorting_hat). На большом количестве пользователей должно быть нормальное распределение.

package sorting_hat

import (
	"github.com/google/uuid"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Hit(userID uuid.UUID, percent int64) bool {
	numID := userID.ID()
	return numID%100 <= uint32(percent)
}
