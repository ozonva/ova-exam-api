package repo
import (
	"ova-exam-api/internal/domain/entity/user"
)

// Repo - интерфейс хранилища для сущности User
type Repo interface {
	AddEntities(entities []user.User) error
	ListEntities(limit, offset uint64) ([]user.User, error)
	DescribeEntity(userId uint64) (*user.User, error)
}