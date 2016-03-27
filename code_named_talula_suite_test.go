package main_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCodeNamedTalula(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CodeNamedTalula Suite")
}

var (
	TestServerPort = 4000
	TestServer     *rizo.RequestRecordingServer
	Cmd            *exec.Cmd
	Stdout         io.ReadCloser
	Stderr         io.ReadCloser
)

func StartApplication() {
	exePath, err := filepath.Abs("./code-named-talula")
	if err != nil {
		panic(err)
	}

	Cmd = exec.Command(exePath)
	Stdout, err = Cmd.StdoutPipe()
	Expect(err).To(BeNil())
	Stderr, err = Cmd.StderrPipe()
	Expect(err).To(BeNil())
	Cmd.Start()
}

func StopApplication() {
	Cmd.Process.Kill()

	stdoutOutput, err := ioutil.ReadAll(Stdout)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----------STDOUT")
	fmt.Println(string(stdoutOutput))

	stderrOutput, err := ioutil.ReadAll(Stderr)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----------STDERR")
	fmt.Println(string(stderrOutput))
}

var _ = BeforeSuite(func() {
	TestServer = rizo.CreateRequestRecordingServer(TestServerPort)
	TestServer.Start()
	StartApplication()
	time.Sleep(time.Second * 1)
})

var _ = AfterSuite(func() {
	TestServer.Stop()
	StopApplication()
})
