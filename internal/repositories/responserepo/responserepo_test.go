package responserepo

import (
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"testing"
	"time"
)

func TestNewResponseRepo(t *testing.T) {
	t.Run("should initiate response repo with empty map when existing responses is nil", func(t *testing.T) {
		responseRepo := NewResponseRepo(nil)
		assert.NotNil(t, responseRepo.responses)
	})
}

func TestResponseRepo_Create(t *testing.T) {
	t.Run("should add response successfully to existing responses", func(t *testing.T) {
		surveyID1 := ksuid.New()
		now := time.Now()
		qID1 := ksuid.New()
		qID2 := ksuid.New()
		existingResponses := map[ksuid.KSUID][]models.Response{
			surveyID1: {
				{
					ID:        ksuid.New(),
					SurveyID:  surveyID1,
					CreatedAt: now,
					Answers: []models.Answer{
						{
							QuestionID: qID1,
							Answer:     true,
						},
						{
							QuestionID: qID2,
							Answer:     false,
						},
					},
				},
			},
		}
		responseRepo := NewResponseRepo(existingResponses)
		newResponse := models.Response{
			ID:        ksuid.New(),
			SurveyID:  surveyID1,
			CreatedAt: now,
			Answers: []models.Answer{
				{
					QuestionID: qID1,
					Answer:     true,
				},
				{
					QuestionID: qID2,
					Answer:     false,
				},
			},
		}
		_, err := responseRepo.Create(&newResponse)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(responseRepo.responses[surveyID1]))
	})
	t.Run("should add new entry when responses are empty", func(t *testing.T) {
		surveyID1 := ksuid.New()
		now := time.Now()
		qID1 := ksuid.New()
		qID2 := ksuid.New()
		existingResponses := map[ksuid.KSUID][]models.Response{}
		responseRepo := NewResponseRepo(existingResponses)
		newResponse := models.Response{
			ID:        ksuid.New(),
			SurveyID:  surveyID1,
			CreatedAt: now,
			Answers: []models.Answer{
				{
					QuestionID: qID1,
					Answer:     true,
				},
				{
					QuestionID: qID2,
					Answer:     false,
				},
			},
		}
		_, err := responseRepo.Create(&newResponse)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(responseRepo.responses[surveyID1]))
	})
	t.Run("should add new entry to new survey when responses are empty", func(t *testing.T) {
		surveyID1 := ksuid.New()
		now := time.Now()
		q1ID1 := ksuid.New()
		q1ID2 := ksuid.New()
		existingResponses := map[ksuid.KSUID][]models.Response{
			surveyID1: {
				{
					ID:        ksuid.New(),
					SurveyID:  surveyID1,
					CreatedAt: now,
					Answers: []models.Answer{
						{
							QuestionID: q1ID1,
							Answer:     true,
						},
						{
							QuestionID: q1ID2,
							Answer:     false,
						},
					},
				},
			},
		}
		surveyID2 := ksuid.New()
		q2ID1 := ksuid.New()
		q2ID2 := ksuid.New()
		responseRepo := NewResponseRepo(existingResponses)
		newResponse := models.Response{
			ID:        ksuid.New(),
			SurveyID:  surveyID2,
			CreatedAt: now,
			Answers: []models.Answer{
				{
					QuestionID: q2ID1,
					Answer:     true,
				},
				{
					QuestionID: q2ID2,
					Answer:     false,
				},
			},
		}
		_, err := responseRepo.Create(&newResponse)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(responseRepo.responses))
		assert.Equal(t, 1, len(responseRepo.responses[surveyID1]))
		assert.Equal(t, 1, len(responseRepo.responses[surveyID2]))
	})
}

func TestResponseRepo_GetBySurveyID(t *testing.T) {
	t.Run("should return all responses for a survey", func(t *testing.T) {
		surveyID1 := ksuid.New()
		now := time.Now()
		qID1 := ksuid.New()
		qID2 := ksuid.New()
		existingResponses := []models.Response{
			{
				ID:        ksuid.New(),
				SurveyID:  surveyID1,
				CreatedAt: now,
				Answers: []models.Answer{
					{
						QuestionID: qID1,
						Answer:     true,
					},
					{
						QuestionID: qID2,
						Answer:     false,
					},
				},
			},
			{
				ID:        ksuid.New(),
				SurveyID:  surveyID1,
				CreatedAt: now,
				Answers: []models.Answer{
					{
						QuestionID: qID1,
						Answer:     false,
					},
					{
						QuestionID: qID2,
						Answer:     true,
					},
				},
			},
		}
		existingEntries := map[ksuid.KSUID][]models.Response{
			surveyID1: existingResponses,
		}
		responseRepo := NewResponseRepo(existingEntries)
		responses, err := responseRepo.GetBySurveyID(surveyID1)
		assert.NoError(t, err)
		assert.Equal(t, existingResponses, responses)
	})
	t.Run("should return error if no responses are found for a surveyID", func(t *testing.T) {
		surveyID1 := ksuid.New()
		now := time.Now()
		qID1 := ksuid.New()
		qID2 := ksuid.New()
		existingResponses := []models.Response{
			{
				ID:        ksuid.New(),
				SurveyID:  surveyID1,
				CreatedAt: now,
				Answers: []models.Answer{
					{
						QuestionID: qID1,
						Answer:     true,
					},
					{
						QuestionID: qID2,
						Answer:     false,
					},
				},
			},
			{
				ID:        ksuid.New(),
				SurveyID:  surveyID1,
				CreatedAt: now,
				Answers: []models.Answer{
					{
						QuestionID: qID1,
						Answer:     false,
					},
					{
						QuestionID: qID2,
						Answer:     true,
					},
				},
			},
		}
		existingEntries := map[ksuid.KSUID][]models.Response{
			surveyID1: existingResponses,
		}
		responseRepo := NewResponseRepo(existingEntries)
		responses, err := responseRepo.GetBySurveyID(ksuid.New())
		assert.Equal(t, repositories.ErrNotFound, err)
		assert.Nil(t, responses)
	})
}

func TestResponseRepo_Entries(t *testing.T) {
	t.Run("should return all entries in the repo", func(t *testing.T) {
		surveyID1 := ksuid.New()
		now := time.Now()
		qID1 := ksuid.New()
		qID2 := ksuid.New()
		existingResponses := []models.Response{
			{
				ID:        ksuid.New(),
				SurveyID:  surveyID1,
				CreatedAt: now,
				Answers: []models.Answer{
					{
						QuestionID: qID1,
						Answer:     true,
					},
					{
						QuestionID: qID2,
						Answer:     false,
					},
				},
			},
			{
				ID:        ksuid.New(),
				SurveyID:  surveyID1,
				CreatedAt: now,
				Answers: []models.Answer{
					{
						QuestionID: qID1,
						Answer:     false,
					},
					{
						QuestionID: qID2,
						Answer:     true,
					},
				},
			},
		}
		existingEntries := map[ksuid.KSUID][]models.Response{
			surveyID1: existingResponses,
		}
		responseRepo := NewResponseRepo(existingEntries)
		entries := responseRepo.Entries()
		assert.Equal(t, existingEntries, entries)
	})
}
