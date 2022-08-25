package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/deweppro/go-app/application/sys"
	"github.com/deweppro/go-app/console"
)

func main() {
	root := console.New("parallel", "Parallel command execution")
	root.RootCommand(RunApp())
	root.Exec()
}

func RunApp() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Flag(func(flags console.FlagsSetter) {
			flags.StringVar("shell", "bash", "default shell (bash, sh, ... etc)")
			flags.Bool("exit", "stop parallel if error or exit 1")
			flags.IntVar("timeout", 3, "restart timeout in sec if error or exit 1")
		})
		setter.ExecFunc(func(args []string, shell string, exit bool, timeout int64) {
			wg := sync.WaitGroup{}
			ctx, cncl := context.WithCancel(context.Background())
			go sys.OnSyscallStop(func() {
				cncl()
			})

			for i, arg := range args {
				wg.Add(1)
				go func(i int, arg string) {
					reTry(ctx, func() error {
						return execContext(ctx, shell, i, arg)
					}, exit, timeout)
					wg.Done()
				}(i+1, arg)
			}

			wg.Wait()
		})
	})
}

func execContext(ctx context.Context, shell string, i int, command string) error {
	cmd := exec.CommandContext(ctx, shell, "-c", command)
	cmd.Env = os.Environ()
	cmd.Dir = pwd()
	stdout, err := cmd.StdoutPipe()
	console.FatalIfErr(err, "stdout init")
	defer stdout.Close()
	console.FatalIfErr(cmd.Start(), "start command")
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			log(i, command, scanner.Text())
		}
	}()
	return cmd.Wait()
}

func reTry(ctx context.Context, call func() error, exit bool, timeout int64) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := call(); err != nil {
				console.Errorf(err.Error())
				if exit {
					return
				}
				time.Sleep(time.Second * time.Duration(timeout))
				continue
			}
			return
		}
	}
}

func log(i int, cmd, message string) {
	fmt.Printf("\u001B[%dm[%s]\t%s\n\u001B[0m", i+31, cmd, message)
}

func pwd() string {
	dir, err := os.Getwd()
	console.FatalIfErr(err, "get current dir")
	return dir
}
