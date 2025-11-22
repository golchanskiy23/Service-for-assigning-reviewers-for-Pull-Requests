package database

import "time"

type Option func(source *DatabaseSource)

func SetMaxPoolSize(size int) Option {
	return func(s *DatabaseSource) {
		s.MaxPoolSize = size
	}
}

func SetMaxConnLifetime(duration time.Duration) Option {
	return func(s *DatabaseSource) {
		s.MaxConnLifetime = duration
	}
}

func SetMaxConnectTimeout(duration time.Duration) Option {
	return func(s *DatabaseSource) {
		s.MaxConnectTimeout = duration
	}
}
