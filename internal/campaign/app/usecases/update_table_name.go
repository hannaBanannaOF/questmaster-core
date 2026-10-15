package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type UpdateTableNameUseCase struct {
	r           campaignApp.CampaignRepository
	playerNames CampaignPlayerNameUpdater
}

func NewUpdateTableName(r campaignApp.CampaignRepository, playerNames CampaignPlayerNameUpdater) *UpdateTableNameUseCase {
	return &UpdateTableNameUseCase{r: r, playerNames: playerNames}
}

// Execute changes the requester's own table name: dm_name for the DM, otherwise their character's player_name
func (uc *UpdateTableNameUseCase) Execute(cmd campaignApp.UpdateTableNameCommand) (rpgDomain.TableName, error) {
	campaign, err := uc.r.FindById(cmd.CampaignID)
	if err != nil {
		return "", err
	}
	if campaign == nil {
		return "", ErrCampaignNotFound
	}

	if campaign.IsDM(cmd.UserID) {
		if err := uc.r.UpdateDmName(campaign.Id, cmd.Name); err != nil {
			return "", err
		}
		return cmd.Name, nil
	}

	updated, err := uc.playerNames.UpdatePlayerName(cmd.UserID, campaign.Id, cmd.Name)
	if err != nil {
		return "", err
	}
	if !updated {
		return "", campaignDomain.ErrNotCampaignMember
	}
	return cmd.Name, nil
}
