package web

import (
	"errors"
	"strings"
)

func ExtractToken(header string) (string, error) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "jwt" {
		return "", errors.New("expected authorization header format: jwt <token>")
	}
	return parts[1], nil
}
