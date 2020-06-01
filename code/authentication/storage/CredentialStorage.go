package storage

type CredentialStorage interface {
	GetPasswordByAppID(appID string)
}