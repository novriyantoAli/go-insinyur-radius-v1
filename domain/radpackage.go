package domain

type Radpackage struct {
	ID        *int64     `json:"id"`
	IDPackage *int64     `json:"id_package"`
	Username  *string    `json:"username"`
	Package   Package    `json:"package"`
	Radcheck  []Radcheck `json:"users"`
}
