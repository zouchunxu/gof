package server

type Server interface {
	Endpoint() (string, error)
	Start() error
	Stop() error
}
