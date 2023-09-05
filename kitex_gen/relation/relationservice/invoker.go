// Code generated by Kitex v0.7.0. DO NOT EDIT.

package relationservice

import (
	server "github.com/cloudwego/kitex/server"
	relation "tiktok_v2/kitex_gen/relation"
)

// NewInvoker creates a server.Invoker with the given service and options.
func NewInvoker(handler relation.RelationService, opts ...server.Option) server.Invoker {
	var options []server.Option

	options = append(options, opts...)

	s := server.NewInvoker(options...)
	if err := s.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	return s
}
