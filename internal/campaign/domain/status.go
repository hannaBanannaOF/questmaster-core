package campaign

type CampaignStatus string

const (
	StatusDraft    CampaignStatus = "DRAFT"
	StatusActive   CampaignStatus = "ACTIVE"
	StatusPaused   CampaignStatus = "PAUSED"
	StatusArchived CampaignStatus = "ARCHIVED"
)

func NewCampaignStatus(value string) (CampaignStatus, error) {
	domainStatus := CampaignStatus(value)
	switch domainStatus {
	case StatusDraft,
		StatusActive,
		StatusPaused,
		StatusArchived:
		return domainStatus, nil
	default:
		return "", ErrInvalidCampaignStatus
	}
}

func (c CampaignStatus) Value() string {
	return string(c)
}

func (c CampaignStatus) CanTransition(to CampaignStatus) bool {
	switch c {
	case StatusDraft:
		return to == StatusActive

	case StatusActive:
		return to == StatusPaused || to == StatusArchived

	case StatusPaused:
		return to == StatusActive || to == StatusArchived

	case StatusArchived:
		return false
	}

	return false
}
