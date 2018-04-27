package common

import (
	"bytes"
	"errors"
	"os/exec"
	"time"
)

//Command  run shell script
func Command(bin string, timeout time.Duration, args ...string) ([]byte, error) {
	bin, err := exec.LookPath(bin)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(bin, args...)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	wait := time.After(timeout)
	chErr := make(chan error)
	go func() {
		chErr <- cmd.Wait()
	}()
	select {
	case err = <-chErr:
	case <-wait:
		cmd.Process.Kill()
		err = errors.New("command exec timed out")
	}
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil

}
