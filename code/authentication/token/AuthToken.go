package token

import "time"

const EXPIRE_INTERVAL = 10 * 60 * 1000


type AuthToken struct {
	createTime int64
	token string
	originalUrl string
}

func (a AuthToken) IsExpired() bool {
	if a.createTime > time.Now().Unix(){
		return true
	}
	return false
}

func (a AuthToken) Match() bool{
	return false
}