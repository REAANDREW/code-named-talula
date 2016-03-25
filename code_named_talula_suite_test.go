package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCodeNamedTalula(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CodeNamedTalula Suite")
}
