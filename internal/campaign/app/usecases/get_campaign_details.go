package campaign

import (
	"log"

	campaignApp "questmaster-core/internal/campaign/app"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
)

type GetCampaignDetailsUseCase struct {
	getCampaignUc           GetCampaignFromIDUseCase
	getCampaignCharactersUc CampaignCharacterFinder
	getCampaignInviteUC     CampaignInviteFinder
	premadeFinder           CampaignPremadeCharactersFinder
	accessRecorder          CampaignAccessRecorder
}

func NewGetCampaignDetails(
	getCampaignUc GetCampaignFromIDUseCase,
	getCampaignCharactersUc CampaignCharacterFinder,
	getCampaignInviteUC CampaignInviteFinder,
	premadeFinder CampaignPremadeCharactersFinder,
	accessRecorder CampaignAccessRecorder,
) *GetCampaignDetailsUseCase {
	return &GetCampaignDetailsUseCase{
		getCampaignUc:           getCampaignUc,
		getCampaignCharactersUc: getCampaignCharactersUc,
		getCampaignInviteUC:     getCampaignInviteUC,
		premadeFinder:           premadeFinder,
		accessRecorder:          accessRecorder,
	}
}

func (uc *GetCampaignDetailsUseCase) Execute(cmd campaignApp.GetCampaignDetailsCommand) (campaignApp.CampaignDetailsReadModel, error) {
	campaign, err := uc.getCampaignUc.FindByID(cmd.ID)
	if err != nil {
		return campaignApp.CampaignDetailsReadModel{}, err
	}

	characters, err := uc.getCampaignCharactersUc.GetByCampaignID(cmd.ID)
	if err != nil {
		return campaignApp.CampaignDetailsReadModel{}, err
	}

	hasCharacter := false
	for _, c := range characters {
		if c.IsPlayer(cmd.UserID) {
			hasCharacter = true
			break
		}
	}

	if err := campaign.CanView(cmd.UserID, hasCharacter); err != nil {
		return campaignApp.CampaignDetailsReadModel{}, err
	}

	// Losing one recency update is harmless; failing the read is not
	if err := uc.accessRecorder.RecordAccess(cmd.UserID, cmd.ID); err != nil {
		log.Printf("recording access to campaign %d: %v", cmd.ID.Value(), err)
	}

	// The invite hash grants access to the campaign and premade characters are the DM preparation,
	// so only the DM gets them
	var invite *inviteDomain.Invite
	var premades []characterDomain.Character
	if campaign.IsDM(cmd.UserID) {
		invite, err = uc.getCampaignInviteUC.GetByCampaignID(cmd.ID)
		if err != nil {
			return campaignApp.CampaignDetailsReadModel{}, err
		}
		premades, err = uc.premadeFinder.GetPremadeByCampaignID(cmd.ID)
		if err != nil {
			return campaignApp.CampaignDetailsReadModel{}, err
		}
	}

	input := campaignApp.CampaignDetailsInput{
		Campaign:          campaign,
		Characters:        characters,
		PremadeCharacters: premades,
		Invite:            invite,
		UserID:            cmd.UserID,
	}

	return campaignApp.MapDomainToDetailReadModel(input), nil
}
