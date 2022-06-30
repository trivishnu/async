# Async library

Async library supports following features to handle ong background tasks

* Running multiple tasks without Wait using ``` GoWithContext ```
* Running multiple tasks with Wait and error if any using ``` GoWithWait ```
* Running multiple tasks without Context and Wait using ``` Go ```

Apart from above featues it also handles

* Timeout event while running tasks
* Panic in tasks
* Logs default Key Value saved in context (e.g. trace-id, taskname etc.) for better task tracking

## Installation

```
go get -u github.com/netbookai/async@v0.2.0
```

## Usage

``` async ``` supports following methods ``` GoWithContext ```, ``` Go ```, ``` GoWithWait ```

### Using GoWithContext


```
  //....
	ctx := context.Background()
	logger := log.NewLogger(zap.NewLogger())
	timeout := 20 * time.Second

	// run as goroutine with timeout and error handlers
	async.New().Add(
		"testTask1",
		func(taskCtx context.Context) error {
			return testTask1(taskCtx, logger)
		},
	).Add(
		"testTask2",
		func(taskCtx context.Context) error {
			return testTask2(taskCtx, logger)
		},
	).GoWithContext(ctx, timeout, logger)

  //...
  
```
### Sample Output

![Alt text](/images/output.jpg?raw=true "Optional Title")

## Using GoWithWait

```
	ctx := context.Background()
	logger := log.NewLogger(zap.NewLogger())
	timeout := 20 * time.Second

	// run as goroutine with timeout and error handlers
	err := async.New().Add(
		"testTask1",
		func(taskCtx context.Context) error {
			return testTask2(taskCtx, logger)
		},
	).Add(
		"testTask2",
		func(taskCtx context.Context) error {
			return testTask2(taskCtx, logger)
		},
	).GoWithWait(ctx, timeout, logger)

	//waiting for tasks to finish

	if err != nil {
		logger.Error(ctx, "error occured", "error", err)
		return
	}

	//resuming next task
	//...
```

### Using Go
```
  //....
	logger := log.NewLogger(zap.NewLogger())
	timeout := 20 * time.Second

	// run as goroutine with timeout and error handlers
	
	async.New().Add(
		"testTask1",
		func() error {
			return testTask1(taskCtx, logger)
		},
	).Add(
		"testTask2",
		func(taskCtx context.Context) error {
			return testTask2(taskCtx, logger)
		},
	).Go(timeout, logger)

  //...
```

## Support
Raise a PR, we will get back to you.

## TODO

* add retry mechanism for task based by filtering permanent vs temporary errors


## Authors and acknowledgment

* Vishnu Tripathi <vishnutiwari612@gmail.com>

## License

   GNU GENERAL PUBLIC LICENSE
