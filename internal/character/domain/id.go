package character

type CharacterID int

func NewCharacterID(value int) CharacterID {
	return CharacterID(value)
}

func (c CharacterID) Value() int {
	return int(c)
}
