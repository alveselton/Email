package campaign

import (
	statusreturnmessage "email/internal/StatusReturnMessage"
	"email/internal/contract"
	"errors"
)

type Service interface {
	Create(newCampaign contract.NewCampaign) (string, error)
	GetBy(id string) (*contract.CampaignResponse, error)
	Delete(id string) error
	Start(id string) error
}

type ServiceImp struct {
	Repository Repository
	SendMail   func(campaign *Campaign) error
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaign) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)

	if err != nil {
		return "", err
	}

	err = s.Repository.Create(campaign)

	if err != nil {
		return "", statusreturnmessage.ErrInternal
	}

	return campaign.ID, nil
}

func (s *ServiceImp) GetBy(id string) (*contract.CampaignResponse, error) {

	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return nil, statusreturnmessage.ProcessErrorToReturn(err)
	}

	return &contract.CampaignResponse{
		ID:                   campaign.ID,
		Name:                 campaign.Name,
		Content:              campaign.Content,
		Status:               campaign.Status,
		AmountOfEmailsToSend: len(campaign.Contacts),
		CreatedBy:            campaign.CreatedBy,
	}, nil
}

func (s *ServiceImp) Delete(id string) error {

	campaignSaved, err := s.getAndValidateStatusIsPending(id)

	if err != nil {
		return err
	}

	campaignSaved.Delete()
	err = s.Repository.Delete(campaignSaved)
	if err != nil {
		return statusreturnmessage.ErrInternal
	}

	return nil
}

func (s *ServiceImp) Start(id string) error {

	campaignSaved, err := s.getAndValidateStatusIsPending(id)

	if err != nil {
		return err
	}

	go func() {
		err = s.SendMail(campaignSaved)
		if err != nil {
			campaignSaved.Fail()
		} else {
			campaignSaved.Done()
		}
		s.Repository.Update(campaignSaved)
	}()

	campaignSaved.Started()
	err = s.Repository.Update(campaignSaved)
	if err != nil {
		return statusreturnmessage.ErrInternal
	}

	return nil
}

func (s *ServiceImp) getAndValidateStatusIsPending(id string) (*Campaign, error) {
	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return nil, statusreturnmessage.ProcessErrorToReturn(err)
	}

	if campaign.Status != Pending {
		return nil, errors.New("Campaign status invalid")
	}
	return campaign, nil
}
