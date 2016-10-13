//Handle os signals in linux.
//Signal syscall.SIGUSR1: restart servers smoothly.
//Signal syscall.SIGTERM: stop servers and exit.

package grace

import (
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/boringding/beekeeper/proc"
)

func handleSignal() {
	var sig os.Signal

	signal.Notify(sigChan, syscall.SIGUSR1, syscall.SIGTERM)

	for {
		sig = <-sigChan
		switch sig {
		case syscall.SIGUSR1:
			sigusr1Handler()
		case syscall.SIGTERM:
			sigtermHandler()
		}
	}

	return
}

func sigusr1Handler() error {
	if isForked == true {
		return nil
	}

	isForked = true

	var files []*os.File

	i := 0
	for k, v := range GracefulSrvs {
		file, err := v.listener.file()
		if err != nil {
			continue
		}

		files = append(files, file)
		//The open file descriptor inherited by child process.
		//See cmd.ExtraFiles.
		envVal := 3 + i
		envValStr := strconv.Itoa(envVal)
		i++

		err = proc.SetEnv(k, envValStr)
		if err != nil {
			continue
		}
	}

	err := proc.SetEnv(BeekeeperChildEnv, "1")
	if err != nil {
		return err
	}

	exeFilePath := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	cmd := exec.Command(exeFilePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = files
	cmd.Env = os.Environ()

	return cmd.Start()
}

func sigtermHandler() error {
	if isClosed == true {
		return nil
	}

	isClosed = true

	var err error
	for _, v := range GracefulSrvs {
		e := v.shutdown()
		if e != nil {
			err = e
		}
	}

	return err
}
