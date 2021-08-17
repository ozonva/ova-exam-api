package answer

import (
	"ova-exam-api/internal/domain/entity"
	"ova-exam-api/internal/domain/entity/question"
)

type Answer struct {
	entity.Entity
	AnswerId uint64
	Text     string
	Question question.Question
	IsRight  bool
}
