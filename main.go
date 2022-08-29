package main

import (
	acme "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"k8s.io/client-go/rest"
	"os"
	"simply-dns-webhook/client"
	"sync"
)

var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	// This will register our custom DNS provider with the webhook serving
	// library, making it available as an API under the provided GroupName.
	// You can register multiple DNS provider implementations with a single
	// webhook, where the Name() method will be used to disambiguate between
	// the different implementations.
	cmd.RunWebhookServer(GroupName,
		&SimplyDnsSolver{},
	)
}

type SimplyDnsSolver struct {
	Username string `json:"username"` //Simply user/account name
	Password string `json:"password"` //Simply api key for corresponding account
	client   client.SimplyClient
	sync.RWMutex
}

func (e *SimplyDnsSolver) Name() string {
	return "simply-dns-solver"
}
func (e *SimplyDnsSolver) Present(ch *acme.ChallengeRequest) error {
	e.Lock()
	e.client.AddTxtRecord(ch.DNSName, ch.Key)
	e.Unlock()
	return nil
}

func (e *SimplyDnsSolver) CleanUp(ch *acme.ChallengeRequest) error {
	e.Lock()
	Id := e.client.GetTxtRecord(ch.DNSName)
	e.client.RemoveTxtRecord(Id)
	e.Unlock()
	return nil
}

func (e *SimplyDnsSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return nil
}
