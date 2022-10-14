package server

type Conn interface {
	Serve() error
	Close() error
}

type Server interface {
	Serve() error
	Accept() (Conn, error)
	Close() error
}
