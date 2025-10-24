package serviceinterface

import (
	"challenge-app/internal/domain/model"
)

type AuthServicer interface {
    // Auth Operations
    RegisterUser(username, email, password, bio string) (*model.UserModel, string, error)
    LoginUser(email, password string) (*model.UserModel, string, error)
}
