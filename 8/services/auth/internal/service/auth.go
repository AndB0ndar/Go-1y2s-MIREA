package service

type AuthError struct {
	msg string
}

func (e *AuthError) Error() string {
	return e.msg
}

var ErrInvalidCredentials = &AuthError{msg: "invalid credentials"}

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(username, password string) (string, error) {
	if username != "student" || password != "student" { // FIXME
		return "", ErrInvalidCredentials
	}
	return "demo-token", nil
}

func (s *AuthService) Verify(token string) (bool, string) {
	if token == "demo-token" {
		return true, "student"
	}
	return false, ""
}
