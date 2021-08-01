package models

import (
	"github.com/segmentio/ksuid"
	"time"
)

type Survey struct {
	ID        ksuid.KSUID `json:"id" example:"-"`
	Name      string      `json:"name" example:"account name"`
	Questions []Question  `json:"questions"`
	CreatedAt time.Time   `json:"created_at" example:"-"`
	UpdatedAt time.Time   `json:"updated_at" example:"-"`
}

type Question struct {
	ID       ksuid.KSUID `json:"id" example:"ksuid"`
	Question string      `json:"question" example:"question"`
}

type Answer struct {
	QuestionID ksuid.KSUID `json:"question_id"`
	Answer     bool        `json:"answer"`
}

type Response struct {
	ID        ksuid.KSUID `json:"id"`
	SurveyID  ksuid.KSUID `json:"survey_id"`
	Answers   []Answer    `json:"answers"`
	CreatedAt time.Time   `json:"created_at"`
}

type DBEntry struct {
	Surveys   map[ksuid.KSUID]Survey     `json:"surveys"`
	Responses map[ksuid.KSUID][]Response `json:"responses"`
}
