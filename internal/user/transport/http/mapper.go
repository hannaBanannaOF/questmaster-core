package user

import userDomain "questmaster-core/internal/user/domain"

func MapUserToUserResponse(user userDomain.User) UserResponse {
	var firstName *string
	if user.Name != nil {
		o := user.Name.FirstName()
		firstName = &o
	}
	return UserResponse{
		Id:       user.Id.Value().String(),
		Name:     firstName,
		Surname:  user.Name.LastName(),
		Username: user.Username.Value(),
	}
}
