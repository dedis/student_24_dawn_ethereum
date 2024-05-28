package miner

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

type vdfWorker struct {
	cancel context.CancelFunc
	result *types.Transaction
	err error
	done chan struct{}
}

func startVdf(tx *types.Transaction, ctx context.Context) (*vdfWorker, error) {
	workerCtx, cancelFunc := context.WithCancel(ctx)
	w := &vdfWorker{
		cancel: cancelFunc,
		done: make(chan struct{}),
	}
	go w.work(workerCtx, tx)
	return w, nil
}

func (w *vdfWorker) work(ctx context.Context, tx *types.Transaction) {
	defer close(w.done)
	w.result, w.err = tx.Decrypt()
}

func (w *vdfWorker) Cancel() {
	w.cancel()
}

func (w *vdfWorker) Wait() (*types.Transaction, error) {
	<-w.done
	return w.result, w.err
}
