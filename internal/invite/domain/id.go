package invite

type InviteID int

func NewInviteID(value int) InviteID {
	return InviteID(value)
}

func (c InviteID) Value() int {
	return int(c)
}
