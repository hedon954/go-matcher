package znet

import "github.com/hedon954/go-matcher/pkg/zinx/ziface"

type BaseRouter struct{}

func (r *BaseRouter) PreHandle(request ziface.IRequest)  {}
func (r *BaseRouter) PostHandle(request ziface.IRequest) {}
