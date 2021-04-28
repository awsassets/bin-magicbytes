package magicbytes

import (
	"context"
	"testing"
)

func Test_checkContextIsAlive(t *testing.T) {
	//Arrange
	ctx := context.Background()
	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Contextes are alive", func(t *testing.T) {
		if got := checkContextIsAlive(ctx, cctx); got != true {
			t.Errorf("checkContextIsAlive() = %v, want %v", got, true)
		}
	})
}

func Test_checkContextIsAlive_FirstIsDone(t *testing.T) {
	//Arrange
	ctx := context.Background()
	cctx, cancel := context.WithCancel(context.Background())

	t.Run("First one is done", func(t *testing.T) {
		cancel()
		if got := checkContextIsAlive(cctx, ctx); got != false {
			t.Errorf("checkContextIsAlive() = %v, want %v", got, false)
		}
	})
}

func Test_checkContextIsAlive_SecondIsDone(t *testing.T) {
	//Arrange
	ctx := context.Background()
	cctx, cancel := context.WithCancel(context.Background())

	t.Run("Second one is done", func(t *testing.T) {
		cancel()
		if got := checkContextIsAlive(ctx, cctx); got != false {
			t.Errorf("checkContextIsAlive() = %v, want %v", got, false)
		}
	})
}
