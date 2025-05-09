package server

type AppBehaviour interface {
	routeSetup()
	middlewareSetup()
	Start() error
	Stop() error
}
