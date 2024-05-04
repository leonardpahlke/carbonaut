package rnd_test

import (
	"carbonaut.dev/pkg/util/compareutils"
	"carbonaut.dev/pkg/util/rnd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rnd", func() {
	for i := 0; i < 5; i++ {
		Describe("neg tests", func() {
			Context("invalid range", func() {
				It("should fail with a null range", func() {
					Expect(rnd.GetNumber(0, 0)).To(Equal(0))
				})
			})
			Context("with min greater than max", func() {
				It("should return -1", func() {
					Expect(rnd.GetNumber(1, 0)).To(Equal(-1))
				})
			})
			Context("with max '1' greater than min '0'", func() {
				It("should return '0' or '1'", func() {
					Expect(compareutils.CheckListContains([]int{0, 1}, rnd.GetNumber(0, 1))).To(BeTrue())
				})
			})
			Context("with max '3' greater than min '1'", func() {
				It("should return '1', '2' or '3'", func() {
					Expect(compareutils.CheckListContains([]int{1, 2, 3}, rnd.GetNumber(1, 3))).To(BeTrue())
				})
			})
		})

		Describe("pos tests", func() {
			Context("with an empty list", func() {
				It("should return an empty list", func() {
					Expect(rnd.GetRandomListSubset([]int{})).To(Equal([]int{}))
				})
			})
			Context("With a list with one element", func() {
				It("should return the exact same list", func() {
					Expect(rnd.GetRandomListSubset([]int{1})).To(Equal([]int{1}))
				})
			})
		})
	}
})
