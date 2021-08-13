package answer

import (
	"ova-exam-api/cmd/ova-exam-api/domain/entity"
	"ova-exam-api/cmd/ova-exam-api/domain/entity/question"
)

type Answer struct {
	entity.Entity
	AnswerId  uint64
	Text      string
	Question  question.Question
	IsRight   bool
}
