package rpg

type System string

const (
	DungeonsAndDragons System = "DUNGEONS_AND_DRAGONS"
	CallOfCthulhu      System = "CALL_OF_CTHULHU"
	CyberpunkRed       System = "CYBERPUNK_RED"
	OrdemParanormal    System = "ORDEM_PARANORMAL"
)

func NewSystem(value string) (System, error) {
	system := System(value)
	switch system {
	case CallOfCthulhu,
		DungeonsAndDragons,
		CyberpunkRed,
		OrdemParanormal:
		return system, nil

	default:
		return "", ErrInvalidSystem
	}
}

func (s System) Value() string {
	return string(s)
}
