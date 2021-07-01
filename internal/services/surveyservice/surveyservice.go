package surveyservice

import (
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
	"log"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"survey-platform/pkg/idgenerator"
	"survey-platform/pkg/timegenerator"
)

type SurveyService struct {
	maxQuestions  int
	surveyRepo    repositories.SurveyRepoInterface
	responseRepo  repositories.ResponseRepoInterface
	idGenerator   idgenerator.IDGenerator
	timeGenerator timegenerator.TimeGenInterface
}

func NewSurveyService(maxQuestions int, surveyRepo repositories.SurveyRepoInterface,
	responseRepo repositories.ResponseRepoInterface, idGenerator idgenerator.IDGenerator,
	timeGenerator timegenerator.TimeGenInterface) *SurveyService {
	return &SurveyService{
		maxQuestions:  maxQuestions,
		surveyRepo:    surveyRepo,
		responseRepo:  responseRepo,
		idGenerator:   idGenerator,
		timeGenerator: timeGenerator,
	}
}

func (s *SurveyService) CreateSurvey(survey models.Survey) (*models.Survey, error) {
	if len(survey.Questions) > s.maxQuestions {
		return nil, errors.New("survey cannot have more than 3 questions")
	}
	if len(survey.Questions) == 0 {
		return nil, errors.New("survey cannot be empty")
	}
	if survey.Name == "" {
		return nil, errors.New("survey needs a name")
	}
	survey.ID = s.idGenerator.Generate()
	questions := survey.Questions
	for i := 0; i < len(questions); i++ {
		questions[i].ID = s.idGenerator.Generate()
	}
	survey.Questions = questions
	now := s.timeGenerator.Now()
	survey.CreatedAt, survey.UpdatedAt = now, now
	return s.surveyRepo.Create(&survey)
}

func (s *SurveyService) GetSurvey(id ksuid.KSUID) (*models.Survey, error) {
	return s.surveyRepo.Get(id)
}

func (s *SurveyService) UpdateSurvey(id ksuid.KSUID, survey models.Survey) (*models.Survey, error) {
	survey.ID = id
	questions := survey.Questions
	for i := 0; i < len(questions); i++ {
		if questions[i].ID.IsNil() || questions[i].ID.String() == "" {
			questions[i].ID = s.idGenerator.Generate()
		}
	}
	survey.Questions = questions
	survey.UpdatedAt = s.timeGenerator.Now()
	return s.surveyRepo.Update(id, &survey)
}

func (s *SurveyService) DeleteSurvey(id ksuid.KSUID) error {
	return s.surveyRepo.Delete(id)
}

func (s *SurveyService) GetAllSurveys() ([]models.Survey, error) {
	return s.surveyRepo.GetAll()
}

func (s *SurveyService) SaveResponse(response models.Response) (*models.Response, error) {
	_, err := s.surveyRepo.Get(response.SurveyID)
	if err != nil {
		log.Println("error while validating survey for response ", err, "surveyID ", response.SurveyID.String())
		return nil, err
	}
	if len(response.Answers) > s.maxQuestions {
		return nil, fmt.Errorf("max number of questions allowed is %d", s.maxQuestions)
	}
	response.ID = s.idGenerator.Generate()
	response.CreatedAt = s.timeGenerator.Now()
	return s.responseRepo.Create(&response)
}

func (s *SurveyService) GetResponses(surveyID ksuid.KSUID) ([]models.Response, error) {
	return s.responseRepo.GetBySurveyID(surveyID)
}

func (s *SurveyService) Entries() *models.DBEntry {
	return &models.DBEntry{
		Responses: s.responseRepo.Entries(),
		Surveys:   s.surveyRepo.Entries(),
	}
}
