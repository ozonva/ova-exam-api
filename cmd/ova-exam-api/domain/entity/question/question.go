package question

import (
	"ova-exam-api/cmd/ova-exam-api/domain/entity"
	"ova-exam-api/cmd/ova-exam-api/domain/entity/answer"
)

type Question struct {
	entity.Entity
	QuestionId 	uint64
	Text        string
	Answers 	[]answer.Answer
}
