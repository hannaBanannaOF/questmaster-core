package character

import "errors"

type HP struct {
	current int
	max     int
}

func NewHP(current, max int) (HP, error) {
	if max <= 0 {
		return HP{}, errors.New("max hp must be > 0")
	}
	if current < 0 || current > max {
		return HP{}, errors.New("current hp out of bounds")
	}
	return HP{current: current, max: max}, nil
}

func (h HP) Current() int { return h.current }
func (h HP) Max() int     { return h.max }
func (h HP) IsAlive() bool {
	return h.current > 0
}
