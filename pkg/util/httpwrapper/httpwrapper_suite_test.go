package httpwrapper_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHttpwrapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httpwrapper Suite")
}
