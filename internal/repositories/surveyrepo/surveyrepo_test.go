package surveyrepo

import (
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"testing"
	"time"
)

func TestNewSurveyRepo(t *testing.T) {
	t.Run("should initiate surverepo with empty map when existing survey is nil", func(t *testing.T) {
		responseRepo := NewSurveyRepo(nil)
		assert.NotNil(t, responseRepo.surveys)
	})
}

func TestSurveyRepo_Create(t *testing.T) {
	t.Run("should successfully create survey", func(t *testing.T) {
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{})
		survey := models.Survey{
			ID:        ksuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       ksuid.New(),
					Question: "is this place good?",
				},
				{
					ID:       ksuid.New(),
					Question: "does this place has parking?",
				},
			}}
		newSurvey, err := surveyRepo.Create(&survey)
		assert.NoError(t, err)
		assert.Equal(t, survey, *newSurvey)
	})
	t.Run("should add new survey to existing surveys", func(t *testing.T) {
		survey1 := models.Survey{
			ID:        ksuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       ksuid.New(),
					Question: "is this place good?",
				},
				{
					ID:       ksuid.New(),
					Question: "does this place has parking?",
				},
			}}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{survey1.ID: survey1})
		survey2 := models.Survey{
			ID:        ksuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       ksuid.New(),
					Question: "is this place good enough?",
				},
				{
					ID:       ksuid.New(),
					Question: "does this place has wheelchair accessible parking?",
				},
			}}
		newSurvey, err := surveyRepo.Create(&survey2)
		assert.NoError(t, err)
		assert.Equal(t, survey2, *newSurvey)
		assert.Equal(t, 2, len(surveyRepo.surveys))
	})
}

func TestSurveyRepo_Get(t *testing.T) {
	t.Run("should successfully get survey by id", func(t *testing.T) {
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		q2ID := ksuid.New()
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       q1ID,
					Question: "is this place good?",
				},
				{
					ID:       q2ID,
					Question: "does this place has parking?",
				},
			},
		}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{
			surveyID: survey,
		})
		expectedSurvey := survey
		newSurvey, err := surveyRepo.Get(surveyID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSurvey, *newSurvey)
	})
	t.Run("should return error if no surveys are found by id", func(t *testing.T) {
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		q2ID := ksuid.New()
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       q1ID,
					Question: "is this place good?",
				},
				{
					ID:       q2ID,
					Question: "does this place has parking?",
				},
			},
		}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{
			surveyID: survey,
		})
		newSurvey, err := surveyRepo.Get(ksuid.New())
		assert.Equal(t, repositories.ErrNotFound, err)
		assert.Nil(t, newSurvey)
	})
}

func TestSurveyRepo_Update(t *testing.T) {
	t.Run("should successfully update survey by id", func(t *testing.T) {
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		q2ID := ksuid.New()
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       q1ID,
					Question: "is this place good?",
				},
				{
					ID:       q2ID,
					Question: "does this place has parking?",
				},
			},
		}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{
			surveyID: survey,
		})
		updatedSurvey := survey
		q2ID1, q2ID2 := ksuid.New(), ksuid.New()
		updatedSurvey.Questions = []models.Question{
			{
				ID:       q2ID1,
				Question: "is apple m1 chip good?",
			},
			{
				ID:       q2ID2,
				Question: "do you prefer linux over MacOS",
			},
		}
		newSurvey, err := surveyRepo.Update(surveyID, &updatedSurvey)
		assert.NoError(t, err)
		assert.NotEqual(t, survey.Questions, newSurvey.Questions)
	})
	t.Run("should return error if survey for id does not exist", func(t *testing.T) {
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		q2ID := ksuid.New()
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       q1ID,
					Question: "is this place good?",
				},
				{
					ID:       q2ID,
					Question: "does this place has parking?",
				},
			},
		}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{
			surveyID: survey,
		})
		updatedSurvey := survey
		q2ID1, q2ID2 := ksuid.New(), ksuid.New()
		updatedSurvey.Questions = []models.Question{
			{
				ID:       q2ID1,
				Question: "is apple m1 chip good?",
			},
			{
				ID:       q2ID2,
				Question: "do you prefer linux over MacOS",
			},
		}
		newSurvey, err := surveyRepo.Update(ksuid.New(), &updatedSurvey)
		assert.Equal(t, repositories.ErrNotFound, err)
		assert.Nil(t, newSurvey)
	})
}

func TestSurveyRepo_GetAll(t *testing.T) {
	t.Run("should return all surveys as array", func(t *testing.T) {
		surveyID1, surveyID2 := ksuid.New(), ksuid.New()
		q1ID1, q1ID2, q2ID1, q2ID2 := ksuid.New(), ksuid.New(), ksuid.New(), ksuid.New()
		survey1 := models.Survey{
			Name:      "new survey1",
			ID:        surveyID1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       q1ID1,
					Question: "is this place good?",
				},
				{
					ID:       q1ID2,
					Question: "does this place has parking?",
				},
			},
		}
		survey2 := models.Survey{
			Name:      "new survey 2",
			ID:        surveyID2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       q2ID1,
					Question: "is apple m1 chip good?",
				},
				{
					ID:       q2ID2,
					Question: "do you prefer linux over MacOS",
				},
			},
		}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{
			surveyID1: survey1,
			surveyID2: survey2,
		})
		surveys, err := surveyRepo.GetAll()
		assert.NoError(t, err)
		assert.EqualValues(t, []models.Survey{survey1, survey2}, surveys)
	})
	t.Run("should return error if no surveys are found", func(t *testing.T) {
		surveyRepo := NewSurveyRepo(nil)
		newSurvey, err := surveyRepo.GetAll()
		assert.Equal(t, repositories.ErrNotFound, err)
		assert.Nil(t, newSurvey)
	})
}

func TestSurveyRepo_Delete(t *testing.T) {
	t.Run("should successfully delete survey by id", func(t *testing.T) {
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		q2ID := ksuid.New()
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       q1ID,
					Question: "is this place good?",
				},
				{
					ID:       q2ID,
					Question: "does this place has parking?",
				},
			},
		}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{
			surveyID: survey,
		})
		err := surveyRepo.Delete(surveyID)
		assert.NoError(t, err)
		assert.Empty(t, surveyRepo.surveys)
	})
	t.Run("should return error if survey for id does not exist", func(t *testing.T) {
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		q2ID := ksuid.New()
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       q1ID,
					Question: "is this place good?",
				},
				{
					ID:       q2ID,
					Question: "does this place has parking?",
				},
			},
		}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{
			surveyID: survey,
		})
		err := surveyRepo.Delete(ksuid.New())
		assert.Error(t, repositories.ErrNotFound, err)
	})
}

func TestSurveyRepo_Entries(t *testing.T) {
	t.Run("should return all entries in the repo", func(t *testing.T) {
		survey1 := models.Survey{
			ID:        ksuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       ksuid.New(),
					Question: "is this place good?",
				},
				{
					ID:       ksuid.New(),
					Question: "does this place has parking?",
				},
			}}
		survey2 := models.Survey{
			ID:        ksuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       ksuid.New(),
					Question: "is this place good enough?",
				},
				{
					ID:       ksuid.New(),
					Question: "does this place has wheelchair accessible parking?",
				},
			}}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{survey1.ID: survey1, survey2.ID: survey2})
		entries := surveyRepo.Entries()
		assert.Equal(t, map[ksuid.KSUID]models.Survey{survey1.ID: survey1, survey2.ID: survey2}, entries)
	})
	t.Run("should return all entries in the repo after creating one survey", func(t *testing.T) {
		survey1 := models.Survey{
			ID:        ksuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       ksuid.New(),
					Question: "is this place good?",
				},
				{
					ID:       ksuid.New(),
					Question: "does this place has parking?",
				},
			}}
		surveyRepo := NewSurveyRepo(map[ksuid.KSUID]models.Survey{survey1.ID: survey1})
		survey2 := models.Survey{
			ID:        ksuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       ksuid.New(),
					Question: "is this place good enough?",
				},
				{
					ID:       ksuid.New(),
					Question: "does this place has wheelchair accessible parking?",
				},
			}}
		_, err := surveyRepo.Create(&survey2)
		assert.NoError(t, err)
		entries := surveyRepo.Entries()
		assert.Equal(t, map[ksuid.KSUID]models.Survey{survey1.ID: survey1, survey2.ID: survey2}, entries)
	})
}
