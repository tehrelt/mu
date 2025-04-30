package dto

type Range struct {
	Min int64
	Max int64
}

func (r *Range) Nil() bool {
	return r.Min == 0 && r.Max == 0
}

type Pagination struct {
	Limit  uint64
	Offset uint64
}

func (p *Pagination) Nil() bool {
	return p.Limit == 0 && p.Offset == 0
}
