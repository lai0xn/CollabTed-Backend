package storage

type Store interface {
	Init()
	Start()
	Stop()
}
