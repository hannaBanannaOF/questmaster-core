package user

import userDomain "questmaster-core/internal/user/domain"

func MapUserToUserResponse(user userDomain.User) UserResponse {
	var firstName *string
	var lastName *string
	if user.Name != nil {
		o := user.Name.FirstName()
		firstName = &o
		lastName = user.Name.LastName()
	}
	return UserResponse{
		Id:       user.Id.Value().String(),
		Name:     firstName,
		Surname:  lastName,
		Username: user.Username.Value(),
	}
}
