package campaign

import "errors"

var ErrInvalidCampaignRole = errors.New("Invalid campaign role")

// CampaignRole is the requester's role in a campaign, used to filter their campaign list
type CampaignRole string

const (
	RoleDM     CampaignRole = "dm"
	RolePlayer CampaignRole = "player"
)

func NewCampaignRole(value string) (CampaignRole, error) {
	role := CampaignRole(value)
	switch role {
	case RoleDM, RolePlayer:
		return role, nil
	default:
		return "", ErrInvalidCampaignRole
	}
}

func (r CampaignRole) Value() string {
	return string(r)
}

// CampaignListFilters narrows the requester's campaigns. Nil fields don't filter.
type CampaignListFilters struct {
	Role   *CampaignRole
	Status *CampaignStatus
}
