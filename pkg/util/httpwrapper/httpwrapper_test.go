package httpwrapper_test

import (
	"net/http"

	"carbonaut.dev/pkg/util/httpwrapper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("httpwrapper", func() {
	doesNotExistPath := "https://example.does-not-exist/"

	Describe("sending a http request to a valid endpoint", func() {
		It("should return no error using post", func() {
			resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
				Method:  http.MethodPost,
				BaseURL: "https://httpbin.org/post",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).To(Not(BeNil()))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
		It("should return no error using get", func() {
			resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
				Method:  http.MethodGet,
				BaseURL: "https://httpbin.org/get",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).To(Not(BeNil()))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Describe("sending a http request to an invalid endpoint", func() {
		It("should return an error using post", func() {
			_, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
				Method:  http.MethodPost,
				BaseURL: doesNotExistPath,
			})
			Expect(err).To(HaveOccurred())
		})
		It("should return an error using get", func() {
			_, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
				Method:  http.MethodGet,
				BaseURL: doesNotExistPath,
			})
			Expect(err).To(HaveOccurred())
		})
	})
})
