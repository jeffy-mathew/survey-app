package services

//go:generate mockgen -source=services.go -destination=./services_mock/services_mock.go -package=services_mock

import (
	"github.com/segmentio/ksuid"
	"survey-platform/internal/models"
)

type SurveyServiceInterface interface {
	CreateSurvey(survey models.Survey) (*models.Survey, error)
	GetSurvey(id ksuid.KSUID) (*models.Survey, error)
	UpdateSurvey(id ksuid.KSUID, survey models.Survey) (*models.Survey, error)
	DeleteSurvey(id ksuid.KSUID) error
	GetAllSurveys() ([]models.Survey, error)
	SaveResponse(response models.Response) (*models.Response, error)
	GetResponses(surveyID ksuid.KSUID) ([]models.Response, error)
	Entries() *models.DBEntry
}
