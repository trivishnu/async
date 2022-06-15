package main

import (
	"context"
	"time"

	"github.com/netbook-ai/async"
	"github.com/netbookai/log"
	"github.com/netbookai/log/loggers/zap"
)

//long time take background task
func testTask(ctx context.Context, logger log.Logger) {
	logger.Info(ctx, "test task started")
	time.Sleep(10 * time.Second) //intermediary steps
	logger.Info(ctx, "test task finished")
}

func main() {

	//....
	wait := make(chan bool) //test channel just to on main routine

	ctx := context.Background()
	logger := log.NewLogger(zap.NewLogger())
	timeout := 20 * time.Second

	// run as goroutine with timeout and error handlers
	async.GoWithContext(
		ctx,
		"testTask", // task name
		func(taskCtx context.Context) {
			testTask(taskCtx, logger)
		},
		timeout,
		logger,
	)

	//...

	<-wait
}
