package blokus

type player struct {
	id      Occupant
	pieces  []*piece
	started bool
	stopped bool
}

func (p *player) hasPiece(id int) bool {
	if id < 0 || id > 20 {
		return false
	}
	return p.pieces[id] != nil
}

func (p *player) getPiece(id int) (*piece, bool) {
	if !p.hasPiece(id) {
		return nil, false
	}
	return p.pieces[id], true
}

func (p *player) takePiece(id int) (*piece, bool) {
	if !p.hasPiece(id) {
		return nil, false
	}
	pc := p.pieces[id]
	p.pieces[id] = nil
	return pc, true
}
