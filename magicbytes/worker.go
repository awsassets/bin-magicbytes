package magicbytes

import (
	"context"
	"log"
	"sync"
)

func worker(ctx context.Context, cancel context.CancelFunc, jobs <-chan string, metas *[]*Meta, onMatch OnMatchFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			if job == "" {
				continue
			}

			meta, err := searchMetasInAFile(ctx, job, metas)

			if err != nil {
				log.Printf("Meta search has failed path: %s, err: %v", job, err)

				continue
			}

			if meta != nil && !onMatch(job, meta.Type) {
				cancel()
			}
		}

	}
}
