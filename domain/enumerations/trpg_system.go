package enum

type TrpgSystem string

const (
	Unknown       = ""
	CallOfCthulhu = "CALL_OF_CTHULHU"
)

func (s *TrpgSystem) Scan(value interface{}) error {
	if value == nil {
		*s = Unknown
		return nil
	}

	switch v := value.(type) {
	case string:
		// Handle string representation if stored as text in DB
		switch v {
		case "CALL_OF_CTHULHU":
			*s = CallOfCthulhu
		default:
			*s = Unknown
		}
	default:
		*s = Unknown
	}
	return nil
}
