// package cmd_test

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"reflect"
// 	"strings"
// 	"testing"

// 	"github.com/spf13/pflag"
// 	"github.com/spf13/cobra"
// )

// func TestSingleCommand(t *testing.T) {
// 	var rootCmdArgs []string
// 	rootCmd := &Command{
// 		Use:  "root",
// 		Args: ExactArgs(2),
// 		Run:  func(_ *Command, args []string) { rootCmdArgs = args },
// 	}
// 	aCmd := &Command{Use: "a", Args: NoArgs, Run: emptyRun}
// 	bCmd := &Command{Use: "b", Args: NoArgs, Run: emptyRun}
// 	rootCmd.AddCommand(aCmd, bCmd)

// 	output, err := executeCommand(rootCmd, "one", "two")
// 	if output != "" {
// 		t.Errorf("Unexpected output: %v", output)
// 	}
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	got := strings.Join(rootCmdArgs, " ")
// 	if got != onetwo {
// 		t.Errorf("rootCmdArgs expected: %q, got: %q", onetwo, got)
// 	}
// }

// func TestChildCommand(t *testing.T) {
// 	var child1CmdArgs []string
// 	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
// 	child1Cmd := &Command{
// 		Use:  "child1",
// 		Args: ExactArgs(2),
// 		Run:  func(_ *Command, args []string) { child1CmdArgs = args },
// 	}
// 	child2Cmd := &Command{Use: "child2", Args: NoArgs, Run: emptyRun}
// 	rootCmd.AddCommand(child1Cmd, child2Cmd)

// 	output, err := executeCommand(rootCmd, "child1", "one", "two")
// 	if output != "" {
// 		t.Errorf("Unexpected output: %v", output)
// 	}
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	got := strings.Join(child1CmdArgs, " ")
// 	if got != onetwo {
// 		t.Errorf("child1CmdArgs expected: %q, got: %q", onetwo, got)
// 	}
// }
