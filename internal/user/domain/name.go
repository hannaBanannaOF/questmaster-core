package user

import "strings"

type Name struct {
	firstName string
	lastName  *string
}

func NewName(firstName string, lastName *string) (Name, error) {
	firstName = strings.TrimSpace(firstName)
	if firstName == "" {
		return Name{}, ErrInvalidFirstname
	}
	return Name{firstName: firstName, lastName: lastName}, nil
}

func (n Name) FirstName() string {
	return n.firstName
}

func (n Name) LastName() *string {
	return n.lastName
}
