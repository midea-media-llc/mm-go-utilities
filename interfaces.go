package utils

type IClaims interface {
	GetId() int64
	GetClientId() int64
	GetUnitId() int64
	GetUsername() string
	GetEmail() string
	GetFullname() string
	GetPhone() string
	GetIsAdmin() bool
	GetIsSystem() bool
	GetLanguage() string
	GetIsBaseLanguage() bool
}
