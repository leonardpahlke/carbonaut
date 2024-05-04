package rnd_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rnd Suite")
}
