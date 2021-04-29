package magicbytes

import (
	"context"
	"fmt"
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

	//Context propagation would be enough for cancelling false result in OnMatchFunc
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan *searchContext)

	//We don't know the length of the recursive files in a dir and all goroutines has to be gracefully ended.
	var waitGroup sync.WaitGroup

	workerCount := runtime.GOMAXPROCS(0)

	//Spawn workers
	waitGroup.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(ctx, cancel, jobs, &waitGroup)
	}

	go func() {
		err := filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				//We don't want to break on individual walk errors.
				log.Println(fmt.Errorf("File walk error on path: %s err: %v", path, err))

				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if d == nil || !d.Type().IsRegular() {
				return nil
			}

			jobs <- &searchContext{filePath: path, metas: &metas, fn: onMatch, wg: &waitGroup}

			return nil
		})

		if err != nil {
			log.Println(fmt.Errorf("File walk error: %v", err))
		}

		close(jobs)
	}()

	waitGroup.Wait()

	return ctx.Err()
}
