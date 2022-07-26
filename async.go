package async

import (
	"context"
	"fmt"
	"time"

	"github.com/netbookai/log"
	"github.com/netbookai/log/loggers"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type taskGroup struct {
	tasks []task
}

type task struct {
	name string
	fn   func(context.Context) error
}

func New() *taskGroup {
	return &taskGroup{}
}

func (g *taskGroup) Add(name string, fn func(context.Context) error) *taskGroup {
	g.tasks = append(g.tasks, task{name, fn})
	return g
}

func (g *taskGroup) GoWithContext(ctx context.Context, timeout time.Duration, logger log.Logger) {

	ctx, _ = context.WithTimeout(ctx, timeout)
	tasks := g.tasks

	for _, task := range tasks {
		taskToRun := task.fn
		taskName := task.name

		go execute(
			ctx,
			func(context.Context) error {
				return taskToRun(ctx)
			},
			taskName,
			logger)
	}
}

func (g *taskGroup) Go(timeout time.Duration, logger log.Logger) {
	g.GoWithContext(context.Background(), timeout, logger)
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

func (g *taskGroup) GoWithWait(ctx context.Context, timeout time.Duration, logger log.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		if ctx.Err() == nil {
			cancel()
		}
	}()

	tasks := g.tasks

	eg, ctx := errgroup.WithContext(ctx)

	for _, task := range tasks {
		taskToRun := task.fn
		taskName := task.name
		eg.Go(func() error {
			if taskErr := execute(ctx, taskToRun, taskName, logger); taskErr != nil {
				return taskErr
			}

			return nil
		})
	}

	return eg.Wait()
}

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
	newCtx := getContextWithTaskName(ctx, taskName)

	start := time.Now()
	logger.Info(newCtx, "task started")

	defer func() {
		elapsed := time.Since(start)
		logger.Info(newCtx, "task finished", "elapsed", elapsed.String())
	}()

	err := run(newCtx, task)
	if err != nil {
		logger.Error(newCtx, "error occured in executing task", "error", err)
		return err
	}

	return nil
}

func getContextWithTaskName(ctx context.Context, taskName string) context.Context {
	return loggers.GetNewLogContextWithValue(ctx, "taskname", taskName)
}
