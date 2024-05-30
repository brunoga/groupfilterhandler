package groupfilterhandler

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestGroupFilterHandler_No_Allowed_groups(t *testing.T) {
	var b bytes.Buffer
	logger := testLogger(&b)

	logger.Info("test")
	checkAndResetBuffer(t, &b, "level=INFO msg=test\n")
}
func TestGroupFilterHandler_Allowed_Groups(t *testing.T) {
	var b bytes.Buffer
	disallowedLogger := testLogger(&b, "allow_group_1", "allow_group_2")

	allowedLogger1 := disallowedLogger.WithGroup("allow_group_1")
	allowedLogger2 := disallowedLogger.WithGroup("allow_group_2")

	disallowedLogger.Info("disallowed")
	checkAndResetBuffer(t, &b, "")

	allowedLogger1.Info("allowed 1")
	checkAndResetBuffer(t, &b, "level=INFO msg=\"allowed 1\"\n")

	allowedLogger2.Info("allowed 2")
	checkAndResetBuffer(t, &b, "level=INFO msg=\"allowed 2\"\n")
}

func TestGroupFilterHandler_Logger_With(t *testing.T) {
	var b bytes.Buffer
	allowedLogger := testLogger(&b, "allow_group_1").WithGroup("allow_group_1").With(
		"extra_key", "extra_value")

	allowedLogger.Info("allowed")
	checkAndResetBuffer(t, &b, "level=INFO msg=allowed allow_group_1.extra_key=extra_value\n")
}

func testLogger(b *bytes.Buffer, allowedGroups ...string) *slog.Logger {
	return slog.New(New(slog.NewTextHandler(b, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				// Drop time attribute to make the test easier.
				return slog.Attr{}
			}

			return a
		},
	}), allowedGroups...))
}

func checkAndResetBuffer(t *testing.T, b *bytes.Buffer, want string) {
	t.Helper()

	if got := b.String(); got != want {
		t.Errorf("log output = %q; want %q", got, want)
	}
	b.Reset()
}
