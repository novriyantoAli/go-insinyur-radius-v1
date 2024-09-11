package domain

import "github.com/go-routeros/routeros/v3"

type SimpleQueue struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Target      string `json:"target"`
	PacketMarks string `json:"packet_marks"`
	MaxLimit    string `json:"max_limit"`
}

type SimpleQueueRepository interface {
	Find(param *SimpleQueue) (res []SimpleQueue, err error)
	Finds(client *routeros.Client, param *SimpleQueue) (res []SimpleQueue, err error)
	First(param *SimpleQueue) (res SimpleQueue, err error)
	Firsts(client *routeros.Client, param *SimpleQueue) (res SimpleQueue, err error)
	Save(param *SimpleQueue) (err error)
	Saves(client *routeros.Client, param *SimpleQueue) (err error)
	Update(param *SimpleQueue) (err error)
	Updates(client *routeros.Client, param *SimpleQueue) (err error)
	Delete(param *SimpleQueue) (err error)
	Deletes(client *routeros.Client, param *SimpleQueue) (err error)
}
