package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

func MapRowToDomain(row CharacterRow) (characterDomain.Character, error) {
	characterName, err := characterDomain.NewCharacterName(row.Name)
	if err != nil {
		return characterDomain.Character{}, err
	}
	var campaignID *campaignDomain.CampaignID

	if row.CampaingID != nil {
		o := campaignDomain.NewCampaignID(*row.CampaingID)
		campaignID = &o
	}

	characterSlug, err := rpgDomain.NewSlug(row.Slug)
	if err != nil {
		return characterDomain.Character{}, err
	}

	var playerID *userDomain.UserID
	if row.PlayerID.Valid {
		id := userDomain.NewUserID(row.PlayerID.UUID)
		playerID = &id
	}

	system, err := rpgDomain.NewSystem(row.System)
	if err != nil {
		return characterDomain.Character{}, err
	}

	return characterDomain.Character{
		Id:         characterDomain.NewCharacterID(row.Id),
		Name:       characterName,
		PlayerID:   playerID,
		System:     system,
		CampaignID: campaignID,
		Slug:       characterSlug,
	}, nil
}
