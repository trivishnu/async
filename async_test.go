package async

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_run_method_with_success(t *testing.T) {
	timeout := 3 * time.Second
	task := func(context.Context) {
		time.Sleep(1 * time.Second)
	}

	err := run(context.Background(), timeout, task)
	if err != nil {
		t.Fatalf("error expected: %v and got: %v", nil, err)
	}
}

func Test_run_method_with_timeout(t *testing.T) {
	timeout := 2 * time.Second
	task := func(context.Context) {
		time.Sleep(3 * time.Second)
	}

	err := run(context.Background(), timeout, task)
	if err == nil {
		t.Fatalf("error expected: %v and got: %v", nil, err)
	}

	containgMsg := "deadline exceeded"
	assert.Containsf(t, err.Error(), containgMsg, "expected error msg must contain \"%s\"", containgMsg)
}

func Test_run_method_with_panic(t *testing.T) {
	timeout := 5 * time.Second
	task := func(context.Context) {
		time.Sleep(3 * time.Second)
		panic("error occured")
	}

	err := run(context.Background(), timeout, task)
	if err == nil {
		t.Fatalf("error expected: %v and got: %v", nil, err)
	}

	containgMsg := "panic"
	assert.Containsf(t, err.Error(), containgMsg, "expected error msg must contain \"%s\"", containgMsg)
}
