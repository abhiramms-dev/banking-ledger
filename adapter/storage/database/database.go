package database

type DBconnection interface {
	Connection() (interface{}, error)
}
