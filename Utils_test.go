package main_test

import (
	. "github.com/reaandrew/code-named-talula"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	It("AdminURL", func() {
		expected := "http://localhost:38765/endpoints/123/transforms/response"
		actual := AdminURL("/endpoints/%d/transforms/response", 123)
		Expect(actual).To(Equal(expected))
	})
	It("TransformURL", func() {
		expected := "http://localhost:48765/endpoints/123/transforms/response"
		actual := TransformURL("/endpoints/%d/transforms/response", 123)
		Expect(actual).To(Equal(expected))
	})
})
