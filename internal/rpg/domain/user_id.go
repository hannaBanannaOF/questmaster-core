package rpg

import "github.com/google/uuid"

type UserID uuid.UUID

func NewUserID(id uuid.UUID) UserID {
	return UserID(id)
}

func (id UserID) Value() uuid.UUID {
	return uuid.UUID(id)
}
