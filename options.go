package golactus

type Option func(*Server)

func Name(name string) Option {
	return func(s *Server) {
		s.name = name
	}
}

func Host(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func Port(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func RequestTimeout(requestTimeout int) Option {
	return func(s *Server) {
		s.requestTimeout = requestTimeout
	}
}
