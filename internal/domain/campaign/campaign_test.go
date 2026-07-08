package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var (
	name      = "Campaign X"
	content   = "Body Hi!"
	contacts  = []string{"email1@e.com", "email2@e.com"}
	fake      = faker.New()
	createdby = "teste@teste.com.br"
)

// Testar o ponto de integração
// Arrange, Action, Assert
func Test_NewCampaign_CreateCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdby)

	assert.NotEmpty(campaign.ID)
	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
	assert.Equal(createdby, campaign.Createdby)
}

func Test_NewCampaign_IDIsNotNil(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdby)

	assert.NotNil(campaign.ID)
}

func Test_NewCampaign_StatusPending(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdby)

	assert.Equal(Pending, campaign.Status)
}

func Test_NewCampaign_StatusStarted(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdby)

	campaign.Create()

	assert.Equal(Started, campaign.Status)
}

func Test_NewCampaign_StatusFinished(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdby)

	campaign.Start()
	err := campaign.Finished()

	assert.Nil(err)
	assert.Equal(Done, campaign.Status)
}

func Test_NewCampaign_CreatedOnMustBeNow(t *testing.T) {
	assert := assert.New(t)

	now := time.Now().Add(-time.Minute)

	campaign, _ := NewCampaign(name, content, contacts, createdby)

	assert.Greater(campaign.CreatedOn, now)
}

func Test_NewCampaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign("", content, contacts, createdby)

	assert.Equal("name is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(30), content, contacts, createdby)

	assert.Equal("name is required with max 24", err.Error())
}

func Test_NewCampaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(1), contacts, createdby)

	assert.Equal("content is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(1040), contacts, createdby)

	assert.Equal("content is required with max 1024", err.Error())
}

func Test_NewCampaign_MustValidateContactsMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, nil, createdby)

	assert.Equal("contacts is required with min 1", err.Error())
}

func Test_NewCampaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{}, createdby)

	assert.Equal("contacts is required", err.Error())
}

func Test_NewCampaign_MustValidateCreatedBy(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, contacts, "")

	assert.Equal("createdby is invalid", err.Error())
}
