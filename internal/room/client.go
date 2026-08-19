package room

import (
	"context"
	"gokus/internal/blokus"
)

type Client struct {
	player blokus.Occupant
	room   *Room
}

func (c *Client) Player() blokus.Occupant {
	return c.player
}

func (c *Client) Place(
	ctx context.Context,
	pieceId int,
	at blokus.Coordinate,
) error {
	return c.room.place(ctx, c.player, pieceId, at)
}
