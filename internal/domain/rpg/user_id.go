package rpg

import "github.com/google/uuid"

type UserID struct {
	value uuid.UUID
}

func NewUserID(id uuid.UUID) UserID {
	return UserID{value: id}
}

func (id UserID) UUID() uuid.UUID {
	return id.value
}
