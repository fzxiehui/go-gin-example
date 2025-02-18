package auth_service

import "github.com/EDDYCJY/go-gin-example/models"

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (bool, error) {
	// Check username and password
	return models.CheckAuth(a.Username, a.Password)
}
