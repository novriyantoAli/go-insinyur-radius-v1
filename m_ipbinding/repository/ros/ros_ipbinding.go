package ros

import (
	"github.com/go-routeros/routeros/v3"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
)

type rosClient struct {
	CLI *routeros.Client
}

func NewROSRepository(rosc *routeros.Client) domain.IPBindingRepository {
	return &rosClient{CLI: rosc}
}

func (r *rosClient) Find(params *domain.IPBinding) (res []domain.IPBinding, err error) {
	result, err := r.CLI.Run("/ip/hotspot/ip-binding/print", ("mac-address=" + params.MacAddress))
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result.Re); i++ {
		ipbinding := domain.IPBinding{
			ID:         result.Re[i].Map[".id"],
			MacAddress: result.Re[i].Map["mac-address"],
			Comment:    result.Re[i].Map["comment"],
			Type:       result.Re[i].Map["type"],
		}
		res = append(res, ipbinding)
	}

	return
}

func (r *rosClient) Finds(client *routeros.Client, param *domain.IPBinding) (res []domain.IPBinding, err error) {
	result, err := client.Run("/ip/hotspot/ip-binding/print", ("mac-address=" + param.MacAddress))
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result.Re); i++ {
		ipbinding := domain.IPBinding{
			ID:         result.Re[i].Map[".id"],
			MacAddress: result.Re[i].Map["mac-address"],
			Comment:    result.Re[i].Map["comment"],
			Type:       result.Re[i].Map["type"],
		}
		res = append(res, ipbinding)
	}
	return
}

func (r *rosClient) Firsts(client *routeros.Client, param *domain.IPBinding) (res domain.IPBinding, err error) {
	result, err := client.Run("/ip/hotspot/ip-binding/print", "=.proplist=.id,mac-address,type,comment")
	if err != nil {
		return domain.IPBinding{}, err
	}

	for i := 0; i < len(result.Re); i++ {
		ipbinding := domain.IPBinding{
			ID:         result.Re[i].Map[".id"],
			MacAddress: result.Re[i].Map["mac-address"],
			Type:       result.Re[i].Map["type"],
			Comment:    result.Re[i].Map["comment"],
		}
		if param.ID != "" && param.ID == ipbinding.ID {
			return ipbinding, nil
		}

		if param.MacAddress != "" && param.MacAddress == ipbinding.MacAddress {
			return ipbinding, nil
		}

		if param.Type != "" && param.Type == ipbinding.Type {
			return ipbinding, nil
		}

		if param.Comment != "" && param.Comment == ipbinding.Comment {
			return ipbinding, nil
		}
	}

	return
}

func (r *rosClient) First(param *domain.IPBinding) (res domain.IPBinding, err error) {
	result, err := r.CLI.Run("/ip/hotspot/ip-binding/print", "=.proplist=.id,mac-address,type,comment")
	if err != nil {
		return domain.IPBinding{}, err
	}

	for i := 0; i < len(result.Re); i++ {
		ipbinding := domain.IPBinding{
			ID:         result.Re[i].Map[".id"],
			MacAddress: result.Re[i].Map["mac-address"],
			Type:       result.Re[i].Map["type"],
			Comment:    result.Re[i].Map["comment"],
		}
		if param.ID != "" && param.ID == ipbinding.ID {
			return ipbinding, nil
		}

		if param.MacAddress != "" && param.MacAddress == ipbinding.MacAddress {
			return ipbinding, nil
		}

		if param.Type != "" && param.Type == ipbinding.Type {
			return ipbinding, nil
		}

		if param.Comment != "" && param.Comment == ipbinding.Comment {
			return ipbinding, nil
		}
	}

	return
}

func (r *rosClient) Saves(client *routeros.Client, params *domain.IPBinding) (err error) {
	tpe := "=type=" + params.Type
	mac := "=mac-address=" + params.MacAddress
	cmt := "=comment=" + params.Comment

	_, err = client.Run("/ip/hotspot/ip-binding/add", tpe, mac, cmt)

	return
}

func (r *rosClient) Save(params *domain.IPBinding) (err error) {
	tpe := "=type=" + params.Type
	mac := "=mac-address=" + params.MacAddress
	cmt := "=comment=" + params.Comment

	_, err = r.CLI.Run("/ip/hotspot/ip-binding/add", tpe, mac, cmt)

	return
}

func (r *rosClient) Updates(client *routeros.Client, params *domain.IPBinding) (err error) {
	id := "=.id=" + params.ID
	tpe := "=type=" + params.Type
	mac := "=mac-address=" + params.MacAddress
	cmt := "=comment=" + params.Comment

	_, err = client.Run("/ip/hotspot/ip-binding/set", id, tpe, mac, cmt)

	return
}

func (r *rosClient) Update(params *domain.IPBinding) (err error) {
	id := "=.id=" + params.ID
	tpe := "=type=" + params.Type
	mac := "=mac-address=" + params.MacAddress
	cmt := "=comment=" + params.Comment

	_, err = r.CLI.Run("/ip/hotspot/ip-binding/set", id, tpe, mac, cmt)

	return
}

func (r *rosClient) Deletes(client *routeros.Client, param *domain.IPBinding) (err error) {
	id := "=.id=" + param.ID
	_, err = r.CLI.Run("/ip/hotspot/ip-binding/remove", id)

	return
}

func (r *rosClient) Delete(params *domain.IPBinding) (err error) {
	id := "=.id=" + params.ID
	_, err = r.CLI.Run("/ip/hotspot/ip-binding/remove", id)

	return
}
