package utils

import "database/sql"

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

type IGormDB interface {
	Raw(sql string, values ...interface{}) IGormDB
	Scan(destination interface{}) IGormDB
	ScanRows(rows *sql.Rows, dest interface{}) error
	Rows() (*sql.Rows, error)
}
