package kar

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/omeid/gonzo/context"
	"github.com/omeid/kargar"
)

// Run setups a build and runs the listed tasks and cancels the
// build on a SIGTREM or INTERRUPT.
// It also calls os.Exit with appropriate code.
func Run(setup func(b *kargar.Build) error, cleanup func(int)) {
	//log.Flags = *level

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	b := kargar.NewBuild(ctx)

	err := setup(b)
	if err != nil {
		ctx.Fatal(err)
	}
	go func() {
		sig := <-interrupts
		// stop watches and clean up.
		fmt.Println() //Next line
		ctx.Warnf("Captured %v, stopping build and exiting...", sig)
		ctx.Warn("Press ctrl+c again to force exit.")
		cancel()
		ret := 0
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				ctx.Error(err)
				ret = 1
			}
		case <-interrupts:
			cancel()
			fmt.Println() //Next line
			ctx.Warn("Force exit.")
			ret = 1
		}
		os.Exit(ret)

	}()

	var wg sync.WaitGroup

	tasks := []string{"default"}
	if len(os.Args) > 1 {
		tasks = os.Args[1:]
	}

	ctx.Info(tasks)
	var ret uint32
	for _, t := range tasks {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			err := b.Run(t)
			if err != nil {
				atomic.StoreUint32(&ret, 1)
				ctx.Error(err)
			}
		}(t)
	}

	wg.Wait()
	//XXX: atomic operation uncessary?
	retcode := int(atomic.LoadUint32(&ret))
	if cleanup != nil {
		cleanup(retcode)
	}
	os.Exit(retcode)
}
