package model

import "github.com/google/uuid"

type (
	RegisterResponse struct {
		Message string `json:"message"`
	}
	LoginResponse struct {
		Message                        string                         `json:"message"`
		Token                          string                         `json:"token"`
		EmbeddedUserDataInLoginSuccess EmbeddedUserDataInLoginSuccess `json:"user_metadata"`
	}
	EmbeddedUserDataInLoginSuccess struct {
		UserId    uuid.UUID `json:"user_id"`
		Email     string    `json:"email"`
		FirstName string    `json:"first_name"`
		RoleName  string    `json:"role_name"`
	}
)
