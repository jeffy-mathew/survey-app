package app

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
	"net/http"
	_ "survey-platform/docs"
	"survey-platform/internal/db"
	"survey-platform/internal/models"
	"survey-platform/internal/repositories"
	"survey-platform/internal/services"
)

var ApiVersion = "1.0.0"

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

// @title Survey app API
// @version 1.0
// @description maintains survey CRUD
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

func (a *SurveyApp) SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/", a.HealthCheck)
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
	router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))
	return router
}

// HealthCheck godoc
// @Summary Check app health
// @Description check app health by hitting at root
// @Accept  json
// @Produce  json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router / [get]
func (a *SurveyApp) HealthCheck(c *gin.Context) {
	c.JSONP(http.StatusOK, Response{Message: "service is up", ApiVersion: ApiVersion})
}

// CreateSurvey godoc
// @Summary creates survey
// @Description creates a survey
// @Accept  json
// @Produce  json
// @Param survey body models.Survey true "survey"
// @success 201 {object} Response{data=models.Survey} "desc"
// @Failure 500 {object} Response
// @Failure 422 {object} Response
// @Router /survey/ [post]
func (a *SurveyApp) CreateSurvey(c *gin.Context) {
	var survey models.Survey
	err := c.ShouldBindJSON(&survey)
	if err != nil {
		log.Println("error while reading survey body", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "malformed body", ApiVersion: ApiVersion})
		return
	}
	newSurvey, err := a.surveyService.CreateSurvey(&survey)
	if err != nil {
		log.Println("error while reading survey body", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while creating survey " + err.Error(), ApiVersion: ApiVersion})
		return
	}
	c.JSONP(http.StatusCreated, Response{Message: "survey created", Data: newSurvey, ApiVersion: ApiVersion})
}

func (a *SurveyApp) GetSurvey(c *gin.Context) {
	id, err := ksuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error while parsing surveyID", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "invalid survey id", ApiVersion: ApiVersion})
		return
	}
	survey, err := a.surveyService.GetSurvey(id)
	if err != nil && err == repositories.ErrNotFound {
		log.Println("survey not found while getting survey", id.String())
		c.JSONP(http.StatusNotFound, Response{Message: "error while reading survey " + err.Error(), ApiVersion: ApiVersion})
		return
	} else if err != nil {
		log.Println("error while getting survey", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while reading survey " + err.Error(), ApiVersion: ApiVersion})
		return
	}
	c.JSONP(http.StatusOK, Response{Message: "success", Data: survey})
}

func (a *SurveyApp) UpdateSurvey(c *gin.Context) {
	id, err := ksuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error while parsing surveyID", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "invalid survey id", ApiVersion: ApiVersion})
		return
	}
	var survey models.Survey
	err = c.ShouldBindJSON(&survey)
	if err != nil {
		log.Println("error while reading survey body", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "malformed body", ApiVersion: ApiVersion})
		return
	}
	updatedSurvey, err := a.surveyService.UpdateSurvey(id, survey)
	if err != nil && err == repositories.ErrNotFound {
		log.Println("survey not found while updating survey", id.String())
		c.JSONP(http.StatusNotFound, Response{Message: "error while updating survey " + err.Error(), ApiVersion: ApiVersion})
		return
	} else if err != nil {
		log.Println("error while updating survey", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while updating survey " + err.Error(), ApiVersion: ApiVersion})
		return
	}
	c.JSONP(http.StatusOK, Response{Message: "survey updated", Data: updatedSurvey, ApiVersion: ApiVersion})
}

func (a *SurveyApp) DeleteSurvey(c *gin.Context) {
	id, err := ksuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("error while parsing surveyID", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "invalid survey id", ApiVersion: ApiVersion})
		return
	}
	err = a.surveyService.DeleteSurvey(id)
	if err != nil && err == repositories.ErrNotFound {
		log.Println("survey not found while deleting survey", id.String())
		c.JSONP(http.StatusNotFound, Response{Message: "error while deleting survey " + err.Error(), ApiVersion: ApiVersion})
		return
	} else if err != nil {
		log.Println("error while getting survey", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while deleting survey " + err.Error(), ApiVersion: ApiVersion})
		return
	}
	c.JSONP(http.StatusNoContent, Response{Message: "survey deleted", ApiVersion: ApiVersion})
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
	if err != nil && err == repositories.ErrNotFound {
		log.Println("survey not found while saving response", response.SurveyID.String())
		c.JSONP(http.StatusNotFound, Response{Message: "error while saving response " + err.Error(), ApiVersion: ApiVersion})
		return
	} else if err != nil {
		log.Println("error while saving survey response", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while saving survey response " + err.Error(), ApiVersion: ApiVersion})
		return
	}
	c.JSONP(http.StatusCreated, Response{Message: "saved response", ApiVersion: ApiVersion})
}

func (a *SurveyApp) GetResponses(c *gin.Context) {
	id, err := ksuid.Parse(c.Query("survey_id"))
	if err != nil {
		log.Println("error while parsing surveyID", err)
		c.JSONP(http.StatusUnprocessableEntity, Response{Message: "invalid survey id", ApiVersion: ApiVersion})
		return
	}
	responses, err := a.surveyService.GetResponses(id)
	if err != nil && err == repositories.ErrNotFound {
		log.Println("survey not found while fetching responses", id.String())
		c.JSONP(http.StatusNotFound, Response{Message: "error while fetching responses " + err.Error(), ApiVersion: ApiVersion})
		return
	} else if err != nil {
		log.Println("error while fetching responses for survey", err)
		c.JSONP(http.StatusInternalServerError, Response{Message: "error while fetching responses", ApiVersion: ApiVersion})
		return
	}
	c.JSONP(http.StatusOK, Response{Message: "success", Data: responses})
}

func (a *SurveyApp) Dump() error {
	return a.db.Dump(a.surveyService.Entries())
}
