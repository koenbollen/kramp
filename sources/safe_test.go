package sources_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/koenbollen/kramp/sources"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestSafe(t *testing.T) {
	t.Parallel()

	a := &single{"a", nil}

	safe := &sources.Safe{
		Source: a,
	}

	got, err := safe.Query(context.Background(), "q")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestSafe_Error(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("panic")
	a := &single{"a", expectedError}

	obs, logs := observer.New(zap.InfoLevel)
	logger := zap.New(obs)

	safe := &sources.Safe{
		Source: a,
		Logger: logger,
	}

	result, err := safe.Query(context.Background(), "q")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("unexpected result: %v", result)
	}

	want := []observer.LoggedEntry{{
		Entry:   zapcore.Entry{Level: zap.ErrorLevel, Message: "safely failed to execute query"},
		Context: []zapcore.Field{zap.String("input", "q"), zap.Error(expectedError)},
	}}
	if !reflect.DeepEqual(logs.AllUntimed(), want) {
		t.Logf("message: %q", logs.AllUntimed()[0].Entry.Message)
		t.Errorf("unexpected logs: %+v", logs.AllUntimed())
	}
}
