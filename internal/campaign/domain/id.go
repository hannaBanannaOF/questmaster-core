package campaign

type CampaignID int

func NewCampaignID(value int) CampaignID {
	return CampaignID(value)
}

func (c CampaignID) Value() int {
	return int(c)
}
