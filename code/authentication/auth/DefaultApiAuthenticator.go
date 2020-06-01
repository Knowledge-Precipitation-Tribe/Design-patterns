package auth

import (
	"code/authentication/apiRequest"
	"code/authentication/storage"
	"fmt"
)

type DefaultApiAuthenticator struct {
	storage storage.CredentialStorage
}

func (DefaultApiAuthenticator) Auth(url string) bool{
	req := apiRequest.GenFakeReq(url)
	fmt.Print(req)
	return true
}