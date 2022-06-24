package main

import (
	"context"
	"time"

	"github.com/netbookai/async"
	"github.com/netbookai/log"
	"github.com/netbookai/log/loggers/zap"
)

func GoWithWait_Example() {
	ctx := context.Background()
	logger := log.NewLogger(zap.NewLogger())
	timeout := 20 * time.Second

	// run as goroutine with timeout and error handlers
	err := async.New().Add(
		"testTask",
		func(taskCtx context.Context) error {
			return testTask(taskCtx, logger)
		},
	).GoWithWait(ctx, timeout, logger)

	//waiting for tasks to finish

	if err != nil {
		logger.Error(ctx, "error occured", "error", err)
		return
	}

	//resuming next task
	//...
}
