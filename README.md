# Async library

Async library wrap long task and run as goroutine. Apart from running long background tasks it handles

* Timeout event while running task
* Panic in task
* Logs default Key Value saved in context (e.g. trace-id, taskname etc.) for better task tracking

## Installation

```
go get -u github.com/netbook-ai/async@v0.1.0
```

You may not be able to access the repo with netbook-devs path in GOPRIVATE,  update it as follows

```
export GOPRIVATE=gitlab.com/*
```

> can update it in your profile settings (.bashrc, .zshrc)

## Usage

``` async ``` supports two functions ``` GoWithContext ``` and ``` Go ```

### Using GoWithContext


```
  //....
	ctx := context.Background()
	logger := log.NewLogger(zap.NewLogger())
	timeout := 20 * time.Second

	// run as goroutine with timeout and error handlers
	async.GoWithContext(
		ctx,
		"testTask", // task name
		func(taskCtx context.Context) {
			testTask(taskCtx, logger) //time taking method to run as task
		},
		timeout,
		logger,
	)

  //...
  
```
### Sample Output

![Alt text](/images/output.jpg?raw=true "Optional Title")


### Using Go
```
  //....
	logger := log.NewLogger(zap.NewLogger())
	timeout := 20 * time.Second

	// run as goroutine with timeout and error handlers
	async.Go(
		"testTask", // task name
		func() {
			testTask(logger) //time taking method to run as task
		},
		timeout,
		logger,
	)

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
