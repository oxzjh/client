package client

import "io"

type options struct {
	maker  func(int) []byte
	parser func(io.Reader) ([]byte, error)
}

type Option func(*options)

func WithMaker(maker func(int) []byte) Option {
	return func(o *options) {
		o.maker = maker
	}
}

func WithParser(parser func(io.Reader) ([]byte, error)) Option {
	return func(o *options) {
		o.parser = parser
	}
}
