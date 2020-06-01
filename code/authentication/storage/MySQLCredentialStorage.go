package storage

type MySQLCredentialStorage struct {

}

func (MySQLCredentialStorage) GetPasswordByAppID(appID string) string{
	return ""
}