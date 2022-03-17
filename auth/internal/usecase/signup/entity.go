package signup

// user represents a domain entity.
type user struct {
	UID         string
	Email       string
	DisplayName string
}

// userToSignUpResponseDto transforms user domain struct into signup response (dto).
func userToSignUpResponseDto(u user) signUpResponseDto {
	dto := signUpResponseDto{
		Email:       u.Email,
		DisplayName: u.DisplayName,
	}
	return dto
}
