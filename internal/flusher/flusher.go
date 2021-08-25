package flusher

import (
	"ova-exam-api/internal/domain/entity/user"
	"ova-exam-api/internal/repo"
	"ova-exam-api/internal/utils"
)

// Flusher - интерфейс для сброса задач в хранилище
type Flusher interface {
	Flush(entities []user.User) []user.User
}

// NewFlusher возвращает Flusher с поддержкой батчевого сохранения
func NewFlusher(
	chunkSize int,
	entityRepo repo.Repo,
) Flusher {
	return &flusher{
		chunkSize: chunkSize,
		entityRepo:  entityRepo,
	}
}

type flusher struct {
	chunkSize int
	entityRepo  repo.Repo
}

func (f *flusher) Flush(users []user.User) []user.User {
	bulks := utils.SplitToBulks(users, uint(f.chunkSize))
	for n, value := range bulks {
		err := f.entityRepo.AddEntities(value)
		if err != nil {
			// Возвращает пользователей которых не удалось сохранить
			return users[n * f.chunkSize:]
		}
	}

	return nil
}
