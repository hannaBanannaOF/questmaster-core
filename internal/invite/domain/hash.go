package invite

import "github.com/google/uuid"

type InviteHash uuid.UUID

func NewHash(hash uuid.UUID) InviteHash {
	return InviteHash(hash)
}

func (id InviteHash) Value() uuid.UUID {
	return uuid.UUID(id)
}
