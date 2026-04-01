package blokus

type player struct {
	id      Occupant
	pieces  []*piece
	started bool
	stopped bool
}

func (p *player) validPieceID(id int) bool {
	return id >= 0 && id <= 20
}

func (p *player) hasPiece(id int) bool {
	pc, ok := p.getPiece(id)
	if !ok {
		return false
	}
	return !pc.used
}

func (p *player) getPiece(id int) (*piece, bool) {
	if !p.validPieceID(id) {
		return nil, false
	}
	return p.pieces[id], true
}

func (p *player) isPieceUsed(id int) bool {
	pc, ok := p.getPiece(id)
	if !ok {
		return false
	}
	return pc.used
}

func (p *player) markPieceUsed(id int) bool {
	if !p.hasPiece(id) {
		return false
	}
	p.pieces[id].used = true
	return true
}
