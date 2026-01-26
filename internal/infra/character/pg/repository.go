package character

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	app "questmaster-core/internal/app/character"
	domain "questmaster-core/internal/domain/character"
	"questmaster-core/internal/domain/rpg"
)

type CharacterRepositoryPG struct {
	db *pgxpool.Pool
}

func NewCharacterRepositoryPG(db *pgxpool.Pool) *CharacterRepositoryPG {
	return &CharacterRepositoryPG{db: db}
}

func (r *CharacterRepositoryPG) GetAllByPlayerId(
	userId rpg.UserID,
) ([]domain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT cs.*
        FROM character_sheet cs
        WHERE cs.player_id = $1
    `, userId.UUID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return nil, err
	}

	domain := make([]domain.Character, 0)

	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		domain = append(domain, val)
	}

	return domain, nil
}

func (r *CharacterRepositoryPG) FindBySlug(slug rpg.Slug) (*domain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT cs.*
		FROM character_sheet cs
		WHERE cs.slug = $1
	`, slug.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *CharacterRepositoryPG) FindById(characterId domain.CharacterID) (*domain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT cs.*
		FROM character_sheet cs
		WHERE cs.id = $1
	`, int(characterId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *CharacterRepositoryPG) Create(input app.CreateCharacterInput) (domain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		INSERT INTO character_sheet(name, player_id, trpg_system, max_hp, current_hp) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING *
	`, input.Name.String(), input.Player.UUID(), input.System, input.Hp.Max(), input.Hp.Current())
	if err != nil {
		return domain.Character{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return domain.Character{}, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return domain.Character{}, err
	}

	return val, nil
}
