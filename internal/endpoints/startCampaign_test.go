package endpoints

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	internalmock "projeto-golang/internal/test/internalMock"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignsStart_Return200(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	campaignID := "xpto"
	service.On("Start", mock.MatchedBy(func(id string) bool {
		return id == campaignID
	})).Return(nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("PATCH", "/", nil)

	// importante
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("id", campaignID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiContext))
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
