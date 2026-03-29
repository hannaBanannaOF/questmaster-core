package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
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
	var characterHp *characterDomain.HP
	if row.CurrentHp != nil || row.MaxHp != nil {
		hp, err := characterDomain.NewHP(*row.CurrentHp, *row.MaxHp)
		if err != nil {
			return characterDomain.Character{}, err
		}
		characterHp = &hp
	}

	system, err := rpgDomain.NewSystem(row.System)
	if err != nil {
		return characterDomain.Character{}, err
	}

	return characterDomain.Character{
		Id:         characterDomain.NewCharacterID(row.Id),
		Name:       characterName,
		PlayerID:   rpgDomain.NewUserID(row.PlayerID),
		System:     system,
		CampaignID: campaignID,
		Slug:       characterSlug,
		Hp:         characterHp,
	}, nil
}
