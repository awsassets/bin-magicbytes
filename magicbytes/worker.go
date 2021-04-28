package magicbytes

import (
	"context"
	"log"
	"sync"
)

func worker(ctx context.Context, inner_ctx context.Context, cancel context.CancelFunc, jobs <-chan *searchContext, wg *sync.WaitGroup) {
loop:
	for {
		select {
		case <-ctx.Done():
			return
		case <-inner_ctx.Done():
			return
		case job := <-jobs:
			if job == nil {
				continue loop
			}
			wg.Add(1)

			meta, err := searchMetasInAFile(ctx, inner_ctx, job.filePath, job.metas)

			if err != nil {
				//TODO: Handle error?
				log.Println(err)
				wg.Done()

				continue
			}

			if meta != nil && !job.fn(job.filePath, meta.Type) {
				cancel()
			}

			wg.Done()
		}

	}
}
