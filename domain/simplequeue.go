package domain

type SimpleQueue struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Target      string `json:"target"`
	PacketMarks string `json:"packet_marks"`
	MaxLimit    string `json:"max_limit"`
}

type SimpleQueueRepository interface {
	Find(param *SimpleQueue) (res []SimpleQueue, err error)
	First(param *SimpleQueue) (res SimpleQueue, err error)
	Save(param *SimpleQueue) (err error)
	Update(param *SimpleQueue) (err error)
	Delete(param *SimpleQueue) (err error)
}
