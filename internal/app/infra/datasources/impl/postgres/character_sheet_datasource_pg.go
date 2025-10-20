package datasource_pg

import (
	"context"
	"questmaster-core/internal/app/infra/db"
	"questmaster-core/internal/app/infra/models"
)

type CharacterSheetDataSourcePG struct {
	Db db.Db
}

func (ds *CharacterSheetDataSourcePG) GetAllByPlayerId(UserId string) ([]models.CharacterSheet, error) {
	ctx := context.Background()
	db, err := ds.Db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close(ctx)
	var data []models.CharacterSheet
	rows, err := db.Query(ctx, "SELECT c.id, c.character_name, c.player_id, c.session_id, c.trpg_system, c.slug FROM character_sheet c WHERE c.player_id = $1", UserId)
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
