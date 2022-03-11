package signup

type newUserDto struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

func dtoToNewUser(dto newUserDto) newUser {
	nu := newUser(dto)
	return nu
}

type userDto struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

func userToDto(u user) userDto {
	dto := userDto{
		Email:       u.Email,
		DisplayName: u.DisplayName,
	}
	return dto
}
