package main

import (
	"context"
	"time"

	"github.com/netbookai/async"
	"github.com/netbookai/log"
	"github.com/netbookai/log/loggers/zap"
)

func GoWithContext_Example() {

	// ...

	ctx := context.Background()
	logger := log.NewLogger(zap.NewLogger())
	timeout := 20 * time.Second

	// run as goroutine with timeout and error handlers
	async.New().Add(
		"testTask",
		func(taskCtx context.Context) error {
			return testTask(taskCtx, logger)
		},
	).GoWithContext(ctx, timeout, logger)

	// ... not waiting
}
