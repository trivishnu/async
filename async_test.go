package async

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_run_method_with_success(t *testing.T) {
	timeout := 3 * time.Second
	task := func(context.Context) error {
		time.Sleep(1 * time.Second)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := run(ctx, task)
	if err != nil {
		t.Fatalf("error expected: %v and got: %v", nil, err)
	}
}

func Test_run_method_with_timeout(t *testing.T) {
	timeout := 2 * time.Second
	task := func(context.Context) error {
		time.Sleep(3 * time.Second)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := run(ctx, task)
	if err == nil {
		t.Fatalf("error expected: %v and got: %v", nil, err)
	}

	containgMsg := "deadline exceeded"
	assert.Containsf(t, err.Error(), containgMsg, "expected error msg must contain \"%s\"", containgMsg)
}

func Test_run_method_with_panic(t *testing.T) {
	timeout := 5 * time.Second
	task := func(context.Context) error {
		time.Sleep(3 * time.Second)
		panic("error occured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := run(ctx, task)
	if err == nil {
		t.Fatalf("error expected: %v and got: %v", nil, err)
	}

	containgMsg := "panic"
	assert.Containsf(t, err.Error(), containgMsg, "expected error msg must contain \"%s\"", containgMsg)
}
