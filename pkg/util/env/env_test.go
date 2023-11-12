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

package env_test

import (
	"os"

	"carbonaut.cloud/pkg/util/env"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Env", func() {
	tmpEnvK1 := "k1"
	tmpEnvV1 := "v1"

	tmpEnvK2 := "k2"
	tmpEnvV2 := "v2"
	Describe("an environment variable is set", func() {
		BeforeEach(func() {
			if err := os.Setenv(tmpEnvK1, tmpEnvV1); err != nil {
				Expect(err).To(BeNil())
			}
		})
		AfterEach(func() {
			if err := os.Setenv(tmpEnvK2, ""); err != nil {
				Expect(err).To(BeNil())
			}
			if err := os.Setenv(tmpEnvK1, ""); err != nil {
				Expect(err).To(BeNil())
			}
		})
		It("should detect an existing environment variable", func() {
			Expect(env.IsSet(tmpEnvK1)).To(BeTrue())
		})
		It("should detect that the environment variable is not set", func() {
			Expect(env.IsSet(tmpEnvK2)).To(BeFalse())
		})
		It("should return the value of the environment variable", func() {
			Expect(env.Default(tmpEnvK2, tmpEnvV2)).To(Equal(tmpEnvV2))
			Expect(env.IsSet(tmpEnvK2)).To(BeTrue())
		})
	})
})
