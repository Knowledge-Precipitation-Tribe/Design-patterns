package auth

type ApiAuthenticator interface {
	Auth(url string) bool
}