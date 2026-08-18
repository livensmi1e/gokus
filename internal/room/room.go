package room

import (
	"context"
	"errors"
	"gokus/internal/blokus"
)

var (
	ErrClosed      = errors.New("room closed")
	ErrInvalidMove = errors.New("invalid move")
	ErrFull        = errors.New("room full")
)

type request interface {
	isRequest()
}

type placeRequest struct {
	pieceId int
	at      blokus.Coordinate
	reply   chan error
}

func (pr *placeRequest) isRequest() {}

type joinResult struct {
	client *Client
	err    error
}

type joinRequest struct {
	reply chan joinResult
}

func (pr *joinRequest) isRequest() {}

type Room struct {
	requests chan request

	cancel context.CancelFunc
	done   chan struct{}
}

func New(parent context.Context) *Room {
	ctx, cancel := context.WithCancel(parent)
	r := &Room{
		requests: make(chan request),
		cancel:   cancel,
		done:     make(chan struct{}),
	}
	go r.run(ctx)
	return r
}

func (r *Room) run(ctx context.Context) {
	defer close(r.done)
	game := blokus.NewDuoGame() // only goroutine run can access game. this is why game not placed in room struct
	joined := 0
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-r.requests:
			switch req := req.(type) {
			case *placeRequest:
				if game.PlacePiece(req.pieceId, req.at) {
					req.reply <- nil
					continue
				}
				req.reply <- ErrInvalidMove
			case *joinRequest:
				if joined == 0 {
					req.reply <- joinResult{
						client: &Client{blokus.Player1},
						err:    nil,
					}
					joined++
					continue
				}
				if joined == 1 {
					req.reply <- joinResult{
						client: &Client{blokus.Player2},
						err:    nil,
					}
					joined++
					continue
				}
				if joined >= 2 {
					req.reply <- joinResult{nil, ErrFull}
				}
			}
		}
	}
}

func (r *Room) Place(
	ctx context.Context,
	pieceId int,
	at blokus.Coordinate,
) error {
	reply := make(chan error, 1) // buffered to avoid room block forever because it send reply but Place caller already canceled
	req := &placeRequest{
		pieceId: pieceId,
		at:      at,
		reply:   reply,
	}
	// send request
	select {
	case r.requests <- req:
		// room received request
	case <-ctx.Done():
		// caller does not want to wait and command room to stop
		return ctx.Err()
	case <-r.done:
		// room has already stopped
		return ErrClosed
	}
	// receive reply
	select {
	case err := <-reply:
		return err
	case <-ctx.Done():
		// caller does not want to wait and command room to stop
		return ctx.Err()
	case <-r.done:
		// room has already stopped
		return ErrClosed
	}
}

func (r *Room) Join(ctx context.Context) (*Client, error) {
	reply := make(chan joinResult, 1)
	req := &joinRequest{
		reply: reply,
	}
	// send join request
	select {
	case r.requests <- req:
		// room received request
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-r.done:
		return nil, ErrClosed
	}
	// receive reply
	select {
	case result := <-reply:
		return result.client, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-r.done:
		return nil, ErrClosed
	}
}

// Safe to call multiple times
// r.Close()
// r.Close() second time is safe
func (r *Room) Close() {
	r.cancel()
	<-r.done // unblock when done chan is closed
}
