package domain

type Firewall struct {
	ID                string `json:"id"`
	Chain             string `json:"chain"`
	Action            string `json:"action"`
	NewConnectionMark string `json:"new_connection_mark"`
	Passthrough       string `json:"passthrough"`
	SrcMacAddress     string `json:"src_mac_address"`
	ConnectionMark    string `json:"connection_mark"`
	NewPacketMark     string `json:"new_packet_mark"`
	Comment           string `json:"comment"`
}

type FirewallRepository interface {
	Find(param *Firewall) (res []Firewall, err error)
	First(param *Firewall) (res Firewall, err error)
	Save(param *Firewall) (err error)
	Update(param *Firewall) (err error)
	Delete(param *Firewall) (err error)
}
