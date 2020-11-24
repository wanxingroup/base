package service

type Interface interface {
	OnStart() error
	OnStop() error
	GetServiceName() string
}
