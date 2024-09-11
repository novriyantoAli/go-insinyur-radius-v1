package domain

import (
	"context"
	"time"
)

// SchedulerUsecase ...
type SchedulerUsecase interface {
	GuardRadius(ctx context.Context, wita *time.Location) (err error)
	IncreasePerform(ctx context.Context, wita *time.Location) (err error)
	GuardMacBinding(ctx context.Context, wita *time.Location) (err error)
	GetUsers(ctx context.Context, delete bool) (resArr []Radcheck, err error)
	DeleteExpireUsers(ctx context.Context, username string) (err error)
	GetOnlineUsers(ctx context.Context, usernameList string) (res []Radacct, err error)
}
