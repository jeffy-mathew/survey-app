package repositories

//go:generate mockgen -source=repositories.go -destination=./repositories_mock/repositories_mock.go -package=repositories_mock

import (
	"errors"
	"github.com/segmentio/ksuid"
	"survey-platform/internal/models"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type SurveyRepoInterface interface {
	Create(survey *models.Survey) (*models.Survey, error)
	Get(id ksuid.KSUID) (*models.Survey, error)
	Update(id ksuid.KSUID, survey *models.Survey) (*models.Survey, error)
	Delete(id ksuid.KSUID) error
	GetAll() ([]models.Survey, error)
	Entries() map[ksuid.KSUID]models.Survey
}

type ResponseRepoInterface interface {
	Create(response *models.Response) (*models.Response, error)
	GetBySurveyID(surveyID ksuid.KSUID) ([]models.Response, error)
	Entries() map[ksuid.KSUID][]models.Response
}
