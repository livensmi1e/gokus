package room

import (
	"context"
	"gokus/internal/blokus"
)

type Client struct {
	player  blokus.Occupant
	room    *Room
	updates chan State
}

func newClient(r *Room, player blokus.Occupant) *Client {
	return &Client{
		player:  player,
		room:    r,
		updates: make(chan State, 1),
	}
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

func (c *Client) State(ctx context.Context) (State, error) {
	return c.room.state(ctx)
}

func (c *Client) Updates() <-chan State {
	return c.updates
}

func (c *Client) Leave(ctx context.Context) error {
	return c.room.leave(ctx)
}
