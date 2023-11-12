/*
Copyright 2023 CARBONAUT AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httpwrapper_test

import (
	"net/http"

	"carbonaut.cloud/pkg/util/httpwrapper"

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
			Expect(err).To(BeNil())
			Expect(resp).To(Not(BeNil()))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
		It("should return no error using get", func() {
			resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
				Method:  http.MethodGet,
				BaseURL: "https://httpbin.org/get",
			})
			Expect(err).To(BeNil())
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
			Expect(err).To(Not(BeNil()))
		})
		It("should return an error using get", func() {
			_, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
				Method:  http.MethodGet,
				BaseURL: doesNotExistPath,
			})
			Expect(err).To(Not(BeNil()))
		})
	})
})
