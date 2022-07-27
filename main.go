package main

import (
	acme "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"k8s.io/client-go/rest"
	"simply-dns-webhook/client"
	"sync"
)

func main() {

}

type SimplySolver struct {
	name   string
	client client.SimplyClient
	sync.RWMutex
}

func (e *SimplySolver) Name() string {
	return e.name
}
func (e *SimplySolver) Present(ch *acme.ChallengeRequest) error {
	e.Lock()
	e.client.AddTxtRecord(ch.DNSName, ch.Key)
	e.Unlock()
	return nil
}

func (e *SimplySolver) CleanUp(ch *acme.ChallengeRequest) error {
	e.Lock()
	Id := e.client.GetTxtRecord(ch.DNSName)
	e.client.RemoveTxtRecord(Id)
	e.Unlock()
	return nil
}

func (e *SimplySolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return nil
}
