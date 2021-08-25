package repo

import (
	"encoding/json"
	"fmt"
	"os"
	"ova-exam-api/internal/domain/entity/user"
	"path"
)

// Repo - интерфейс хранилища для сущности User
type Repo interface {
	AddEntities(entities []user.User) error
	ListEntities(limit, offset uint64) ([]user.User, error)
	DescribeEntity(userId uint64) (*user.User, error)
}

// NewRepo возвращает Repo с поддержкой записи в файл
func NewRepo(
	fileName string,
) Repo {
	return &repo{
		fileName: fileName,
	}
}

type repo struct {
	fileName string
}

func (r *repo) AddEntities(entities []user.User) error {
	pwd, pathErr := os.Getwd()
	fileName := path.Join(pwd, r.fileName)
	if pathErr != nil {
		return pathErr
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("close file", file)
		file.Close()
	}()

	fmt.Println("open file", file)

	data, marshalErr := json.MarshalIndent(entities, "", " ")
	if marshalErr != nil {
		return marshalErr
	}

	if _, err := file.Write(data); err != nil {
		return err
	}
	fmt.Printf("users %d wrote to file\n", len(entities))
	return nil
}

func (r repo) ListEntities(limit, offset uint64) ([]user.User, error) {
	panic("implement me")
}

func (r repo) DescribeEntity(userId uint64) (*user.User, error) {
	panic("implement me")
}
