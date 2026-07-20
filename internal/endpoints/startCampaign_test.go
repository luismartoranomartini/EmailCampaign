package endpoints

import (
	"errors"
	"net/http"
	"net/http/httptest"
	internalmock "projeto-golang/internal/test/internalMock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignsStart_Return200(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	service.On("Start", mock.Anything).Return(nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("PATCH", "/", nil)
	rr := httptest.NewRecorder()

	_, status, err := handler.CampaignStart(rr, req)

	assert.Equal(200, status)
	assert.Nil(err)
}
func Test_CampaignsStart_ReturnErr(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	errExpected := errors.New("Something wrong")
	service.On("Start", mock.Anything).Return(errExpected)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("PATCH", "/", nil)
	rr := httptest.NewRecorder()

	_, _, err := handler.CampaignStart(rr, req)

	assert.Equal(errExpected, err)
}
