package endpoints

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignsStart_200(t *testing.T) {
	setup()
	campaignId := "xpto"
	service.On("Start", mock.MatchedBy(func(id string) bool {
		return id == campaignId
	})).Return(nil)
	req, rr := newHttpTest("PATCH", "/", nil)
	req = addParameter(req, "id", campaignId)

	_, status, err := handler.CampaignStart(rr, req)

	assert.Equal(t, 200, status)
	assert.Nil(t, err)
}

func Test_CampaignStart_Err(t *testing.T) {
	setup()

	errExpected := errors.New("something wrong")
	service.On("Start", mock.Anything).Return(errExpected)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("PATCH", "/", nil)
	rr := httptest.NewRecorder()

	_, _, err := handler.CampaignStart(rr, req)

	assert.Equal(t, errExpected, err)
}
