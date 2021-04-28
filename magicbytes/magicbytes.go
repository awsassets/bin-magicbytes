package magicbytes

import (
	"context"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"sync"
)

type searchContext struct {
	filePath string
	metas    *[]*Meta
	fn       OnMatchFunc
	wg       *sync.WaitGroup
}

// Search searches the given target directory to find files recursively using meta information.
// For every match, onMatch callback is called concurrently.
func Search(ctx context.Context, targetDir string, metas []*Meta, onMatch OnMatchFunc) error {

	if targetDir == "" {
		return &ArgumentError{"targetDir is empty"}
	}

	if metas == nil || len(metas) > 1000 {
		return &ArgumentError{"metas has to be provided and must length be 0-100"}
	}

	if ctx == nil {
		return &ArgumentError{"ctx must be provided!"}
	}

	inner_ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	watcher := make(chan error)
	jobs := make(chan *searchContext)

	//We don't know the length of the recursive files in a dir and all goroutines has to be gracefully ended.
	var waitGroup sync.WaitGroup

	workerCount := runtime.GOMAXPROCS(0)
	//Spawn workers
	for i := 0; i < workerCount; i++ {
		go worker(ctx, inner_ctx, cancel, jobs, &waitGroup)
	}

	go func() {
		err := filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {

			if !checkContextIsAlive(ctx, inner_ctx) {
				return &ContextCancelledError{}
			}

			//TODO: handle individual walk error

			if info == nil || info.IsDir() {
				return nil
			}

			jobs <- &searchContext{filePath: path, metas: &metas, fn: onMatch, wg: &waitGroup}

			return nil
		})

		watcher <- err
	}()

	err := <-watcher

	if err != nil {
		log.Println(err)
	}

	//Ignore wait group when cancellation
	if checkContextIsAlive(ctx, inner_ctx) {
		waitGroup.Wait()
	}

	//Close channels
	close(watcher)
	close(jobs)

	return ctx.Err()
}
