package domain

import "context"

type Radreply struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Attribute string `json:"attribute"`
	OP        string `json:"op"`
	Value     string `json:"value"`
}

func (Radreply) TableName() string {
	return "radreply"
}

type RadreplyRepository interface {
	First(ctx context.Context, radreply *Radreply) (res Radreply, err error)
	Delete(ctx context.Context, radreply *Radreply) (err error)
}
