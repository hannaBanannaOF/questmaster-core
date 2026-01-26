package bootstrap

import (
	usecases "questmaster-core/internal/app/character/usecases"
	infra "questmaster-core/internal/infra/character/pg"
	transport "questmaster-core/internal/transport/character/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildCharacterHandler(db *pgxpool.Pool) *transport.CharactersHandler {
	repo := infra.NewCharacterRepositoryPG(db)
	fetchUc := usecases.NewFetchMyCharacter(repo)
	createUc := usecases.NewCreateCharacter(repo)
	resolveSlugUc := usecases.NewResolveCharacterSlug(repo)
	getDetailUc := usecases.NewGetCharacterDetail(repo)
	return transport.NewCharactersHandler(fetchUc, createUc, resolveSlugUc, getDetailUc)
}
