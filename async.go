package async

import (
	"context"
	"fmt"
	"time"

	"github.com/netbookai/log"
	"github.com/netbookai/log/loggers"
	"github.com/pkg/errors"
)

func run(ctx context.Context, timeout time.Duration, fn func(context.Context)) error {
	errCh := make(chan error)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		cancel()
		close(errCh)
	}()

	go func() {
		defer panicHandler(ctx, errCh)
		fn(ctx)
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func execute(ctx context.Context, task func(context.Context), taskName string, timeout time.Duration, logger log.Logger) {
	ctx = getContextWithTaskName(ctx, taskName)

	start := time.Now()
	logger.Info(ctx, "task started")

	defer func() {
		elapsed := time.Since(start)
		logger.Info(ctx, "task finished", "elapsed", elapsed.String())
	}()

	err := run(ctx, timeout, task)
	if err != nil {
		logger.Error(ctx, "error occured in executing task", "error", err)
		return
	}
}

func GoWithContext(ctx context.Context, taskName string, task func(context.Context), timeout time.Duration, logger log.Logger) {
	go execute(
		ctx,
		task,
		taskName,
		timeout,
		logger)
}

func Go(taskName string, task func(), timeout time.Duration, logger log.Logger) {
	go execute(
		context.Background(),
		func(context.Context) {
			task()
		},
		taskName,
		timeout,
		logger)
}

func panicHandler(ctx context.Context, errCh chan error) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}

		err = errors.Wrap(err, "panic occured in task")

		//case: task can still run after timeout and panic somewhere
		//way to check closed channel after timeout
		if ctx.Err() == nil {
			errCh <- err
		}
	}
}

func getContextWithTaskName(ctx context.Context, taskName string) context.Context {
	return loggers.AddToLogContext(ctx, "taskname", taskName)
}
