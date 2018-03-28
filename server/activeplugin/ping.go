package activeplugin

//ping需要root权限，（或加suid） 所以用exec.Command
import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"time"
)

func waitTimeout(c *exec.Cmd, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	done := make(chan error)
	go func() { done <- c.Wait() }()
	select {
	case err := <-done:
		timer.Stop()
		return err
	case <-timer.C:
		if err := c.Process.Kill(); err != nil {
			log.Printf("E! FATAL error killing process: %s", err)
			return err
		}
		// wait for the command to return after killing it
		<-done

		return errors.New("command exec timed out")
	}
}

func combinedOutputTimeout(c *exec.Cmd, timeout time.Duration) ([]byte, error) {
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b
	if err := c.Start(); err != nil {
		return nil, err
	}
	err := waitTimeout(c, timeout)
	return b.Bytes(), err
}

func HostPinger(timeout int, url string) (string, error) {
	var args []string
	args = append(args, "-i 0.2", "-c 4", url)
	bin, err := exec.LookPath("ping")
	if err != nil {
		return "", err
	}
	c := exec.Command(bin, args...)
	out, err := combinedOutputTimeout(c,
		time.Millisecond*time.Duration(timeout+1000))
	return string(out), err
}
