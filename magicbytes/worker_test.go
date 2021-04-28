package magicbytes

import (
	"context"
	"sync"
	"testing"
	"time"
)

func Test_worker(t *testing.T) {
	//Arrange
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan *searchContext)
	var wg sync.WaitGroup

	path := saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==")

	matchResult := false
	m := &Meta{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}}
	sc := &searchContext{filePath: path, metas: &[]*Meta{m}, wg: &wg, fn: func(path, metaType string) bool {
		matchResult = true
		return true
	}}

	type args struct {
		ctx       context.Context
		inner_ctx context.Context
		cancel    context.CancelFunc
		jobs      chan *searchContext
		wg        *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Init", args: args{ctx: ctx, inner_ctx: ctx, cancel: cancel, jobs: jobs, wg: &wg}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go worker(tt.args.ctx, tt.args.inner_ctx, tt.args.cancel, tt.args.jobs, tt.args.wg)

			tt.args.jobs <- sc

			wg.Wait()
			close(tt.args.jobs)

			if !matchResult {
				t.Errorf("worker() matchresult %v", matchResult)
				return
			}
		})
	}
}

func Test_worker_OnMatchReturnFalse(t *testing.T) {
	//Arrange
	ctx := context.Background()
	inner_ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan *searchContext, 2)
	var wg sync.WaitGroup

	path := saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==")
	hitCount := 0

	sc := &searchContext{filePath: path, metas: &[]*Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
	}, wg: &wg, fn: func(path, metaType string) bool {
		hitCount++
		return false
	}}

	t.Run("OnMatch function returns false", func(t *testing.T) {
		go worker(ctx, inner_ctx, cancel, jobs, &wg)

		jobs <- sc
		jobs <- sc

	loop:
		for {
			select {
			case <-inner_ctx.Done():
				break loop
			}
		}

		close(jobs)

		if hitCount != 1 {
			t.Errorf("worker() hitCount must be 1 but its %v", hitCount)
			return
		}

		if inner_ctx.Err() == nil {
			t.Errorf("worker() context %v", ctx.Err())
			return
		}
	})
}

func Test_worker_ContextCancel(t *testing.T) {
	//Arrange
	ctx, cancel := context.WithCancel(context.Background())
	inner_ctx := context.Background()
	defer cancel()

	jobs := make(chan *searchContext, 2)
	var wg sync.WaitGroup

	path := saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==")
	hitCount := 0

	sc := &searchContext{filePath: path, metas: &[]*Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
	}, wg: &wg, fn: func(path, metaType string) bool {
		hitCount++

		//Let's put a sleep for context cancellation.
		time.Sleep(2 * time.Second)

		return true
	}}

	t.Run("Context cancellation", func(t *testing.T) {
		go worker(ctx, inner_ctx, cancel, jobs, &wg)

		jobs <- sc
		jobs <- sc

		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()

	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			}
		}

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
