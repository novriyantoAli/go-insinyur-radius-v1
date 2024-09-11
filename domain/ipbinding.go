package domain

import "github.com/go-routeros/routeros/v3"

type IPBinding struct {
	ID         string `json:"id"`
	MacAddress string `json:"mac-address"`
	Comment    string `json:"comment"`
	Type       string `json:"type"`
}

type IPBindings struct {
	RName    string
	Bindings []IPBinding
}

type IPBind struct {
	Rname string
	Bind  IPBinding
}

type IPBindingRepository interface {
	Find(params *IPBinding) (res []IPBinding, err error)
	Finds(client *routeros.Client, params *IPBinding) (res []IPBinding, err error)
	First(param *IPBinding) (res IPBinding, err error)
	Firsts(client *routeros.Client, params *IPBinding) (res IPBinding, err error)
	Save(params *IPBinding) (err error)
	Saves(client *routeros.Client, params *IPBinding) (err error)
	Update(params *IPBinding) (err error)
	Updates(client *routeros.Client, params *IPBinding) (err error)
	Delete(params *IPBinding) (err error)
	Deletes(client *routeros.Client, params *IPBinding) (err error)
}
