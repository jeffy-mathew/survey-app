// Code generated by MockGen. DO NOT EDIT.
// Source: services.go

// Package services_mock is a generated GoMock package.
package services_mock

import (
	reflect "reflect"
	models "survey-platform/internal/models"

	gomock "github.com/golang/mock/gomock"
	ksuid "github.com/segmentio/ksuid"
)

// MockSurveyServiceInterface is a mock of SurveyServiceInterface interface.
type MockSurveyServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSurveyServiceInterfaceMockRecorder
}

// MockSurveyServiceInterfaceMockRecorder is the mock recorder for MockSurveyServiceInterface.
type MockSurveyServiceInterfaceMockRecorder struct {
	mock *MockSurveyServiceInterface
}

// NewMockSurveyServiceInterface creates a new mock instance.
func NewMockSurveyServiceInterface(ctrl *gomock.Controller) *MockSurveyServiceInterface {
	mock := &MockSurveyServiceInterface{ctrl: ctrl}
	mock.recorder = &MockSurveyServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSurveyServiceInterface) EXPECT() *MockSurveyServiceInterfaceMockRecorder {
	return m.recorder
}

// CreateSurvey mocks base method.
func (m *MockSurveyServiceInterface) CreateSurvey(survey *models.Survey) (*models.Survey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSurvey", survey)
	ret0, _ := ret[0].(*models.Survey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSurvey indicates an expected call of CreateSurvey.
func (mr *MockSurveyServiceInterfaceMockRecorder) CreateSurvey(survey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSurvey", reflect.TypeOf((*MockSurveyServiceInterface)(nil).CreateSurvey), survey)
}

// DeleteSurvey mocks base method.
func (m *MockSurveyServiceInterface) DeleteSurvey(id ksuid.KSUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSurvey", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSurvey indicates an expected call of DeleteSurvey.
func (mr *MockSurveyServiceInterfaceMockRecorder) DeleteSurvey(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSurvey", reflect.TypeOf((*MockSurveyServiceInterface)(nil).DeleteSurvey), id)
}

// Entries mocks base method.
func (m *MockSurveyServiceInterface) Entries() *models.DBEntry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Entries")
	ret0, _ := ret[0].(*models.DBEntry)
	return ret0
}

// Entries indicates an expected call of Entries.
func (mr *MockSurveyServiceInterfaceMockRecorder) Entries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Entries", reflect.TypeOf((*MockSurveyServiceInterface)(nil).Entries))
}

// GetAllSurveys mocks base method.
func (m *MockSurveyServiceInterface) GetAllSurveys() ([]models.Survey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSurveys")
	ret0, _ := ret[0].([]models.Survey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSurveys indicates an expected call of GetAllSurveys.
func (mr *MockSurveyServiceInterfaceMockRecorder) GetAllSurveys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSurveys", reflect.TypeOf((*MockSurveyServiceInterface)(nil).GetAllSurveys))
}

// GetResponses mocks base method.
func (m *MockSurveyServiceInterface) GetResponses(surveyID ksuid.KSUID) ([]models.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResponses", surveyID)
	ret0, _ := ret[0].([]models.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResponses indicates an expected call of GetResponses.
func (mr *MockSurveyServiceInterfaceMockRecorder) GetResponses(surveyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResponses", reflect.TypeOf((*MockSurveyServiceInterface)(nil).GetResponses), surveyID)
}

// GetSurvey mocks base method.
func (m *MockSurveyServiceInterface) GetSurvey(id ksuid.KSUID) (*models.Survey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSurvey", id)
	ret0, _ := ret[0].(*models.Survey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSurvey indicates an expected call of GetSurvey.
func (mr *MockSurveyServiceInterfaceMockRecorder) GetSurvey(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSurvey", reflect.TypeOf((*MockSurveyServiceInterface)(nil).GetSurvey), id)
}

// SaveResponse mocks base method.
func (m *MockSurveyServiceInterface) SaveResponse(response models.Response) (*models.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveResponse", response)
	ret0, _ := ret[0].(*models.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveResponse indicates an expected call of SaveResponse.
func (mr *MockSurveyServiceInterfaceMockRecorder) SaveResponse(response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveResponse", reflect.TypeOf((*MockSurveyServiceInterface)(nil).SaveResponse), response)
}

// UpdateSurvey mocks base method.
func (m *MockSurveyServiceInterface) UpdateSurvey(id ksuid.KSUID, survey models.Survey) (*models.Survey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSurvey", id, survey)
	ret0, _ := ret[0].(*models.Survey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSurvey indicates an expected call of UpdateSurvey.
func (mr *MockSurveyServiceInterfaceMockRecorder) UpdateSurvey(id, survey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSurvey", reflect.TypeOf((*MockSurveyServiceInterface)(nil).UpdateSurvey), id, survey)
}
