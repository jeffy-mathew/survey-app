package responserepo

import (
	"github.com/segmentio/ksuid"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"sync"
)

type ResponseRepo struct {
	mu        *sync.RWMutex
	responses map[ksuid.KSUID][]models.Response
}

func NewResponseRepo(existingResponses map[ksuid.KSUID][]models.Response) *ResponseRepo {
	if existingResponses == nil {
		existingResponses = make(map[ksuid.KSUID][]models.Response)
	}
	return &ResponseRepo{
		mu:        &sync.RWMutex{},
		responses: existingResponses,
	}
}

func (r *ResponseRepo) Create(response *models.Response) (*models.Response, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existingResponses, ok := r.responses[response.SurveyID]
	if !ok {
		r.responses[response.SurveyID] = []models.Response{*response}
		return response, nil
	}
	existingResponses = append(existingResponses, *response)
	r.responses[response.SurveyID] = existingResponses
	return response, nil
}

func (r *ResponseRepo) GetBySurveyID(surveyID ksuid.KSUID) ([]models.Response, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	responses, ok := r.responses[surveyID]
	if !ok {
		return nil, repositories.ErrNotFound
	}
	return responses, nil
}

func (r *ResponseRepo) Entries() map[ksuid.KSUID][]models.Response {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.responses
}
