package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"survey-platform/internal/db/db_mock"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"survey-platform/internal/services/services_mock"
	"testing"
	"time"
)

func TestSurveyApp_HealthCheck(t *testing.T) {
	t.Run("should return status ok(200) on hitting health check endpoint", func(t *testing.T) {
		surveyApp := NewSurveyApp(nil, nil)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestSurveyApp_CreateSurvey(t *testing.T) {
	t.Run("should successfully create survey and return status created (201)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurvey := models.Survey{
			Name: "new survey",
			Questions: []models.Question{
				{
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
			},
		}
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().CreateSurvey(mockSurvey).Return(&mockSurvey, nil)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		marshalledSurvey, _ := json.Marshal(&mockSurvey)
		body := bytes.NewReader(marshalledSurvey)
		req, _ := http.NewRequest(http.MethodPost, "/survey/", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})
	t.Run("should return internal server error(500) when service returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurvey := models.Survey{
			Name: "new survey",
			Questions: []models.Question{
				{
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
			},
		}
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().CreateSurvey(mockSurvey).Return(nil, errors.New("something went wrong"))
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		marshalledSurvey, _ := json.Marshal(&mockSurvey)
		body := bytes.NewReader(marshalledSurvey)
		req, _ := http.NewRequest(http.MethodPost, "/survey/", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
	t.Run("should return unprocessable entity(422) with invalid json input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		body := bytes.NewReader([]byte("hello"))
		req, _ := http.NewRequest(http.MethodPost, "/survey/", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})
}

func TestSurveyApp_GetSurvey(t *testing.T) {
	t.Run("should return statusOK(200) on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockSurvey := models.Survey{
			Name: "new survey",
			Questions: []models.Question{
				{
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
			},
		}
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().GetSurvey(surveyID).Return(&mockSurvey, nil)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/survey/%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("should return statusUnprocessableEntity(422) when id is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, "/survey/xxx", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})
	t.Run("should return statusNotFound(404) when given surveyID is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().GetSurvey(surveyID).Return(nil, repositories.ErrNotFound)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/survey/%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
	t.Run("should return StatusInternalServerError(500) when service return other errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().GetSurvey(surveyID).Return(nil, errors.New("something went wrong"))
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/survey/%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSurveyApp_UpdateSurvey(t *testing.T) {
	t.Run("should return statusOK(200) on successful update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockSurvey := models.Survey{
			ID:   surveyID,
			Name: "updated survey",
			Questions: []models.Question{
				{
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
			},
		}
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().UpdateSurvey(surveyID, mockSurvey).Return(&mockSurvey, nil)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		marshalledSurvey, _ := json.Marshal(&mockSurvey)
		body := bytes.NewReader(marshalledSurvey)
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/survey/%s", surveyID.String()), body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("should return statusUnprocessableEntity(422) when id is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodPut, "/survey/xxx", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})
	t.Run("should return statusUnprocessableEntity(422) when body in malformed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		body := bytes.NewReader([]byte("hello"))
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/survey/%s", surveyID.String()), body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})
	t.Run("should return statusNotFound(404) when given surveyID is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockSurvey := models.Survey{
			ID:   surveyID,
			Name: "updated survey",
			Questions: []models.Question{
				{
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
			},
		}
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().UpdateSurvey(surveyID, mockSurvey).Return(nil, repositories.ErrNotFound)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		marshalledSurvey, _ := json.Marshal(&mockSurvey)
		body := bytes.NewReader(marshalledSurvey)
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/survey/%s", surveyID.String()), body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
	t.Run("should return StatusInternalServerError(500) when error on service layer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockSurvey := models.Survey{
			ID:   surveyID,
			Name: "updated survey",
			Questions: []models.Question{
				{
					Question: "is this place good?",
				},
				{
					Question: "does this place has parking?",
				},
			},
		}
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().UpdateSurvey(surveyID, mockSurvey).Return(nil, errors.New("something went wrong"))
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		marshalledSurvey, _ := json.Marshal(&mockSurvey)
		body := bytes.NewReader(marshalledSurvey)
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/survey/%s", surveyID.String()), body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSurveyApp_DeleteSurvey(t *testing.T) {
	t.Run("should return statusNoContent(204) on successful deletion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()

		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().DeleteSurvey(surveyID).Return(nil)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/survey/%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})
	t.Run("should return statusUnprocessableEntity(422) when id is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodDelete, "/survey/xxxx", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})
	t.Run("should return statusNotFound(404) when given surveyID is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().DeleteSurvey(surveyID).Return(repositories.ErrNotFound)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/survey/%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
	t.Run("should return StatusInternalServerError(500) when service returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().DeleteSurvey(surveyID).Return(errors.New("something went wrong"))
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/survey/%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSurveyApp_GetAllSurveys(t *testing.T) {
	t.Run("should return statusOK(200) on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSurveys := []models.Survey{
			{
				Name: "new survey",
				Questions: []models.Question{
					{
						Question: "is this place good?",
					},
					{
						Question: "does this place has parking?",
					},
				},
			},
			{
				Name: "good survey",
				Questions: []models.Question{
					{
						Question: "is this place good?",
					},
					{
						Question: "does this place has parking?",
					},
				},
			},
		}
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().GetAllSurveys().Return(mockSurveys, nil)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, "/survey/", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("should return StatusInternalServerError(500) when service return other errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().GetAllSurveys().Return(nil, errors.New("something went wrong"))
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, "/survey/", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSurveyApp_SaveResponse(t *testing.T) {
	t.Run("should successfully save response and return status created (201)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		mockResponse := models.Response{
			SurveyID: surveyID,
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
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().SaveResponse(mockResponse).Return(&mockResponse, nil)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		marshalledSurvey, _ := json.Marshal(&mockResponse)
		body := bytes.NewReader(marshalledSurvey)
		req, _ := http.NewRequest(http.MethodPost, "/response/", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})
	t.Run("should not save response and return statusNotFound(404) when surveyID is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		mockResponse := models.Response{
			SurveyID: surveyID,
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
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().SaveResponse(mockResponse).Return(nil, repositories.ErrNotFound)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		marshalledSurvey, _ := json.Marshal(&mockResponse)
		body := bytes.NewReader(marshalledSurvey)
		req, _ := http.NewRequest(http.MethodPost, "/response/", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
	t.Run("should return internal server error(500) when service returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		mockResponse := models.Response{
			SurveyID: surveyID,
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
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().SaveResponse(mockResponse).Return(nil, errors.New("something went wrong"))
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		marshalledSurvey, _ := json.Marshal(&mockResponse)
		body := bytes.NewReader(marshalledSurvey)
		req, _ := http.NewRequest(http.MethodPost, "/response/", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
	t.Run("should return unprocessable entity(422) with invalid json input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		body := bytes.NewReader([]byte("hello"))
		req, _ := http.NewRequest(http.MethodPost, "/response/", body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})
}

func TestSurveyApp_GetResponses(t *testing.T) {
	t.Run("should return statusOK(200) on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		mockResponses := []models.Response{
			{
				SurveyID: surveyID,
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
			{
				SurveyID: surveyID,
				Answers: []models.Answer{
					{
						QuestionID: qID1,
						Answer:     true,
					},
					{
						QuestionID: qID2,
						Answer:     true,
					},
				},
			},
		}
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().GetResponses(surveyID).Return(mockResponses, nil)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/response/?survey_id=%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("should return statusUnprocessableEntity(422) when survey id is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, "/response/?survey_id=xxx", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})
	t.Run("should return statusNotFound(404) when given surveyID has no responses", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().GetResponses(surveyID).Return(nil, repositories.ErrNotFound)
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/response/?survey_id=%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
	t.Run("should return StatusInternalServerError(500) when service return other errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		surveyID := ksuid.New()
		mockService := services_mock.NewMockSurveyServiceInterface(ctrl)
		mockService.EXPECT().GetResponses(surveyID).Return(nil, errors.New("something went wrong"))
		surveyApp := NewSurveyApp(nil, mockService)
		router := surveyApp.SetupRoutes()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/response/?survey_id=%s", surveyID.String()), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSurveyApp_Dump(t *testing.T) {
	t.Run("should return services entries on dump", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockDB := db_mock.NewMockDB(ctrl)
		mockSurveyService := services_mock.NewMockSurveyServiceInterface(ctrl)
		surveyID, qID1, qID2 := ksuid.New(), ksuid.New(), ksuid.New()
		dbEntry := models.DBEntry{
			Responses: map[ksuid.KSUID][]models.Response{
				surveyID: []models.Response{
					{
						SurveyID: surveyID,
						ID:       ksuid.New(),
						Answers: []models.Answer{
							{
								QuestionID: qID1,
								Answer:     true,
							},
							{
								QuestionID: qID2,
								Answer:     true,
							},
						},
					},
				},
			},
			Surveys: map[ksuid.KSUID]models.Survey{
				surveyID: models.Survey{
					ID:   surveyID,
					Name: "new survey",
					Questions: []models.Question{
						{
							ID:       qID1,
							Question: "is this a good product?",
						},
						{
							ID:       qID2,
							Question: "would you recommend this to your friend?",
						},
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}
		mockSurveyService.EXPECT().Entries().Return(&dbEntry)
		mockDB.EXPECT().Dump(&dbEntry).Return(nil)
		surveyApp := NewSurveyApp(mockDB, mockSurveyService)
		err := surveyApp.Dump()
		assert.NoError(t, err)
	})
}
