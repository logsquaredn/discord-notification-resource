package resource_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiscordNotificationResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DiscordNotificationResource Suite")
}
