package simplehttp

func MiddlewareChain(middlewares ...RouteMiddleware) []RouteMiddleware {
	return middlewares
}
