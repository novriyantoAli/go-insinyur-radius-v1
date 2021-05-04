package domain

import "context"

type Radcheck struct {
	ID        *int64  `json:"id"`
	Username  *string `json:"username"`
	Attribute *string `json:"attribute"`
	OP        *string `json:"op"`
	Value     *string `json:"value"`
}

type RadcheckRepository interface {
	Get(ctx context.Context, radcheck Radcheck) (res []Radcheck, err error)
	FetchWithValueExpiration(ctx context.Context, delete bool) (res []Radcheck, err error)
	Update(ctx context.Context, radcheck Radcheck) (err error)
	DeleteWithUsername(ctx context.Context, username string) (err error)
}
