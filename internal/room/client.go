package room

import "gokus/internal/blokus"

type Client struct {
	player blokus.Occupant
}

func (c Client) Player() blokus.Occupant {
	return c.player
}
