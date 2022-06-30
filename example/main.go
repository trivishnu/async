package main

import (
	"context"
	"time"

	"github.com/netbookai/log"
)

//long time take background task
func testTask(ctx context.Context, logger log.Logger) error {
	logger.Info(ctx, "test task started")
	time.Sleep(10 * time.Second) //intermediary steps
	logger.Info(ctx, "test task finished")

	return nil
}

func main() {

	//....
	wait := make(chan bool) //test channel just to on main routine

	GoWithContext_Example()

	GoWithWait_Example()

	//...

	<-wait
}
