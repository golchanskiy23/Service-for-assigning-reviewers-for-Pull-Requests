package postgres

type Option func(source *DatabaseSource)

func SetMaxPoolSize(size int) Option {
	return func(s *DatabaseSource) {
		s.MaxPoolSize = size
	}
}
