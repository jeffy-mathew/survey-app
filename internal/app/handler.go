package app

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
	"survey-platform/internal/db"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"survey-platform/internal/services"
)

type Response struct {
	Message    string      `json:"success,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	ApiVersion string      `json:"api_version,omitempty"`
}

// SurveyApp handles the hit and dump from high level
type SurveyApp struct {
	db            db.DB
	surveyService services.SurveyServiceInterface
}

// NewSurveyApp returns app configured with passed surveyService
func NewSurveyApp(persistence db.DB, surveyService services.SurveyServiceInterface) *SurveyApp {
	return &SurveyApp{
		db:            persistence,
		surveyService: surveyService,
	}
}

func (a *SurveyApp) SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSONP(http.StatusCreated, Response{Message: "service is up"})
	})
	surveyRouter := router.Group("/survey")
	{
		surveyRouter.GET("/", a.GetAllSurveys)
		surveyRouter.POST("/", a.CreateSurvey)
		surveyRouter.GET("/:id", a.GetSurvey)
		surveyRouter.PUT("/:id", a.UpdateSurvey)
		surveyRouter.DELETE("/:id", a.DeleteSurvey)
	}
	responseRouter := router.Group("/response")
	{
		responseRouter.POST("/", a.SaveResponse)
		responseRouter.GET("/", a.GetResponses)
	}
	return router
}

// CreateSurvey is the http handler function for handling the request
func (a *SurveyApp) CreateSurvey(c *gin.Context) {
	var survey models.Survey
	err := c.ShouldBindJSON(&survey)
	if err != nil {
		log.Println("error while reading survey body", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "malformed body"})
		return
	}
	newSurvey, err := a.surveyService.CreateSurvey(survey)
	if err != nil {
		log.Println("error while reading survey body", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while creating survey " + err.Error()})
		return
	}
	c.JSONP(http.StatusCreated, Response{Message: "survey created", Data: newSurvey})
}

func (a *SurveyApp) GetSurvey(c *gin.Context) {
	id, err := ksuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error while parsing surveyID", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "invalid survey id"})
		return
	}
	survey, err := a.surveyService.GetSurvey(id)
	if err != nil && err == repositories.ErrNotFound {
		log.Println("survey not found while getting survey", id.String())
		c.JSONP(http.StatusNotFound, Response{Message: "error while reading survey " + err.Error()})
		return
	} else if err != nil {
		log.Println("error while getting survey", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while reading survey " + err.Error()})
		return
	}
	c.JSONP(http.StatusOK, Response{Message: "success", Data: survey})
}

func (a *SurveyApp) UpdateSurvey(c *gin.Context) {
	id, err := ksuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error while parsing surveyID", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "invalid survey id"})
		return
	}
	var survey models.Survey
	err = c.ShouldBindJSON(&survey)
	if err != nil {
		log.Println("error while reading survey body", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "malformed body"})
		return
	}
	updatedSurvey, err := a.surveyService.UpdateSurvey(id, survey)
	if err != nil && err == repositories.ErrNotFound {
		log.Println("survey not found while updating survey", id.String())
		c.JSONP(http.StatusNotFound, Response{Message: "error while updating survey " + err.Error()})
		return
	} else if err != nil {
		log.Println("error while updating survey", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while updating survey " + err.Error()})
		return
	}
	c.JSONP(http.StatusOK, Response{Message: "survey updated", Data: updatedSurvey})
}

func (a *SurveyApp) DeleteSurvey(c *gin.Context) {
	id, err := ksuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error while parsing surveyID", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "invalid survey id"})
		return
	}
	err = a.surveyService.DeleteSurvey(id)
	if err != nil && err == repositories.ErrNotFound {
		log.Println("survey not found while deleting survey", id.String())
		c.JSONP(http.StatusNotFound, Response{Message: "error while deleting survey " + err.Error()})
		return
	} else if err != nil {
		log.Println("error while getting survey", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while deleting survey " + err.Error()})
		return
	}
	c.JSONP(http.StatusNoContent, Response{Message: "survey deleted"})
}

func (a *SurveyApp) GetAllSurveys(c *gin.Context) {
	surveys, err := a.surveyService.GetAllSurveys()
	if err != nil {
		log.Println("error while getting surveys", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while reading surveys " + err.Error()})
		return
	}
	c.JSONP(http.StatusOK, Response{Message: "survey created", Data: surveys})
}

func (a *SurveyApp) SaveResponse(c *gin.Context) {
	var response models.Response
	err := c.ShouldBindJSON(&response)
	if err != nil {
		log.Println("error while reading response body", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "malformed body"})
		return
	}
	_, err = a.surveyService.SaveResponse(response)
	if err != nil {
		log.Println("error while saving survey response", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while saving survey response " + err.Error()})
		return
	}
	c.JSONP(http.StatusCreated, Response{Message: "saved response"})
}

func (a *SurveyApp) GetResponses(c *gin.Context) {
	id, err := ksuid.Parse(c.Query("survey_id"))
	if err != nil {
		log.Println("error while parsing surveyID", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "invalid survey id"})
		return
	}
	responses, err := a.surveyService.GetResponses(id)
	if err != nil {
		log.Println("error while getting responses for survey", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while reading responses"})
		return
	}
	c.JSONP(http.StatusOK, Response{Message: "success", Data: responses})
}

func (a *SurveyApp) Dump() error {
	return a.db.Dump(a.surveyService.Entries())
}
