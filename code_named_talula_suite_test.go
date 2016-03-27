package main_test

import (
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
)

var _ = BeforeSuite(func() {
	TestServer = rizo.CreateRequestRecordingServer(TestServerPort)
	TestServer.Start()
})

var _ = AfterSuite(func() {
	TestServer.Stop()
})
