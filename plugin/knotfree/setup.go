package knotfree

import (
	"fmt"
	"os"

	"github.com/awootton/knotfreeiot/iot"
	"github.com/awootton/knotfreeiot/tokens"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

var count = 0

// setup is the function that gets called when the config parser see the token "knotfree". Setup is responsible
// for parsing any extra options the knotfree plugin may have. The first token this function sees is "knotfree".
func setup(c *caddy.Controller) error {

	log.Info("knotfree setup", count)

	c.Next() // Ignore "knotfree" and give us the next token.
	if c.NextArg() {
		// If there was another token, return an error, because we don't have any configuration.
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error("knotfree", c.ArgErr())
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {

		if next != nil {
			log.Info("knotfree AddPlugin next", next.Name(), count)
		} else {
			log.Info("knotfree AddPlugin next is nil", count)
		}

		kf := Knotfree{Next: next, instanceCount: count}
		count++

		var err error
		kf.sc = nil
		_ = kf.sc
		// init the service contact
		host := os.Getenv("TARGET_CLUSTER") // "localhost:8384"
		if host == "" {
			host = "knotfree.io:8384"
			fmt.Println("No TARGET_CLUSTER environment variable set, using default host", host)

		}
		token := os.Getenv("KNOTFREE_TOKEN")
		if token == "" {
			fmt.Println("No KNOTFREE_TOKEN environment variable set, using default token")
			token, _ = tokens.GetImpromptuGiantTokenLocal("", "")
		}
		kf.sc, err = iot.StartNewServiceContactTcp(host, token)
		if err != nil {
			log.Error("knotfree setup failed to start service contact", err)
			return nil
		}
		return kf
	})

	log.Info("knotfree setup OK")

	// All OK, return a nil error.
	return nil
}
