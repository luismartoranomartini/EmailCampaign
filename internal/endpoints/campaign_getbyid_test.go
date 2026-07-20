package endpoints

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"projeto-golang/internal/contract"
	internalmock "projeto-golang/internal/test/internalMock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignsGetByID_ReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	campaign := contract.CampaignRespose{
		ID:      "343",
		Name:    "Test",
		Content: "Hi!",
		Status:  "Pending",
	}
	service := new(internalmock.CampaignServiceMock)
	service.On("GetBy", mock.Anything).Return(&campaign, nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	response, status, _ := handler.CampaignGetByID(rr, req)

	assert.Equal(200, status)
	assert.Equal(campaign.ID, response.(*contract.CampaignRespose).ID)
	assert.Equal(campaign.Name, response.(*contract.CampaignRespose).Name)
}

func Test_CampaignsGetByID_ReturnError(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	handler := Handler{CampaignService: service}
	errExpected := errors.New("something wrong")
	service.On("GetBy", mock.Anything).Return(nil, errExpected)
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	_, _, errReturned := handler.CampaignGetByID(rr, req)

	assert.Equal(errExpected.Error(), errReturned.Error())

}
