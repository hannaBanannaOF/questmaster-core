package campaign

type PlayerCount int

func NewPlayerCount(value int) PlayerCount {
	return PlayerCount(value)
}

func (p PlayerCount) Value() int {
	return int(p)
}
