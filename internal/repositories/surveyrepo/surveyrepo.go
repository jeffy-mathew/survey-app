package surveyrepo

import (
	"github.com/segmentio/ksuid"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"sync"
)

type SurveyRepo struct {
	mu      *sync.RWMutex
	surveys map[ksuid.KSUID]models.Survey
}

func NewSurveyRepo(existingSurveys map[ksuid.KSUID]models.Survey) *SurveyRepo {
	if existingSurveys == nil {
		existingSurveys = make(map[ksuid.KSUID]models.Survey)
	}
	return &SurveyRepo{
		mu:      &sync.RWMutex{},
		surveys: existingSurveys,
	}
}

func (s *SurveyRepo) Create(survey *models.Survey) (*models.Survey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.surveys[survey.ID] = *survey
	return survey, nil
}

func (s *SurveyRepo) Get(id ksuid.KSUID) (*models.Survey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	survey, ok := s.surveys[id]
	if !ok {
		return nil, repositories.ErrNotFound
	}
	return &survey, nil
}

func (s *SurveyRepo) Update(id ksuid.KSUID, survey *models.Survey) (*models.Survey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.surveys[id]; !ok {
		return nil, repositories.ErrNotFound
	}
	s.surveys[id] = *survey
	return survey, nil
}

func (s *SurveyRepo) Delete(id ksuid.KSUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.surveys[id]; !ok {
		return repositories.ErrNotFound
	}
	delete(s.surveys, id)
	return nil
}

func (s *SurveyRepo) GetAll() ([]models.Survey, error) {
	s.mu.RLock()
	s.mu.RUnlock()
	if len(s.surveys) == 0 {
		return nil, repositories.ErrNotFound
	}
	var surveys []models.Survey
	for _, survey := range s.surveys {
		surveys = append(surveys, survey)
	}
	return surveys, nil
}

func (s *SurveyRepo) Entries() map[ksuid.KSUID]models.Survey {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.surveys
}
