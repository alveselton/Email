package endpoints

import "email/internal/domain/campaign"

type Handler struct {
	CampaignService campaign.Service
}
