package campaign

import (
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	"questmaster-core/internal/shared/pagination"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeListRepo struct {
	campaignApp.CampaignRepository
	result     pagination.Result[campaignDomain.Campaign]
	gotUser    userDomain.UserID
	gotFilters campaignDomain.CampaignListFilters
	gotPage    pagination.Page
}

func (f *fakeListRepo) ListForUser(
	user userDomain.UserID,
	filters campaignDomain.CampaignListFilters,
	page pagination.Page,
) (pagination.Result[campaignDomain.Campaign], error) {
	f.gotUser, f.gotFilters, f.gotPage = user, filters, page
	return f.result, nil
}

type fakeMyCharactersFinder struct {
	characters []characterDomain.Character
	gotIDs     []campaignDomain.CampaignID
	called     bool
}

func (f *fakeMyCharactersFinder) GetByPlayerInCampaigns(_ userDomain.UserID, ids []campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	f.called, f.gotIDs = true, ids
	return f.characters, nil
}

func campaignRef(id int) *campaignDomain.CampaignID {
	c := campaignDomain.NewCampaignID(id)
	return &c
}

// Listing each campaign once, ordering and filtering happen in the query (see the repository tests)
func TestGetCurrentUserCampaigns(t *testing.T) {
	user := userDomain.NewUserID(uuid.New())

	t.Run("passes filters and page, and attaches the requester characters per campaign", func(t *testing.T) {
		role := campaignDomain.RolePlayer
		filters := campaignDomain.CampaignListFilters{Role: &role}
		page := pagination.Page{Limit: 5, Offset: 10}
		repo := &fakeListRepo{result: pagination.Result[campaignDomain.Campaign]{
			Items: []campaignDomain.Campaign{{Id: 1}, {Id: 2}},
			Total: 12,
		}}
		finder := &fakeMyCharactersFinder{characters: []characterDomain.Character{
			{Name: "Alva", Slug: "alva", CampaignID: campaignRef(1)},
			{Name: "Bram", Slug: "bram", CampaignID: campaignRef(1)},
		}}

		got, err := NewGetCurrentUserMyCampaigns(repo, finder).Execute(campaignApp.GetCurrentUserCampaignsCommand{
			UserID: user, Filters: filters, Page: page,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.gotUser != user || repo.gotFilters.Role != &role || repo.gotPage != page {
			t.Fatalf("expected requester, filters and page to reach the repository")
		}
		if len(finder.gotIDs) != 2 || finder.gotIDs[0] != 1 || finder.gotIDs[1] != 2 {
			t.Fatalf("expected characters to be loaded for the page campaigns, got %v", finder.gotIDs)
		}
		if got.Total != 12 || len(got.Items) != 2 {
			t.Fatalf("unexpected page %+v", got)
		}
		if mine := got.Items[0].MyCharacters; len(mine) != 2 || mine[0].Name != "Alva" || mine[1].Name != "Bram" {
			t.Fatalf("expected both characters in campaign 1, got %v", mine)
		}
		if len(got.Items[1].MyCharacters) != 0 {
			t.Fatalf("expected no characters in campaign 2, got %v", got.Items[1].MyCharacters)
		}
	})

	t.Run("empty page doesn't load characters", func(t *testing.T) {
		finder := &fakeMyCharactersFinder{}
		got, err := NewGetCurrentUserMyCampaigns(&fakeListRepo{}, finder).Execute(campaignApp.GetCurrentUserCampaignsCommand{UserID: user})
		if err != nil || len(got.Items) != 0 || finder.called {
			t.Fatalf("expected an empty page without loading characters, got %+v called=%v err=%v", got, finder.called, err)
		}
	})
}
