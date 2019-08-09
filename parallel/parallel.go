package parallel

import (
	"context"
	"sync"

	"github.com/ppltools/utils/recovery"
)

type DoWorkPieceFunc func(piece int)

// Parallelize is a very simple framework that allows for parallelizing
// N independent pieces of work.
func Parallelize(workers, pieces int, doWorkPiece DoWorkPieceFunc) {
	ParallelizeUntil(nil, workers, pieces, doWorkPiece)
}

// ParallelizeUntil is a framework that allows for parallelizing N
// independent pieces of work until done or the context is canceled.
func ParallelizeUntil(ctx context.Context, workers, pieces int, doWorkPiece DoWorkPieceFunc) {
	var stop <-chan struct{}
	if ctx != nil {
		stop = ctx.Done()
	}

	toProcess := make(chan int, pieces)
	for i := 0; i < pieces; i++ {
		toProcess <- i
	}
	close(toProcess)

	if pieces < workers {
		workers = pieces
	}

	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer recovery.Recovery(ctx, nil, "parallel")
			defer wg.Done()
			for piece := range toProcess {
				select {
				case <-stop:
					return
				default:
					doWorkPiece(piece)
				}
			}
		}()
	}
	wg.Wait()
}
