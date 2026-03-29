package campaign

type CampaignOverview string

func NewCampaignOverview(value string) CampaignOverview {
	return CampaignOverview(value)
}

func (o CampaignOverview) Value() string {
	return string(o)
}
