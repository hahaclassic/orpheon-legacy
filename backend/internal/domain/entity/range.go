package entity

type Range struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func (r Range) Len() int {
	return r.End - r.Start
}
