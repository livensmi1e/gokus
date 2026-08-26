package room

import (
	"context"
	"errors"
	"gokus/internal/blokus"
)

var (
	ErrClosed             = errors.New("room closed")
	ErrInvalidMove        = errors.New("invalid move")
	ErrFull               = errors.New("room full")
	ErrOutOfTurn          = errors.New("out of turn")
	ErrWaitingForOpponent = errors.New("waiting for opponent")
)

type request interface {
	isRequest()
}

type placeRequest struct {
	player  blokus.Occupant
	pieceId int
	at      blokus.Coordinate
	reply   chan error
}

func (*placeRequest) isRequest() {}

type joinResult struct {
	client *Client
	err    error
}

type joinRequest struct {
	reply chan joinResult
}

func (*joinRequest) isRequest() {}

type stateRequest struct {
	reply chan State
}

func (*stateRequest) isRequest() {}

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
	game := blokus.NewDuoGame() // only goroutine run can access game. this is why game not placed in room struct
	clients := make([]*Client, 0, 2)
	defer func() {
		for _, client := range clients {
			close(client.updates)
		}
		close(r.done)
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-r.requests:
			switch req := req.(type) {
			case *placeRequest:
				if len(clients) < 2 {
					req.reply <- ErrWaitingForOpponent
					continue
				}
				if req.player != game.CurrentPlayer() {
					req.reply <- ErrOutOfTurn
					continue
				}
				if game.PlacePiece(req.pieceId, req.at) {
					for _, client := range clients {
						publishLatest(client.updates, State{
							Board:         game.Board(),
							CurrentPlayer: game.CurrentPlayer(),
							PlayerCount:   len(clients),
						})
					}
					req.reply <- nil
					continue
				}
				req.reply <- ErrInvalidMove
			case *joinRequest:
				if len(clients) >= 2 {
					req.reply <- joinResult{err: ErrFull}
					continue
				}
				var client *Client
				if len(clients) == 0 {
					client = newClient(r, blokus.Player1)

				}
				if len(clients) == 1 {
					client = newClient(r, blokus.Player2)
				}
				clients = append(clients, client)
				// publish first
				for _, client := range clients {
					publishLatest(client.updates, State{
						Board:         game.Board(),
						CurrentPlayer: game.CurrentPlayer(),
						PlayerCount:   len(clients),
					})
				}
				// then return to join request
				req.reply <- joinResult{
					client: client,
					err:    nil,
				}
			case *stateRequest:
				req.reply <- State{
					Board:         game.Board(),
					CurrentPlayer: game.CurrentPlayer(),
					PlayerCount:   len(clients),
				}
			}
		}
	}
}

func (r *Room) place(
	ctx context.Context,
	player blokus.Occupant,
	pieceId int,
	at blokus.Coordinate,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	reply := make(chan error, 1) // buffered to avoid room block forever because it send reply but Place caller already canceled
	req := &placeRequest{
		player:  player,
		pieceId: pieceId,
		at:      at,
		reply:   reply,
	}
	// send request
	select {
	case r.requests <- req:
		// room received request
	case <-ctx.Done():
		// caller does not want to wait and command room to close
		return ctx.Err()
	case <-r.done:
		// room has already closed
		return ErrClosed
	}
	// receive reply
	// always wait the operation to be completed
	return <-reply
}

func (r *Room) state(ctx context.Context) (State, error) {
	if err := ctx.Err(); err != nil {
		return State{}, err
	}
	reply := make(chan State, 1)
	req := &stateRequest{reply: reply}
	select {
	case r.requests <- req:
		// room received request
	case <-ctx.Done():
		return State{}, ctx.Err()
	case <-r.done:
		return State{}, ErrClosed
	}
	return <-reply, nil
}

func (r *Room) Join(ctx context.Context) (*Client, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	reply := make(chan joinResult, 1)
	req := &joinRequest{reply: reply}
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
	// always wait the operation to be completed
	result := <-reply
	return result.client, result.err

}

// Safe to call multiple times
// r.Close()
// r.Close() second time is safe
func (r *Room) Close() {
	r.cancel()
	<-r.done // unblock when done chan is closed
}

func publishLatest(updates chan State, s State) {
	// try publish state
	select {
	case updates <- s:
		// success
		return
	default:
		// channel full
	}
	// try recevie old state
	select {
	case <-updates:
		// success
	default:
		// old state has been received by client
	}
	// publish again, channel always empty so always success
	updates <- s
}
