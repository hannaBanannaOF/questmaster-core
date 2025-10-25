package datasource_pg

import (
	"context"
	enum "questmaster-core/domain/enumerations"
	"questmaster-core/internal/app/infra/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CharacterSheetDataSourcePG struct {
	Db *pgxpool.Pool
}

func (ds *CharacterSheetDataSourcePG) GetAllByPlayerId(UserId string) ([]models.CharacterSheet, error) {
	ctx := context.Background()
	var data []models.CharacterSheet
	rows, err := ds.Db.Query(ctx, "SELECT c.id, c.character_name, c.player_id, c.session_id, c.trpg_system, c.slug FROM character_sheet c WHERE c.player_id = $1", UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s models.CharacterSheet
		err := rows.Scan(&s.Id, &s.CharacterName, &s.PlayerId, &s.SessionId, &s.TrpgSystem, &s.Slug)
		if err != nil {
			return nil, err
		}
		data = append(data, s)
	}
	return data, nil
}

func (ds *CharacterSheetDataSourcePG) GetAllBySessionId(SessionId int) ([]models.CharacterSheet, error) {
	ctx := context.Background()
	var data []models.CharacterSheet
	rows, err := ds.Db.Query(ctx, "SELECT c.id, c.character_name, c.player_id, c.session_id, c.trpg_system, c.slug, c.max_hp, c.current_hp FROM character_sheet c WHERE c.session_id = $1", SessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s models.CharacterSheet
		err := rows.Scan(&s.Id, &s.CharacterName, &s.PlayerId, &s.SessionId, &s.TrpgSystem, &s.Slug, &s.MaxHp, &s.CurrentHp)
		if err != nil {
			return nil, err
		}
		data = append(data, s)
	}
	return data, nil
}

func (ds *CharacterSheetDataSourcePG) GetOne(CharacterSheetId int) (*models.CharacterSheet, error) {
	ctx := context.Background()
	row := ds.Db.QueryRow(ctx, "SELECT cs.id, cs.character_name, cs.player_id, cs.session_id, cs.trpg_system, cs.slug, cs.max_hp, cs.current_hp FROM character_sheet cs WHERE cs.id = $1", CharacterSheetId)
	var s models.CharacterSheet
	err := row.Scan(&s.Id, &s.CharacterName, &s.PlayerId, &s.SessionId, &s.TrpgSystem, &s.Slug, &s.MaxHp, &s.CurrentHp)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &s, nil
}

func (ds *CharacterSheetDataSourcePG) ResolveSlug(Slug string) (*int, error) {
	ctx := context.Background()
	row := ds.Db.QueryRow(ctx, "SELECT cs.id FROM character_sheet cs WHERE cs.slug = $1", Slug)
	var s int
	err := row.Scan(&s)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &s, nil
}

func (ds *CharacterSheetDataSourcePG) CreateCharacterSheet(CharacterName string, MaxHp int, System enum.TrpgSystem, UserId string) (*models.CharacterSheet, error) {
	ctx := context.Background()
	row := ds.Db.QueryRow(ctx, "INSERT INTO character_sheet(character_name, player_id, trpg_system, max_hp, current_hp) VALUES($1, $2, $3, $4, $4) RETURNING *", CharacterName, UserId, System, MaxHp)
	var s models.CharacterSheet
	err := row.Scan(&s.Id, &s.CharacterName, &s.PlayerId, &s.SessionId, &s.TrpgSystem, &s.Slug, &s.MaxHp, &s.CurrentHp)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &s, nil
}
