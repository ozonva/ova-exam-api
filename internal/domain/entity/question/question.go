package question

import (
	"ova-exam-api/internal/domain/entity"
	"ova-exam-api/internal/domain/entity/answer"
)

type Question struct {
	entity.Entity
	QuestionId 	uint64
	Text        string
	Answers 	[]answer.Answer
}
