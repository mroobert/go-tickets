package signup

import "github.com/mroobert/go-tickets/auth/internal/usecase/signup/vstruct"

// signUpRequestDto represents the payload request contract.
type signUpRequestDto struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

// dtoToUser transforms signup payload (dto) into user domain struct.
func dtoToSignUpUser(dto signUpRequestDto) (vstruct.SignUpUser, error) {
	su, err := vstruct.NewSignUpUser(dto.Email, dto.Password, dto.DisplayName)
	if err != nil {
		return vstruct.SignUpUser{}, err
	}
	return su, nil
}

// signUpResponseDto represents the payload response contract.
type signUpResponseDto struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}
