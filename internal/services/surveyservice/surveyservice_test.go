package surveyservice

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"survey-platform/internal/repositories/repositories_mock"
	"survey-platform/pkg/idgenerator/idgenerator_mock"
	"survey-platform/pkg/timegenerator/timegenerator_mock"
	"testing"
	"time"
)

func TestSurveyService_CreateSurvey(t *testing.T) {
	t.Run("should call repo and successfully create survey", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		idGeneratorMock := idgenerator_mock.NewMockIDGenerator(ctrl)
		gomock.InOrder(
			idGeneratorMock.EXPECT().Generate().Return(surveyID),
			idGeneratorMock.EXPECT().Generate().Return(qID1),
			idGeneratorMock.EXPECT().Generate().Return(qID2),
		)
		now := time.Now()
		survey := models.Survey{
			ID:        surveyID,
			Name:      "new survey",
			CreatedAt: now,
			UpdatedAt: now,
			Questions: []models.Question{
				{
					ID:       qID1,
					Question: "is this place good?",
				},
				{
					ID:       qID2,
					Question: "does this place has parking?",
				},
			}}
		mockSurveyRepo.EXPECT().Create(&survey).Return(&survey, nil)
		timeGeneratorMock := timegenerator_mock.NewMockTimeGenInterface(ctrl)
		timeGeneratorMock.EXPECT().Now().Return(now)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, idGeneratorMock, timeGeneratorMock)
		createdSurvey, err := surveyService.CreateSurvey(&survey)
		assert.NoError(t, err)
		assert.Equal(t, survey, *createdSurvey)
	})

	t.Run("should return error when repo layer returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		now := time.Now()
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		idGeneratorMock := idgenerator_mock.NewMockIDGenerator(ctrl)
		gomock.InOrder(
			idGeneratorMock.EXPECT().Generate().Return(surveyID),
			idGeneratorMock.EXPECT().Generate().Return(qID1),
			idGeneratorMock.EXPECT().Generate().Return(qID2),
		)
		survey := models.Survey{
			Name:      "new survey",
			ID:        surveyID,
			CreatedAt: now,
			UpdatedAt: now,
			Questions: []models.Question{
				{
					ID:       qID1,
					Question: "is this place good?",
				},
				{
					ID:       qID2,
					Question: "does this place has parking?",
				},
			}}
		timeGeneratorMock := timegenerator_mock.NewMockTimeGenInterface(ctrl)
		timeGeneratorMock.EXPECT().Now().Return(now)
		mockSurveyRepo.EXPECT().Create(&survey).Return(nil, errors.New("something went wrong"))
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, idGeneratorMock, timeGeneratorMock)
		createdSurvey, err := surveyService.CreateSurvey(&survey)
		assert.Error(t, err)
		assert.Nil(t, createdSurvey)
	})

	t.Run("should return error when number of questions is greater than max questions", func(t *testing.T) {
		survey := models.Survey{
			Name:      "new survey",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
				{
					Question: "does this place serve coffee?",
				},
				{
					Question: "does this place offer dine in?",
				},
			}}
		surveyService := NewSurveyService(3, nil, nil, nil, nil)
		createdSurvey, err := surveyService.CreateSurvey(&survey)
		assert.Error(t, err)
		assert.Nil(t, createdSurvey)
	})

	t.Run("should return error when there are no questions", func(t *testing.T) {
		survey := models.Survey{
			Name:      "new survey",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{}}
		surveyService := NewSurveyService(3, nil, nil, nil, nil)
		createdSurvey, err := surveyService.CreateSurvey(&survey)
		assert.Error(t, err)
		assert.Nil(t, createdSurvey)
	})

	t.Run("should return error when survey name is empty", func(t *testing.T) {
		survey := models.Survey{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
				{
					Question: "does this place serve coffee?",
				},
			}}
		surveyService := NewSurveyService(3, nil, nil, nil, nil)
		createdSurvey, err := surveyService.CreateSurvey(&survey)
		assert.Error(t, err)
		assert.Nil(t, createdSurvey)
	})

	t.Run("should add survey id, question id and timestamps while creating survey", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		now := time.Now()
		surveyID, qID1, qID2, qID3 := ksuid.New(), ksuid.New(), ksuid.New(), ksuid.New()
		idGeneratorMock := idgenerator_mock.NewMockIDGenerator(ctrl)
		gomock.InOrder(
			idGeneratorMock.EXPECT().Generate().Return(surveyID),
			idGeneratorMock.EXPECT().Generate().Return(qID1),
			idGeneratorMock.EXPECT().Generate().Return(qID2),
			idGeneratorMock.EXPECT().Generate().Return(qID3),
		)
		timeGeneratorMock := timegenerator_mock.NewMockTimeGenInterface(ctrl)
		timeGeneratorMock.EXPECT().Now().Return(now)
		survey := models.Survey{
			Name:      "new survey",
			ID:        surveyID,
			CreatedAt: now,
			UpdatedAt: now,
			Questions: []models.Question{
				{
					ID:       qID1,
					Question: "is this place good?",
				},
				{
					ID:       qID2,
					Question: "does this place has parking?",
				},
				{

					Question: "does this place serve coffee?",
				},
			}}
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().Create(&survey).Return(&survey, nil)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, idGeneratorMock, timeGeneratorMock)
		createdSurvey, err := surveyService.CreateSurvey(&survey)
		assert.NoError(t, err)
		assert.NotNil(t, createdSurvey.ID)
		assert.NotNil(t, createdSurvey.Questions[0].ID)
		assert.NotZero(t, createdSurvey.CreatedAt)
		assert.NotZero(t, createdSurvey.UpdatedAt)
	})
}

func TestSurveyService_GetSurvey(t *testing.T) {
	t.Run("should call survey repo and return survey", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		survey := models.Survey{
			ID:        ksuid.KSUID{},
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
			Name:      "global survey",
			Questions: []models.Question{
				{
					Question: "is this good?",
					ID:       ksuid.New(),
				},
				{
					Question: "are you liking this product?",
					ID:       ksuid.New(),
				},
			},
		}
		mockSurveyRepo.EXPECT().Get(survey.ID).Return(&survey, nil)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, nil, nil)
		returnedSurvey, err := surveyService.GetSurvey(survey.ID)
		assert.NoError(t, err)
		assert.Equal(t, survey, *returnedSurvey)
	})
	t.Run("should return error when repo layer returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		surveyID := ksuid.New()
		mockSurveyRepo.EXPECT().Get(surveyID).Return(nil, errors.New("something went wrong"))
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, nil, nil)
		returnedSurvey, err := surveyService.GetSurvey(surveyID)
		assert.Error(t, err)
		assert.Nil(t, returnedSurvey)
	})
}

func TestSurveyService_UpdateSurvey(t *testing.T) {
	t.Run("should successfully update survey", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		q2ID := ksuid.New()
		now := time.Now()
		timeGeneratorMock := timegenerator_mock.NewMockTimeGenInterface(ctrl)
		timeGeneratorMock.EXPECT().Now().Return(now)
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: now,
			UpdatedAt: now,
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
		oldSurvey := survey
		oldSurvey.UpdatedAt = now.AddDate(0, 0, -1)
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().Update(surveyID, &survey).Return(&survey, nil)
		idGeneratorMock := idgenerator_mock.NewMockIDGenerator(ctrl)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, idGeneratorMock, timeGeneratorMock)
		updatedSurvey, err := surveyService.UpdateSurvey(surveyID, oldSurvey)
		assert.NoError(t, err)
		assert.NotEqual(t, oldSurvey.UpdatedAt, updatedSurvey.UpdatedAt)
	})
	t.Run("should successfully update survey with new question", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		now := time.Now()
		timeGeneratorMock := timegenerator_mock.NewMockTimeGenInterface(ctrl)
		timeGeneratorMock.EXPECT().Now().Return(now)
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: now,
			UpdatedAt: now,
			Questions: []models.Question{
				{
					ID:       q1ID,
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
			},
		}
		oldSurvey := survey
		oldSurvey.UpdatedAt = now.AddDate(0, 0, -1)
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().Update(surveyID, &survey).Return(&survey, nil)
		qID2 := ksuid.New()
		idGeneratorMock := idgenerator_mock.NewMockIDGenerator(ctrl)
		idGeneratorMock.EXPECT().Generate().Return(qID2)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, idGeneratorMock, timeGeneratorMock)
		updatedSurvey, err := surveyService.UpdateSurvey(surveyID, oldSurvey)
		assert.NoError(t, err)
		assert.NotEqual(t, oldSurvey.UpdatedAt, updatedSurvey.UpdatedAt)
		assert.False(t, updatedSurvey.Questions[1].ID.IsNil())
	})
	t.Run("should return error when survey repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		q1ID := ksuid.New()
		q2ID := ksuid.New()
		now := time.Now()
		timeGeneratorMock := timegenerator_mock.NewMockTimeGenInterface(ctrl)
		timeGeneratorMock.EXPECT().Now().Return(now)
		survey := models.Survey{
			ID:        surveyID,
			CreatedAt: now,
			UpdatedAt: now,
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
		oldSurvey := survey
		oldSurvey.UpdatedAt = now.AddDate(0, 0, -1)
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().Update(surveyID, &survey).Return(&survey, nil)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, nil, timeGeneratorMock)
		updatedSurvey, err := surveyService.UpdateSurvey(surveyID, oldSurvey)
		assert.NoError(t, err)
		assert.NotEqual(t, oldSurvey.UpdatedAt, updatedSurvey.UpdatedAt)
	})
}

func TestSurveyService_DeleteSurvey(t *testing.T) {
	t.Run("should successfully delete survey", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		surveyID := ksuid.New()
		mockSurveyRepo.EXPECT().Delete(surveyID).Return(nil)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, nil, nil)
		err := surveyService.DeleteSurvey(surveyID)
		assert.NoError(t, err)
	})
	t.Run("should return error when repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		surveyID := ksuid.New()
		mockSurveyRepo.EXPECT().Delete(surveyID).Return(repositories.ErrNotFound)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, nil, nil)
		err := surveyService.DeleteSurvey(surveyID)
		assert.Equal(t, repositories.ErrNotFound, err)
	})
}

func TestSurveyService_GetAllSurveys(t *testing.T) {
	t.Run("should successfully return all surveys", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveys := []models.Survey{{
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
			},
		}, {
			ID:        ksuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Questions: []models.Question{
				{
					ID:       ksuid.New(),
					Question: "is apple m1 chip good?",
				},
				{
					ID:       ksuid.New(),
					Question: "do you prefer linux over MacOS",
				},
			},
		}}
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().GetAll().Return(mockSurveys, nil)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, nil, nil)
		surveys, err := surveyService.GetAllSurveys()
		assert.NoError(t, err)
		assert.Equal(t, mockSurveys, surveys)
	})
	t.Run("should return error when repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().GetAll().Return(nil, repositories.ErrNotFound)
		surveyService := NewSurveyService(3, mockSurveyRepo, nil, nil, nil)
		surveys, err := surveyService.GetAllSurveys()
		assert.Equal(t, repositories.ErrNotFound, err)
		assert.Nil(t, surveys)
	})
}

func TestSurveyService_SaveResponse(t *testing.T) {
	t.Run("should successfully save response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID, qID1, qID2, responseID := ksuid.New(), ksuid.New(), ksuid.New(), ksuid.New()
		now := time.Now()
		timeGeneratorMock := timegenerator_mock.NewMockTimeGenInterface(ctrl)
		timeGeneratorMock.EXPECT().Now().Return(now)
		mockResponse := models.Response{
			ID:        responseID,
			SurveyID:  surveyID,
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
		}
		mockResponseRepo := repositories_mock.NewMockResponseRepoInterface(ctrl)
		mockResponseRepo.EXPECT().Create(&mockResponse).Return(&mockResponse, nil)
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().Get(surveyID).Return(&models.Survey{}, nil)
		mockIDGenerator := idgenerator_mock.NewMockIDGenerator(ctrl)
		mockIDGenerator.EXPECT().Generate().Return(responseID)
		surveyService := NewSurveyService(3, mockSurveyRepo, mockResponseRepo, mockIDGenerator, timeGeneratorMock)
		response, err := surveyService.SaveResponse(mockResponse)
		assert.NoError(t, err)
		assert.Equal(t, mockResponse, *response)
	})
	t.Run("should return error when survey is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		mockResponse := models.Response{
			SurveyID:  surveyID,
			CreatedAt: time.Now(),
			Answers: []models.Answer{
				{
					QuestionID: qID1,
					Answer:     false,
				},
				{
					QuestionID: qID2,
					Answer:     true,
				},
				{
					QuestionID: qID1,
					Answer:     false,
				},
				{
					QuestionID: qID2,
					Answer:     true,
				},
			},
		}
		mockResponseRepo := repositories_mock.NewMockResponseRepoInterface(ctrl)
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().Get(surveyID).Return(nil, repositories.ErrNotFound)
		surveyService := NewSurveyService(3, mockSurveyRepo, mockResponseRepo, nil, nil)
		response, err := surveyService.SaveResponse(mockResponse)
		assert.Equal(t, repositories.ErrNotFound, err)
		assert.Nil(t, response)
	})
	t.Run("should return error when answers are more than allowed questions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		mockResponse := models.Response{
			ID:        ksuid.New(),
			SurveyID:  surveyID,
			CreatedAt: time.Now(),
			Answers: []models.Answer{
				{
					QuestionID: qID1,
					Answer:     false,
				},
				{
					QuestionID: qID2,
					Answer:     true,
				},
				{
					QuestionID: qID1,
					Answer:     false,
				},
				{
					QuestionID: qID2,
					Answer:     true,
				},
			},
		}
		mockResponseRepo := repositories_mock.NewMockResponseRepoInterface(ctrl)
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().Get(surveyID).Return(&models.Survey{}, nil)
		surveyService := NewSurveyService(3, mockSurveyRepo, mockResponseRepo, nil, nil)
		response, err := surveyService.SaveResponse(mockResponse)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestSurveyService_GetResponses(t *testing.T) {
	t.Run("should successfully get responses for a survey", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockResponseRepo := repositories_mock.NewMockResponseRepoInterface(ctrl)
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		mockResponses := []models.Response{
			{
				ID:        ksuid.New(),
				SurveyID:  surveyID,
				CreatedAt: time.Now(),
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
		mockResponseRepo.EXPECT().GetBySurveyID(surveyID).Return(mockResponses, nil)
		surveyService := NewSurveyService(3, nil, mockResponseRepo, nil, nil)
		responses, err := surveyService.GetResponses(surveyID)
		assert.NoError(t, err)
		assert.Equal(t, mockResponses, responses)
	})
	t.Run("should return error when repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockResponseRepo := repositories_mock.NewMockResponseRepoInterface(ctrl)
		surveyID := ksuid.New()
		mockResponseRepo.EXPECT().GetBySurveyID(surveyID).Return(nil, repositories.ErrNotFound)

		surveyService := NewSurveyService(3, nil, mockResponseRepo, nil, nil)
		responses, err := surveyService.GetResponses(surveyID)
		assert.Equal(t, repositories.ErrNotFound, err)
		assert.Nil(t, responses)
	})
}

func TestSurveyService_Dump(t *testing.T) {
	t.Run("should combine the entries from both survey and response repo and return", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		surveyEntries := map[ksuid.KSUID]models.Survey{
			surveyID: {
				ID: surveyID,
				Questions: []models.Question{
					{ID: qID1, Question: "is this a good product?"},
					{ID: qID2, Question: "would you recommend this product to a friend?"},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      "new survey",
			},
		}
		responseEntries := map[ksuid.KSUID][]models.Response{
			surveyID: {
				{
					ID:        ksuid.New(),
					SurveyID:  surveyID,
					CreatedAt: time.Now(),
					Answers: []models.Answer{
						{QuestionID: qID1, Answer: true},
						{QuestionID: qID2, Answer: false},
					},
				},
			},
		}
		mockResponseRepo := repositories_mock.NewMockResponseRepoInterface(ctrl)
		mockSurveyRepo := repositories_mock.NewMockSurveyRepoInterface(ctrl)
		mockSurveyRepo.EXPECT().Entries().Return(surveyEntries)
		mockResponseRepo.EXPECT().Entries().Return(responseEntries)
		surveyService := NewSurveyService(3, mockSurveyRepo, mockResponseRepo, nil, nil)
		repoEntries := surveyService.Entries()
		assert.Equal(t, models.DBEntry{Surveys: surveyEntries, Responses: responseEntries}, *repoEntries)
	})
}
