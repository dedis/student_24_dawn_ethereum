package miner

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

type decryptionWorker struct {
	cancel context.CancelFunc
	result *types.Transaction
	err error
	done chan struct{}
}

func startDecryption(tx *types.Transaction, ctx context.Context) (*decryptionWorker, error) {
	workerCtx, cancelFunc := context.WithCancel(ctx)
	w := &decryptionWorker{
		cancel: cancelFunc,
		done: make(chan struct{}),
	}
	go w.work(workerCtx, tx)
	return w, nil
}

func (w *decryptionWorker) work(ctx context.Context, tx *types.Transaction) {
	defer close(w.done)
	w.result, w.err = tx.Decrypt()
}

func (w *decryptionWorker) Cancel() {
	w.cancel()
}

func (w *decryptionWorker) Wait() (*types.Transaction, error) {
	<-w.done
	return w.result, w.err
}
