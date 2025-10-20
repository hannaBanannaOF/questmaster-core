package enum

import "fmt"

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
			return fmt.Errorf("unknown status string: %s", string(v))
		}
	default:
		return fmt.Errorf("unsupported type for Status: %T", value)
	}
	return nil
}
