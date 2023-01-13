package main

import (
	"github.com/cert-manager/cert-manager/test/acme/dns"
	"os"
	"testing"
	"time"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.
	//

	// Uncomment the below fixture when implementing your custom DNS provider
	propLimit, _ := time.ParseDuration("10m")
	pollInterval, _ := time.ParseDuration("30s")
	fixture := dns.NewFixture(&SimplyDnsSolver{},
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("testdata/simply-dns-webhook"),
		dns.SetPropagationLimit(propLimit),
		dns.SetPollInterval(pollInterval),
		//dns.SetStrict(true), // concurrent challenges on same domain is obsolete
		//dns.SetBinariesPath("_test/kubebuilder/bin"),
	)
	//fixture := dns.NewFixture(&SimplyDnsSolver{},
	//	dns.SetResolvedZone("example.com."),
	//	dns.SetManifestPath("testdata/simply-dns-webhook"),
	//	dns.SetDNSServer("127.0.0.1:59351"),
	//	dns.SetUseAuthoritative(false),
	//)
	//need to uncomment and  RunConformance delete runBasic and runExtended once https://github.com/cert-manager/cert-manager/pull/4835 is merged
	fixture.RunConformance(t)
	//fixture.RunBasic(t)
	//fixture.RunExtended(t)
}
