package room

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"gokus/internal/blokus"
)

func mustJoin(t *testing.T, r *Room) *Client {
	t.Helper()
	client, err := r.Join(context.Background())
	if err != nil {
		t.Fatalf("join room: %v", err)
	}
	if client == nil {
		t.Fatal("join room returned nil client without error")
	}
	return client
}

func mustReceiveState(t *testing.T, updates <-chan State) State {
	t.Helper()
	if updates == nil {
		t.Fatal("expected update channel, got nil")
	}
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	select {
	case state, ok := <-updates:
		if !ok {
			t.Fatal("update channel closed before delivering state")
		}
		return state
	case <-timer.C:
		t.Fatal("timed out waiting for state update")
	}
	return State{}
}

func TestRoomClose(t *testing.T) {
	ctx := context.Background()
	r := New(ctx)
	r.Close()
	select {
	case <-r.done:
		// done channel is closed
	default:
		t.Fatal("room should be closed already")
	}
}

func TestJoinAssignsFirstPlayer(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client := mustJoin(t, r)
	if client.Player() != blokus.Player1 {
		t.Fatalf("expected Player1, got %v", client.Player())
	}
}

func TestJoinAssignsSecondPlayer(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	_ = mustJoin(t, r)
	client2, err := r.Join(context.Background())
	if err != nil {
		t.Fatalf("join second player: %v", err)
	}
	if client2 == nil {
		t.Fatal("expected a client, got nil")
	}
	if client2.Player() != blokus.Player2 {
		t.Fatalf("expected Player2, got %v", client2.Player())
	}
}

func TestJoinRejectsWhenRoomFull(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	_ = mustJoin(t, r)
	_ = mustJoin(t, r)
	client, err := r.Join(context.Background())
	if !errors.Is(err, ErrFull) {
		t.Fatalf("expected ErrFull, got: %v", err)
	}
	if client != nil {
		t.Fatal("expected nil client")
	}
}

func TestClientPlace(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	_ = mustJoin(t, r)
	err := client1.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if err != nil {
		t.Fatalf("expected valid move, got: %v", err)
	}
}

func TestClientPlaceInvalidMove(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	_ = mustJoin(t, r)
	err := client1.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(0, 0),
	)
	if !errors.Is(err, ErrInvalidMove) {
		t.Fatalf("expected ErrInvalidMove, got %v", err)
	}
}

func TestClientPlaceAfterRoomClose(t *testing.T) {
	r := New(context.Background())
	client := mustJoin(t, r)
	r.Close()
	err := client.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("expected ErrClosed, got %v", err)
	}
}

// This test is flaky because room maybe wating for request when request is sent
// Making 2 option is valid for select in Place
// Go will randomly choose one to continue
// Therefor not a good test because it does not test the expected contract of the API
//
// After adding pre-check context cancelled in Place, this test is not flaky anymore
func TestClientPlaceWhenCallerContextAlreadyCancel(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := mustJoin(t, r)
	err := client.Place(
		ctx,
		0,
		blokus.NewCoordinate(4, 4),
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestClientPlaceRoomRejectPlayerOutOfTurn(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	_ = mustJoin(t, r)
	player2 := mustJoin(t, r)
	// precondition of test
	if player2.Player() != blokus.Player2 {
		t.Fatalf("expected Player2, got: %v", player2.Player())
	}
	err := player2.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(9, 9),
	)
	if !errors.Is(err, ErrOutOfTurn) {
		t.Fatalf("expected ErrOutOfTurn, got: %v", err)
	}
}

func TestClientWaitForOpponent(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	err := client1.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if !errors.Is(err, ErrWaitingForOpponent) {
		t.Fatalf("expected ErrWaitingForOpponent, got: %v", err)
	}
	_ = mustJoin(t, r)
	err = client1.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if err != nil {
		t.Fatalf("expected valid move after two players join, got: %v", err)
	}
}

func TestClientStateReflectsPlacedMove(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	client2 := mustJoin(t, r)
	err := client1.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if err != nil {
		t.Fatalf("expected valid move, got: %v", err)
	}
	state, err := client2.State(context.Background())
	if err != nil {
		t.Fatalf("get client state: %v", err)
	}
	if state.Board[4][4] != blokus.Player1 || state.CurrentPlayer != blokus.Player2 {
		t.Fatalf("expected player1 piece at (4,4) and current player is player2")
	}
}

func TestStateMutationDoesNotAffectRoom(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	client2 := mustJoin(t, r)
	err := client1.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if err != nil {
		t.Fatalf("expected valid move, got: %v", err)
	}
	state, err := client2.State(context.Background())
	if err != nil {
		t.Fatalf("get client state: %v", err)
	}
	if state.Board[4][4] != blokus.Player1 || state.CurrentPlayer != blokus.Player2 {
		t.Fatal("expected player1 piece at (4,4) and current player is player2")
	}
	state.Board[4][4] = blokus.Player2
	state.CurrentPlayer = blokus.Player1
	state, err = client2.State(context.Background())
	if err != nil {
		t.Fatalf("get client state: %v", err)
	}
	if state.Board[4][4] != blokus.Player1 || state.CurrentPlayer != blokus.Player2 {
		t.Fatal("expected room's game state does not change")
	}
}

// Subject:   Client
// Behavior:  Receives State
// Condition: After Opponent Places
func TestClientReceivesStateAfterOpponentPlaces(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	client2 := mustJoin(t, r)
	updates := client2.Updates()
	err := client1.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if err != nil {
		t.Fatalf("expected valid move, got: %v", err)
	}
	state := mustReceiveState(t, updates)
	if state.Board[4][4] != blokus.Player1 {
		t.Fatal("expected board at 4,4 to be Player1")
	}
	if state.CurrentPlayer != blokus.Player2 {
		t.Fatal("expected current player to be Player2")
	}
}

func TestSlowClientDoesNotBlockRoom(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	client2 := mustJoin(t, r)
	err := client1.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if err != nil {
		t.Fatalf("place first move: %v", err)
	}
	result := make(chan error, 1)
	go func() {
		result <- client2.Place(
			context.Background(),
			0,
			blokus.NewCoordinate(9, 9),
		)
	}()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("place second move: %v", err)
		}
	case <-timer.C:
		// timeout proves that Room has been blocked by updates channel
		// cleanup work
		// make sure the room does not block
		<-client1.Updates()
		<-client2.Updates()
		// make sure Place has finished
		<-result
		t.Fatal("place blocked by clients that had not consumed updates")
	}
}

func TestPublishLatestReplacePendingState(t *testing.T) {
	updates := make(chan State, 1)
	updates <- State{
		CurrentPlayer: blokus.Player1,
	}
	publishLatest(updates, State{
		CurrentPlayer: blokus.Player2,
	})
	got := <-updates
	if got.CurrentPlayer != blokus.Player2 {
		t.Fatalf(
			"expected latest state for Player2, got current player %v",
			got.CurrentPlayer,
		)
	}
}

func TestClientUpdatesCloseWhenRoomClose(t *testing.T) {
	r := New(context.Background())
	client := mustJoin(t, r)
	updates := client.Updates()
	r.Close()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-updates:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("updates channel did not close after room closed")
		}
	}
}

// be notice in code Room.Run() that
// if change reply join request first then publish state, this test maybe failed
func TestClientReceiveStateWhenOpponentJoins(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	updates := client1.Updates()
	_ = mustJoin(t, r)
	state := mustReceiveState(t, updates)
	if state.PlayerCount != 2 {
		t.Fatalf(
			"expected player count 2, got %d",
			state.PlayerCount,
		)
	}
}

func TestRoomClosesWhenClientLeaves(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	client2 := mustJoin(t, r)
	updates := client2.Updates()
	if err := client1.Leave(context.Background()); err != nil {
		t.Fatalf("leave room: %v", err)
	}
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-updates:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("updates channel did not close after room closed")
		}
	}
}

func TestRegistryJoinsClientsToSameRoom(t *testing.T) {
	registry := NewRegistry()
	defer registry.Close()
	client1, err := registry.Join(context.Background(), "alpha")
	if err != nil {
		t.Fatalf("join first client: %v", err)
	}
	client2, err := registry.Join(context.Background(), "alpha")
	if err != nil {
		t.Fatalf("join second client: %v", err)
	}
	if client1.Player() != blokus.Player1 {
		t.Fatalf("expected first client to be Player1, got %v", client1.Player())
	}
	if client2.Player() != blokus.Player2 {
		t.Fatalf("expected second client to be Player2, got %v", client2.Player())
	}
}

func TestRegistryJoinAfterClose(t *testing.T) {
	registry := NewRegistry()
	registry.Close()
	client, err := registry.Join(context.Background(), "alpha")
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("expected ErrClosed, got %v", err)
	}
	if client != nil {
		t.Fatal("expected nil client")
	}
}

func TestRegistryCloseIsIdempotent(t *testing.T) {
	registry := NewRegistry()
	registry.Close()
	registry.Close()
}

// todo
func TestRegistryConcurrentJoinsSameRoom(t *testing.T) {

}

func TestClientRotateUpdatesCurrentPieceShape(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client1 := mustJoin(t, r)
	_ = mustJoin(t, r)
	err := client1.Rotate(
		context.Background(),
		3,
	)
	if err != nil {
		t.Fatalf("rotate piece: %v", err)
	}
	state, err := client1.State(context.Background())
	if err != nil {
		t.Fatalf("get state: %v", err)
	}
	got := state.CurrentPieceShapes[3]
	want := []blokus.Coordinate{
		blokus.NewCoordinate(0, 0),
		blokus.NewCoordinate(1, 0),
		blokus.NewCoordinate(1, 1),
	}
	if !slices.Equal(got, want) {
		t.Fatalf(
			"expected rotated shape %v, got %v",
			want,
			got,
		)
	}
}
