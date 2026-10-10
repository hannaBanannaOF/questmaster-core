package campaign

import (
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/pagination"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeListRepo struct {
	campaignApp.CampaignRepository
	result      pagination.Result[campaignDomain.Campaign]
	gotUser     userDomain.UserID
	gotFilters  campaignDomain.CampaignListFilters
	gotPage     pagination.Page
}

func (f *fakeListRepo) ListForUser(
	user userDomain.UserID,
	filters campaignDomain.CampaignListFilters,
	page pagination.Page,
) (pagination.Result[campaignDomain.Campaign], error) {
	f.gotUser, f.gotFilters, f.gotPage = user, filters, page
	return f.result, nil
}

// Listing each campaign once, ordering and filtering happen in the query (see the repository tests);
// the use case passes the requester, filters and page through.
func TestGetCurrentUserCampaignsPassesFiltersAndPage(t *testing.T) {
	user := userDomain.NewUserID(uuid.New())
	role := campaignDomain.RolePlayer
	filters := campaignDomain.CampaignListFilters{Role: &role}
	page := pagination.Page{Limit: 5, Offset: 10}
	repo := &fakeListRepo{result: pagination.Result[campaignDomain.Campaign]{
		Items: []campaignDomain.Campaign{{Id: 1}},
		Total: 11,
	}}

	got, err := NewGetCurrentUserMyCampaigns(repo).Execute(campaignApp.GetCurrentUserCampaignsCommand{
		UserID: user, Filters: filters, Page: page,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.gotUser != user || repo.gotFilters.Role != &role || repo.gotPage != page {
		t.Fatalf("expected requester, filters and page to reach the repository")
	}
	if got.Total != 11 || len(got.Items) != 1 {
		t.Fatalf("expected the repository result, got %+v", got)
	}
}
