package database

type TokenRepository interface {
	SaveAccessToken(chatId int64, token string) error
	GetAccessToken(chatId int64) (string, error)
	SaveRequestToken(chatId int64, token string) error
	GetRequestToken(chatId int64) (string, error)
}
