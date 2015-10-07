package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/omeid/gonzo/context"
)

var (
	pkgdir string
	imp    string
	cwd    string

	ctx = context.Background()
)

func init() {
	maxprocs := runtime.NumCPU()
	if maxprocs > 2 {
		runtime.GOMAXPROCS(maxprocs / 2)
	}

	var err error
	cwd, err = os.Getwd()
	if err != nil {
		ctx.Fatal(err)
	}

	gopaths := filepath.SplitList(os.Getenv("GOPATH"))

	if len(gopaths) == 0 {
		ctx.Fatal("Kar requires $GOAPTH to be set.")
	}

	for _, gopath := range gopaths {
		imp, err = filepath.Rel(filepath.Join(gopath, "src"), cwd)

		//package outside go path.
		if err != nil {
			ctx.Fatal("Dir outside go path.")
		}

		//package outside go path.
		if base := filepath.Base(imp); base == "." || base == ".." {
			continue
		}

		pkgdir = filepath.Join(gopath, "pkg", fmt.Sprintf("%s_%s_%s", runtime.GOOS, runtime.GOARCH, "kar"))
		return
	}

	ctx.Fatal("No Go Path found.")
}

func main() {

	files, err := filepath.Glob("*.go")
	if err != nil {
		ctx.Fatal(err)
	}
	if len(files) == 0 {
		ctx.Fatal("Not a go project.")
	}

	if !hasKargar(files) {
		ctx.Fatalf("No kar files. Found: %#v", files)
	}

	params := os.Args[1:]

	if len(params) == 0 {
		params = []string{"run"}
	}

	command, params := params[0], params[1:]

	args := []string{command, "-tags=kar", "-pkgdir=" + pkgdir}

	switch command {
	case "run":
		args = append(args, files...)
		args = append(args, params...)

	case "install":

		flags := flag.NewFlagSet("kar", flag.ExitOnError)

		var fv, fx bool
		flags.BoolVar(&fv, "v", false, "Print the name of packages as they are compiled.")
		flags.BoolVar(&fx, "x", false, "Print the commands.")

		flags.Parse(os.Args[2:])
		if fv {
			args = append(args, "-v")
		}

		if fx {
			args = append(args, "-x")
		}

		deps, err := deps(cwd, imp)
		if err != nil {
			ctx.Fatal(err)
		}

		args = append(args, deps...)

	default:
		ctx.Fatal("Invalid command. Expects `run` or `install`")
	}

	err = run(args)
	if err != nil {
		ctx.Fatal("BAD")
	}
}

func hasKargar(files []string) bool {
	for _, file := range files {
		if strings.HasSuffix("kar.go", file) {
			return true
		}
	}
	return false
}

func run(args []string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//Proxy OS Signals.
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range interrupts {
			cmd.Process.Signal(sig)
		}
	}()

	return cmd.Run()
}
