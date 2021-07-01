package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"survey-platform/internal/services/services_mock"
	"testing"
)

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
		mockService.EXPECT().GetSurvey(surveyID).Return(nil, errors.New("error not found"))
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
}

func TestSurveyApp_DeleteSurvey(t *testing.T) {

}

func TestSurveyApp_GetAllSurveys(t *testing.T) {

}

func TestSurveyApp_SaveResponse(t *testing.T) {

}

func TestSurveyApp_GetResponses(t *testing.T) {

}

func TestSurveyApp_Dump(t *testing.T) {

}

// helper function
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}
