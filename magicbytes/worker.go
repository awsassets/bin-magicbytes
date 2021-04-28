package magicbytes

import (
	"context"
	"fmt"
	"log"
	"sync"
)

func worker(ctx context.Context, cancel context.CancelFunc, jobs <-chan *searchContext, wg *sync.WaitGroup) {
	defer wg.Done()
loop:
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			if job == nil {
				continue loop
			}

			meta, err := searchMetasInAFile(ctx, job.filePath, job.metas)

			if err != nil {
				log.Println(fmt.Errorf("Meta search has failed path: %s, err: %v", job.filePath, err))

				continue
			}

			if meta != nil && !job.fn(job.filePath, meta.Type) {
				cancel()
			}
		}

	}
}
