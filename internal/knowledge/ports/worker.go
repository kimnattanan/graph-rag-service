package ports

import (
	"context"
	"time"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/command"
)

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc

	app *app.Application

	workerCount    int
	workerInterval time.Duration
}

func NewWorker(ctx context.Context, app *app.Application, workerCount int, workerInterval time.Duration) *Worker {
	ctx, cancel := context.WithCancel(ctx)
	return &Worker{
		ctx:            ctx,
		cancel:         cancel,
		app:            app,
		workerCount:    workerCount,
		workerInterval: workerInterval,
	}
}

func (w *Worker) RunWorkers() {
	if w.workerCount <= 0 {
		return
	}
	for i := 0; i < w.workerCount; i++ {
		go w.runWorker()
	}
	go w.runSweepOrphans()
}

func (w *Worker) runWorker() {
	ticker := time.NewTicker(w.workerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			for {
				if w.ctx.Err() != nil {
					return
				}
				noPending := false
				err := w.app.Commands.IndexNextDocument.Handle(w.ctx, command.IndexNextDocument{
					NoPending: &noPending,
				})
				if err != nil || noPending {
					break
				}
			}
		}
	}
}

func (w *Worker) runSweepOrphans() {
	ticker := time.NewTicker(w.workerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.app.Commands.SweepOrphans.Handle(w.ctx, command.SweepOrphans{})
		}
	}
}

func (w *Worker) Shutdown() {
	w.cancel()
}
