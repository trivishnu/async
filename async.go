package async

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/netbookai/log"
	"github.com/netbookai/log/loggers"
	"github.com/pkg/errors"
)

func run(ctx context.Context, fn func(context.Context) error) error {
	errCh := make(chan error)
	defer func() {
		close(errCh)
	}()

	go func() {
		defer panicHandler(ctx, errCh)
		errCh <- fn(ctx)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func execute(ctx context.Context, task func(context.Context) error, taskName string, logger log.Logger) error {
	ctx = getContextWithTaskName(ctx, taskName)

	start := time.Now()
	logger.Info(ctx, "task started")

	defer func() {
		elapsed := time.Since(start)
		logger.Info(ctx, "task finished", "elapsed", elapsed.String())
	}()

	err := run(ctx, task)
	if err != nil {
		logger.Error(ctx, "error occured in executing task", "error", err)
		return err
	}

	return nil
}

func GoWithContext(ctx context.Context, taskName string, task func(context.Context) error, timeout time.Duration, logger log.Logger) {

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go execute(
		ctx,
		func(context.Context) error {
			defer func() {
				if ctx.Err() == nil {
					cancel()
				}
			}()
			return task(ctx)
		},
		taskName,
		logger)
}

func Go(taskName string, task func() error, timeout time.Duration, logger log.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	go execute(
		ctx,
		func(context.Context) error {
			defer func() {
				if ctx.Err() == nil {
					cancel()
				}
			}()
			return task()
		},
		taskName,
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

func GoWithWait(ctx context.Context,
	taskNames []string,
	tasks []func(context.Context) error,
	timeout time.Duration,
	logger log.Logger) error {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		if ctx.Err() == nil {
			cancel()
		}
	}()

	errs := make(chan error, len(tasks))

	wg := new(sync.WaitGroup)
	waitCh := make(chan struct{})
	wg.Add(len(tasks))

	go func() {
		for i, task := range tasks {
			taskToRun := task
			taskName := taskNames[i]
			go func() {
				defer wg.Done()
				taskErr := execute(ctx, taskToRun, taskName, logger)
				if taskErr != nil && ctx.Err() == nil {
					errs <- taskErr
					cancel()
				}
			}()
		}

		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh: // successfully completed all tasks
		close(waitCh)
		return nil
	case err := <-errs:
		close(errs)
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
