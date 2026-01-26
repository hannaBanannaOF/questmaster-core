package character

import (
	campaignDomain "questmaster-core/internal/domain/campaign"
	domain "questmaster-core/internal/domain/character"
	"questmaster-core/internal/domain/rpg"
)

func MapRowToDomain(row CharacterRow) (domain.Character, error) {
	characterName, err := domain.NewCharacterName(row.Name)
	if err != nil {
		return domain.Character{}, err
	}
	var campaignId *campaignDomain.CampaignID

	if row.CampaingId != nil {
		o := campaignDomain.CampaignID(*row.CampaingId)
		campaignId = &o
	}

	characterSlug, err := rpg.NewSlug(row.Slug)
	if err != nil {
		return domain.Character{}, err
	}
	var characterHp *domain.HP
	if row.CurrentHp != nil || row.MaxHp != nil {
		hp, err := domain.NewHP(*row.CurrentHp, *row.MaxHp)
		if err != nil {
			return domain.Character{}, err
		}
		characterHp = &hp
	}
	return domain.Character{
		Id:         domain.CharacterID(row.Id),
		Name:       characterName,
		PlayerId:   rpg.NewUserID(row.PlayerId),
		System:     rpg.System(row.System),
		CampaingId: campaignId,
		Slug:       characterSlug,
		Hp:         characterHp,
	}, nil
}
