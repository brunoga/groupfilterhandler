package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/brunoga/groupfilterhandler"
)

type stringList []string

// Implement the flag.Value interface
func (s *stringList) String() string {
	return fmt.Sprint(*s)
}

func (s *stringList) Set(value string) error {
	if len(value) > 0 {
		*s = strings.Split(value, ",")
	}

	return nil
}

func main() {
	var logGroups stringList

	flag.Var(&logGroups, "log-groups", "Comma-separated list of log groups to "+
		"allow (empty means allow all groups)")

	flag.Parse()

	l := slog.New(groupfilterhandler.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{}), logGroups...)).WithGroup("allowed")

	l.Error("This will only show if \"allowed\" is listed in the -log-groups " +
		"flag (or if the flag is empty)")
}
