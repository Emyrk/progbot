package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var nodepath = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "Emyrk", "fortnitediscord", "matchwatcher", "gonode", "index.js")

var _ = fmt.Println

func ExecuteElixirWithTimeout(code string, timeout int) (string, error) {
	c1 := make(chan string, 1)
	e1 := make(chan error, 1)
	go func() {
		res, err := ExecuteElixir(code)
		if err != nil {
			e1 <- err
			return
		}
		c1 <- res

	}()

	// Here's the `select` implementing a timeout.
	// `res := <-c1` awaits the result and `<-Time.After`
	// awaits a value to be sent after the timeout of
	// 1s. Since `select` proceeds with the first
	// receive that's ready, we'll take the timeout case
	// if the operation takes more than the allowed 1s.
	select {
	case res := <-c1:
		return res, nil
	case err := <-e1:
		return "", err
	case <-time.After(time.Duration(timeout) * time.Second):
		return "", fmt.Errorf("Timeout on node execute")

	}
	return "", fmt.Errorf("impossible")
}

func ExecuteElixir(code string) (string, error) {
	os.Remove("code.exs")
	file, err := os.OpenFile("code.exs", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return "", err
	}
	file.Write([]byte(code))
	file.Close()
	data := exec.Command("elixir", "code.exs")

	data.Wait()
	d, err := data.Output()
	// os.Remove("code.exs")
	if err != nil {
		fmt.Println("Error output:", string(d))
		return "", err
	}

	return string(d), nil
}

func ExecuteNodeWithTimeout(code string, timeout int) (string, error) {
	c1 := make(chan string, 1)
	e1 := make(chan error, 1)
	go func() {
		res, err := ExecuteNode(code)
		if err != nil {
			e1 <- err
			return
		}
		c1 <- res

	}()

	// Here's the `select` implementing a timeout.
	// `res := <-c1` awaits the result and `<-Time.After`
	// awaits a value to be sent after the timeout of
	// 1s. Since `select` proceeds with the first
	// receive that's ready, we'll take the timeout case
	// if the operation takes more than the allowed 1s.
	select {
	case res := <-c1:
		return res, nil
	case err := <-e1:
		return "", err
	case <-time.After(time.Duration(timeout) * time.Second):
		return "", fmt.Errorf("Timeout on node execute")

	}
	return "", fmt.Errorf("impossible")
}

func ExecuteNode(code string) (string, error) {
	os.Remove("code.js")
	file, err := os.OpenFile("code.js", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return "", err
	}
	file.Write([]byte(code))
	file.Close()
	data := exec.Command("node", "code.js")

	data.Wait()
	d, err := data.Output()
	// os.Remove("code.exs")
	if err != nil {
		fmt.Println("Error output:", string(d))
		return "", err
	}

	return string(d), nil
}
