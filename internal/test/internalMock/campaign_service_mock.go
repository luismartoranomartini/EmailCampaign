package mock

import (
	"github.com/stretchr/testify/mock"
	"projeto-golang/internal/contract"
)

type CampaignServiceMock struct {
	mock.Mock
}

func (r *CampaignServiceMock) Create(newCampaign contract.NewCampaign) (string, error) {
	args := r.Called(newCampaign)
	return args.String(0), args.Error(1)
}

func (r *CampaignServiceMock) GetBy(id string) (*contract.CampaignRespose, error) {
	args := r.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contract.CampaignRespose), args.Error(1)
}

func (r *CampaignServiceMock) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *CampaignServiceMock) Start(id string) error {
	args := r.Called(id)
	return args.Error(0)
}
