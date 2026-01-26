package campaign

import (
	"errors"
	"strings"
)

type CampaignName string

func NewCampaignName(value string) (CampaignName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("Empty campaign name")
	}
	return CampaignName(value), nil
}

func (n CampaignName) String() string {
	return string(n)
}
