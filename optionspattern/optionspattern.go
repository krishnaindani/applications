package optionspattern

import "time"

type Server struct {
	Host    string
	Port    int
	Timeout time.Duration
}

type Option func(*Server)

func defaultServerConfig() *Server {
	return &Server{
		Host:    "localhost",
		Port:    80,
		Timeout: time.Second,
	}
}

func NewServer(opts ...Option) *Server {
	s := defaultServerConfig()

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithHost(host string) Option {
	return func(s *Server) {
		s.Host = host
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

func main() {
	s := NewServer(
		WithHost("localhost"),
		WithPort(80),
		WithTimeout(time.Minute),
	)

	_ = s
}
