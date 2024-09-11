package ros

import (
	"slices"

	"github.com/go-routeros/routeros/v3"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
)

type rosClient struct {
	CLI *routeros.Client
}

func NewROSRepository(rosc *routeros.Client) domain.FirewallRepository {
	return &rosClient{CLI: rosc}
}

func (r *rosClient) Find(params *domain.Firewall) (res []domain.Firewall, err error) {
	result, err := r.CLI.Run("/ip/firewall/mangle/print", "=.proplist=.id,chain,action,new-connection-mark,passthrough,src-mac-address,connection-mark,new-packet-mark")
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result.Re); i++ {
		firewall := domain.Firewall{
			ID:                result.Re[i].Map[".id"],
			Chain:             result.Re[i].Map["chain"],
			Action:            result.Re[i].Map["action"],
			NewConnectionMark: result.Re[i].Map["new-connection-mark"],
			Passthrough:       result.Re[i].Map["passthrough"],
			SrcMacAddress:     result.Re[i].Map["src-mac-address"],
			ConnectionMark:    result.Re[i].Map["connection-mark"],
			NewPacketMark:     result.Re[i].Map["new-packet-mark"],
			Comment:           result.Re[i].Map["comment"],
		}
		res = append(res, firewall)
	}

	resultIDs := slices.DeleteFunc(res, func(i domain.Firewall) bool { return params.ID != "" && params.ID == i.ID })
	if len(resultIDs) > 0 {
		return resultIDs, err
	}

	resultChain := slices.DeleteFunc(res, func(i domain.Firewall) bool { return params.Chain != "" && params.Chain == i.Chain })
	if len(resultChain) > 0 {
		return resultChain, err
	}

	resultAction := slices.DeleteFunc(res, func(i domain.Firewall) bool { return params.Action != "" && params.Action == i.Action })
	if len(resultAction) > 0 {
		return resultAction, err
	}

	resultNewConnectionMark := slices.DeleteFunc(res, func(i domain.Firewall) bool {
		return params.NewConnectionMark != "" && params.NewConnectionMark == i.NewConnectionMark
	})
	if len(resultNewConnectionMark) > 0 {
		return resultNewConnectionMark, err
	}

	resultPassthrough := slices.DeleteFunc(res, func(i domain.Firewall) bool { return params.Passthrough != "" && params.Passthrough == i.Passthrough })
	if len(resultPassthrough) > 0 {
		return resultPassthrough, err
	}

	resultSrcMacAddress := slices.DeleteFunc(res, func(i domain.Firewall) bool {
		return params.SrcMacAddress != "" && params.SrcMacAddress == i.SrcMacAddress
	})
	if len(resultSrcMacAddress) > 0 {
		return resultSrcMacAddress, err
	}

	resultConnectionMark := slices.DeleteFunc(res, func(i domain.Firewall) bool {
		return params.ConnectionMark != "" && params.ConnectionMark == i.ConnectionMark
	})
	if len(resultConnectionMark) > 0 {
		return resultConnectionMark, err
	}

	resultNewPacketMark := slices.DeleteFunc(res, func(i domain.Firewall) bool {
		return params.NewPacketMark != "" && params.NewPacketMark == i.NewPacketMark
	})
	if len(resultNewPacketMark) > 0 {
		return resultNewPacketMark, err
	}

	return

}

func (r *rosClient) First(param *domain.Firewall) (res domain.Firewall, err error) {
	firewalls := make([]domain.Firewall, 0)
	result, err := r.CLI.Run("/ip/firewall/mangle/print", "=.proplist=.id,chain,action,new-connection-mark,passthrough,src-mac-address,connection-mark,new-packet-mark")
	if err != nil {
		return domain.Firewall{}, err
	}

	for i := 0; i < len(result.Re); i++ {
		firewall := domain.Firewall{
			ID:                result.Re[i].Map[".id"],
			Chain:             result.Re[i].Map["chain"],
			Action:            result.Re[i].Map["action"],
			NewConnectionMark: result.Re[i].Map["new-connection-mark"],
			Passthrough:       result.Re[i].Map["passthrough"],
			SrcMacAddress:     result.Re[i].Map["src-mac-address"],
			ConnectionMark:    result.Re[i].Map["connection-mark"],
			NewPacketMark:     result.Re[i].Map["new-packet-mark"],
			Comment:           result.Re[i].Map["comment"],
		}
		firewalls = append(firewalls, firewall)

		if param.ID != "" && param.ID == firewall.ID {
			return firewall, nil
		}

		if param.Chain != "" && param.Chain == firewall.Chain {
			return firewall, nil
		}

		if param.Action != "" && param.Action == firewall.Action {
			return firewall, nil
		}

		if param.NewConnectionMark != "" && param.NewConnectionMark == firewall.NewConnectionMark {
			return firewall, nil
		}

		if param.SrcMacAddress != "" && param.SrcMacAddress == firewall.SrcMacAddress {
			return firewall, nil
		}

		if param.ConnectionMark != "" && param.ConnectionMark == firewall.ConnectionMark {
			return firewall, nil
		}

		if param.NewPacketMark != "" && param.NewPacketMark == firewall.NewPacketMark {
			return firewall, nil
		}
	}

	return

	// idx := slices.IndexFunc(firewalls, func(f domain.Firewall) bool { return f.ID != "" && param.ID != "" && f.ID == param.ID })
	// if idx != -1 {
	// 	return firewalls[idx], nil
	// }

	// idx = slices.IndexFunc(firewalls, func(f domain.Firewall) bool { return f.Chain != "" && param.Chain != "" && f.Chain == param.Chain })
	// if idx != -1 {
	// 	return firewalls[idx], nil
	// }

	// idx = slices.IndexFunc(firewalls, func(f domain.Firewall) bool { return f.Action != "" && param.Action != "" && f.Action == param.Action })
	// if idx != -1 {
	// 	return firewalls[idx], nil
	// }

	// idx = slices.IndexFunc(firewalls, func(f domain.Firewall) bool { return f.NewConnectionMark != "" && param.NewConnectionMark != "" && f.NewConnectionMark == param.NewConnectionMark })
	// if idx != -1 {
	// 	return firewalls[idx], nil
	// }

	// idx = slices.IndexFunc(firewalls, func(f domain.Firewall) bool  { return f.Passthrough == param.Passthrough })
	// if idx != -1 {
	// 	return firewalls[idx], nil
	// }

	// idx = slices.IndexFunc(firewalls, )
}

func (r *rosClient) Save(param *domain.Firewall) (err error) {
	params := make([]string, 0)
	params = append(params, "/ip/firewall/mangle/add")
	if param.Chain != "" {
		params = append(params, ("=chain=" + param.Chain))
	}
	if param.Action != "" {
		params = append(params, ("=action=" + param.Action))
	}
	if param.NewConnectionMark != "" {
		params = append(params, ("=new-connection-mark=" + param.NewConnectionMark))
	}
	if param.Passthrough != "" {
		params = append(params, ("=passthrough=" + param.Passthrough))
	}
	if param.SrcMacAddress != "" {
		params = append(params, ("=src-mac-address=" + param.SrcMacAddress))
	}
	if param.ConnectionMark != "" {
		params = append(params, ("=connection-mark=" + param.ConnectionMark))
	}
	if param.NewPacketMark != "" {
		params = append(params, ("=new-packet-mark=" + param.NewPacketMark))
	}
	if param.Comment != "" {
		params = append(params, "=comment="+param.Comment)
	}
	_, err = r.CLI.Run(params...)

	return
}

func (r *rosClient) Update(param *domain.Firewall) (err error) {
	params := make([]string, 0)
	params = append(params, "/ip/firewall/mangle/set")
	if param.ID != "" {
		params = append(params, ("=.id=" + param.ID))
	}
	if param.Chain != "" {
		params = append(params, ("=chain=" + param.Chain))
	}
	if param.Action != "" {
		params = append(params, ("=action=" + param.Action))
	}
	if param.NewConnectionMark != "" {
		params = append(params, ("=new-connection-mark=" + param.NewConnectionMark))
	}
	if param.Passthrough != "" {
		params = append(params, ("=passthrough=" + param.Passthrough))
	}
	if param.SrcMacAddress != "" {
		params = append(params, ("=src-mac-address=" + param.SrcMacAddress))
	}
	if param.ConnectionMark != "" {
		params = append(params, ("=connection-mark=" + param.ConnectionMark))
	}
	if param.NewPacketMark != "" {
		params = append(params, ("=new-packet-mark=" + param.NewPacketMark))
	}
	if param.Comment != "" {
		params = append(params, "=comment="+param.Comment)
	}
	_, err = r.CLI.Run(params...)
	return
}

func (r *rosClient) Delete(params *domain.Firewall) (err error) {
	id := "=.id=" + params.ID
	_, err = r.CLI.Run("/ip/firewall/mangle/remove", id)

	return
}
