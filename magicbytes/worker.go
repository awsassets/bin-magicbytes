package magicbytes

import (
	"context"
	"log"
	"sync"
)

func worker(ctx context.Context, inner_ctx context.Context, cancel context.CancelFunc, jobs <-chan *searchContext, wg *sync.WaitGroup) {
	defer wg.Done()
loop:
	for {
		select {
		case <-ctx.Done():
			return
		case <-inner_ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			if job == nil {
				continue loop
			}

			meta, err := searchMetasInAFile(ctx, inner_ctx, job.filePath, job.metas)

			if err != nil {
				//TODO: Handle error?
				log.Println(err)

				continue
			}

			if meta != nil && !job.fn(job.filePath, meta.Type) {
				cancel()
			}

		}

	}
}
