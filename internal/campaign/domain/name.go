package campaign

import (
	"strings"
)

type CampaignName string

func NewCampaignName(value string) (CampaignName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrEmptyCampaignName
	}
	return CampaignName(value), nil
}

func (n CampaignName) Value() string {
	return string(n)
}
