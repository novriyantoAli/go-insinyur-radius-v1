package domain

import "context"

type ResellerBalance struct {
	Balance int64 `json:"balance"`
}

type ResellerRepository interface {
	Transaction(ctx context.Context, idUsers int64, username string, password string, packages Package) (idTransaction int64, err error)
}

type ResellerUsecase interface {
	Transaction(c context.Context, idUsers int64, idPackage int64) (res Transaction, err error)
	Balance(c context.Context, idUsers int64) (res ResellerBalance, err error)
	ChangePackage(c context.Context, voucher string, profile string) (err error)
	ChangeProfile(c context.Context, voucher string, profile string) (res string, err error)
}
