package client

import "github.com/oxzjh/stream"

type options struct {
	maker  stream.IMaker
	parser stream.IParser
}

type Option func(*options)

func WithMaker(maker stream.IMaker) Option {
	return func(o *options) {
		o.maker = maker
	}
}

func WithParser(parser stream.IParser) Option {
	return func(o *options) {
		o.parser = parser
	}
}
