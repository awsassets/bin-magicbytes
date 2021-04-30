package magicbytes

import (
	"context"
	"sync"
	"testing"
	"time"
)

func Test_worker(t *testing.T) {
	//Arrange
	ctxParent, cancelParent := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(ctxParent)

	defer cancel()
	defer cancelParent()

	jobs := make(chan string)
	var wg sync.WaitGroup

	path := saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==")

	matchResult := false
	metas := &[]*Meta{{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}}}
	onMatch := func(p, metaType string) bool {
		matchResult = true
		return true
	}

	tests := []struct {
		name string
	}{
		{name: "Init"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg.Add(1)
			go worker(ctx, cancel, jobs, metas, onMatch, &wg)

			jobs <- path
			close(jobs)

			wg.Wait()

			if !matchResult {
				t.Errorf("worker() matchresult %v", matchResult)
				return
			}
		})
	}
}

func Test_worker_OnMatchReturnFalse(t *testing.T) {
	//Arrange
	ctxParent, cancelParent := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(ctxParent)

	defer cancel()
	defer cancelParent()

	jobs := make(chan string, 2)
	var wg sync.WaitGroup

	path := saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==")
	hitCount := 0

	metas := &[]*Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
	}

	onMatch := func(path, metaType string) bool {
		hitCount++
		return false
	}

	t.Run("OnMatch function returns false", func(t *testing.T) {
		wg.Add(1)
		go worker(ctx, cancelParent, jobs, metas, onMatch, &wg)

		jobs <- path
		jobs <- path

		<-ctx.Done()

		close(jobs)

		if hitCount != 1 {
			t.Errorf("worker() hitCount must be 1 but its %v", hitCount)
			return
		}

		if ctx.Err() == nil {
			t.Errorf("worker() context %v", ctx.Err())
			return
		}
	})
}

func Test_worker_ContextCancel(t *testing.T) {
	//Arrange
	ctxParent, cancelParent := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(ctxParent)

	defer cancel()
	defer cancelParent()

	jobs := make(chan string, 2)
	var wg sync.WaitGroup

	path := saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==")
	hitCount := 0

	metas := &[]*Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
	}

	onMatch := func(path, metaType string) bool {
		hitCount++

		//Let's put a sleep for context cancellation.
		time.Sleep(2 * time.Second)

		return true
	}

	t.Run("Context cancellation", func(t *testing.T) {
		wg.Add(1)
		go worker(ctx, cancel, jobs, metas, onMatch, &wg)

		jobs <- path
		jobs <- path

		go func() {
			time.Sleep(1 * time.Second)
			cancelParent()
		}()

		<-ctxParent.Done()

		close(jobs)

		if hitCount != 1 {
			t.Errorf("worker() hitCount must be 1 but its %v", hitCount)
			return
		}

		if ctxParent.Err() == nil {
			t.Errorf("worker() context %v", ctx.Err())
			return
		}
	})
}
