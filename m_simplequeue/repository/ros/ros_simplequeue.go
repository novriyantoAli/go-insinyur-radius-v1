package ros

import (
	"slices"
	"strings"

	"github.com/go-routeros/routeros/v3"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
)

type rosClient struct {
	CLI *routeros.Client
}

func NewROSRepository(rosc *routeros.Client) domain.SimpleQueueRepository {
	return &rosClient{CLI: rosc}
}

func (r *rosClient) Find(params *domain.SimpleQueue) (res []domain.SimpleQueue, err error) {
	result, err := r.CLI.Run("/queue/simple/print", "=.proplist=.id,name,target,packet-marks,max-limit")
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result.Re); i++ {
		simpleQueue := domain.SimpleQueue{
			ID:          result.Re[i].Map[".id"],
			Name:        result.Re[i].Map["name"],
			Target:      result.Re[i].Map["target"],
			PacketMarks: result.Re[i].Map["packet-marks"],
			MaxLimit:    result.Re[i].Map["max-limit"],
		}
		res = append(res, simpleQueue)
	}

	resultIDs := slices.DeleteFunc(res, func(i domain.SimpleQueue) bool { return params.ID != "" && params.ID == i.ID })
	if len(resultIDs) > 0 {
		return resultIDs, err
	}

	resultName := slices.DeleteFunc(res, func(i domain.SimpleQueue) bool { return params.Name != "" && params.Name == i.Name })
	if len(resultName) > 0 {
		return resultName, err
	}
	resultPacketMarks := slices.DeleteFunc(res, func(i domain.SimpleQueue) bool {
		return params.PacketMarks != "" && params.PacketMarks == i.PacketMarks
	})
	if len(resultPacketMarks) > 0 {
		return resultPacketMarks, err
	}

	resultTarget := slices.DeleteFunc(res, func(i domain.SimpleQueue) bool { return params.Target != "" && params.Target == i.Target })
	if len(resultTarget) > 0 {
		return resultTarget, err
	}

	resultMaxLimit := slices.DeleteFunc(res, func(i domain.SimpleQueue) bool { return params.MaxLimit != "" && params.MaxLimit == i.MaxLimit })
	if len(resultMaxLimit) > 0 {
		return resultMaxLimit, err
	}

	return

}

func (r *rosClient) First(param *domain.SimpleQueue) (res domain.SimpleQueue, err error) {
	result, err := r.CLI.Run("/queue/simple/print", "=.proplist=.id,name,target,packet-marks,max-limit")
	if err != nil {
		return
	}

	for i := 0; i < len(result.Re); i++ {
		simpleQueue := domain.SimpleQueue{
			ID:          result.Re[i].Map[".id"],
			Name:        result.Re[i].Map["name"],
			Target:      result.Re[i].Map["target"],
			PacketMarks: result.Re[i].Map["packet-marks"],
			MaxLimit:    result.Re[i].Map["max-limit"],
		}

		if param.ID != "" && param.ID == simpleQueue.ID {
			return simpleQueue, nil
		}

		if param.Name != "" && param.Name == simpleQueue.Name {
			return simpleQueue, nil
		}

		if param.PacketMarks != "" && param.PacketMarks == simpleQueue.PacketMarks {
			return simpleQueue, nil
		}

		if param.MaxLimit != "" && param.MaxLimit == simpleQueue.MaxLimit {
			return simpleQueue, nil
		}
	}

	return

}

func (r *rosClient) Save(param *domain.SimpleQueue) (err error) {
	params := make([]string, 0)
	params = append(params, "/queue/simple/add")
	if param.Name != "" {
		params = append(params, ("=name=" + param.Name))
	}
	if param.Target != "" {
		params = append(params, ("=target=" + param.Target))
	}
	if param.PacketMarks != "" {
		params = append(params, ("=packet-marks=" + param.PacketMarks))
	}
	if param.MaxLimit != "" {
		params = append(params, ("=max-limit=" + strings.ToUpper(param.MaxLimit)))
	}
	_, err = r.CLI.Run(params...)

	return
}

func (r *rosClient) Update(param *domain.SimpleQueue) (err error) {
	params := make([]string, 0)
	params = append(params, "/queue/simple/set")
	if param.ID != "" {
		params = append(params, "=.id="+param.ID)
	}
	if param.Name != "" {
		params = append(params, ("=name=" + param.Name))
	}
	if param.Target != "" {
		params = append(params, ("=target=" + param.Target))
	}
	if param.PacketMarks != "" {
		params = append(params, ("=packet-marks=" + param.PacketMarks))
	}
	if param.MaxLimit != "" {
		params = append(params, ("=max-limit=" + strings.ToUpper(param.MaxLimit)))
	}
	_, err = r.CLI.Run(params...)

	return
}

func (r *rosClient) Delete(params *domain.SimpleQueue) (err error) {
	id := "=.id=" + params.ID
	_, err = r.CLI.Run("/queue/simple/remove", id)

	return
}
