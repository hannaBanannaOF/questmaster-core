package campaign

type CampaignOverview string

func NewCampaignOverview(value string) (CampaignOverview, error) {
	return CampaignOverview(value), nil
}

func (o CampaignOverview) String() string {
	return string(o)
}
