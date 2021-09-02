package repo

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
	"ova-exam-api/internal/domain/entity"
	"ova-exam-api/internal/domain/entity/user"
	"time"
)

// Repo - интерфейс хранилища для сущности User
type Repo interface {
	AddEntities(entities []user.User) error
	ListEntities(limit, offset uint64) ([]user.User, error)
	DescribeEntity(userId uint64) (*user.User, error)
	RemoveEntity(userId uint64) error
}

// NewRepo возвращает Repo
func NewRepo(
	db sq.BaseRunner,
) Repo {
	return &repo{
		tableName: "users",
		db: db,
	}
}

type repo struct {
	tableName string
	db sq.BaseRunner
}

func (r *repo) AddEntities(users []user.User) error {
	query := sq.Insert(r.tableName).
		Columns("Email", "Password", "createdAt").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	for _, newUser := range users {
		query = query.
			Values(newUser.Email, newUser.Password, time.Now())
	}

	_, err := query.Exec()

	return err
}

func (r repo) ListEntities(limit, offset uint64) ([]user.User, error) {
	query := sq.Select("Id", "Email", "Password", "createdAt", "updatedAt").
		From("users").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar).
		Offset(offset).
		Limit(limit)
	log.Print(query.ToSql())

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]user.User, 0)
	for rows.Next() {
		var id uint64
		var email, password string
		var createat time.Time
		var updateatDb sql.NullTime
		if err := rows.Scan(&id, &email, &password, &createat, &updateatDb); err != nil {
			return nil, err
		}

		var updateat time.Time
		if updateatDb.Valid{
			updateat = updateatDb.Time
		}

		existUser := user.User{
			Entity:   entity.Entity{
				CreatedAt: createat,
				UpdatedAt: updateat,
			},
			UserId:   id,
			Email:    email,
			Password: password,
		}
		result = append(result, existUser)
		// обрабатываем строку
		fmt.Printf("%d %s %s %s %s\n", id, email, password, createat, updateat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r repo) DescribeEntity(userId uint64) (*user.User, error) {
	query := sq.Select("Id", "Email", "Password", "createdAt", "updatedAt").
		From("users").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userId})

	log.Print(query.ToSql())

	var id uint64
	var email, password string
	var createat time.Time
	var updateatDb sql.NullTime

	err := query.QueryRow().Scan(&id, &email, &password, &createat, &updateatDb)
	if err != nil {
		return nil, err
	}
	var updateat time.Time
	if updateatDb.Valid{
		updateat = updateatDb.Time
	}

	existUser := user.User{
		Entity:   entity.Entity{
			CreatedAt: createat,
			UpdatedAt: updateat,
		},
		UserId:   id,
		Email:    email,
		Password: password,
	}
	fmt.Printf("%d %s %s %s %s\n", id, email, password, createat, updateat)

	return &existUser, nil
}

func (r repo) RemoveEntity(userId uint64) error {
	query := sq.Delete(r.tableName).
		Where(sq.Eq{"id": userId}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	_, err := query.Exec()

	return err
}