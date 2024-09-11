package domain

import "github.com/go-routeros/routeros/v3"

type ROSClient struct {
	Name   string
	Client *routeros.Client
}
