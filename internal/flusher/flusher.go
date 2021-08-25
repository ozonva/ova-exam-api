package flusher

import (
	"ova-exam-api/internal/repo"
	"ova-exam-api/internal/utils"
	"ova-exam-api/internal/domain/entity/user"
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
	storageUsers []user.User
	chunkSize int
	entityRepo  repo.Repo
}

func (f flusher) Flush(users []user.User) []user.User {
	f.storageUsers = append(f.storageUsers, users...)

	bulks := utils.SplitToBulks(f.storageUsers, uint(f.chunkSize))
	for _, value := range bulks {
		if len(value) < f.chunkSize{
			f.storageUsers = value

			return f.storageUsers
		}
		f.entityRepo.AddEntities(value)
	}

	return nil
}
