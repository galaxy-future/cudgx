package clients_test

import (
	"github.com/galaxy-future/cudgx/internal/clients"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestXclient(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Xclient Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	clients.InitializeBridgxClient("http://bridgx-api.internal.galaxy-future.org")
	clients.InitializeSchedulxClient("http://10.16.23.96:9090")
})
