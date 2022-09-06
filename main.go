package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"simply-dns-webhook/client"
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

type SimplyDNSProviderConfig struct {
	// Change the two fields below according to the format of the configuration
	// to be decoded.
	// These fields will be set by users in the
	// `issuer.spec.acme.dns01.providers.webhook.config` field.
	SecretRef string `json:"secretName"`
}

type SimplyDnsSolver struct {
	client     client.SimplyClient
	kubeClient *kubernetes.Clientset
}

func (e *SimplyDnsSolver) Name() string {
	return "simply-dns-solver"
}
func (e *SimplyDnsSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	cred, err2 := loadCredFromSecret(ch, e)
	if err2 != nil {
		return err2
	}
	e.client.AddTxtRecord(ch.DNSName, ch.Key, cred)
	return nil
}

func (e *SimplyDnsSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	cred, err2 := loadCredFromSecret(ch, e)
	if err2 != nil {
		return err2
	}
	Id := e.client.GetTxtRecord(ch.DNSName, cred)
	e.client.RemoveTxtRecord(Id, cred)
	return nil
}

func (e *SimplyDnsSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}

	e.kubeClient = cl

	return nil
}

func loadConfig(cfgJSON *extapi.JSON) (SimplyDNSProviderConfig, error) {
	cfg := SimplyDNSProviderConfig{}

	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}
	return cfg, nil
}

func stringFromSecretData(secretData *map[string][]byte, key string) (string, error) {
	data, ok := (*secretData)[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret data", key)
	}
	return string(data), nil
}

func loadCredFromSecret(ch *v1alpha1.ChallengeRequest, e *SimplyDnsSolver) (client.Credentials, error) {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return client.Credentials{}, err
	}
	secretName := cfg.SecretRef
	sec, err := e.kubeClient.CoreV1().Secrets(ch.ResourceNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return client.Credentials{}, err
	}

	accountName, err := stringFromSecretData(&sec.Data, "account-name")
	apiKey, err := stringFromSecretData(&sec.Data, "api-key")
	if err != nil {
		return client.Credentials{}, err
	}

	cred := client.Credentials{
		AccountName: accountName, ApiKey: apiKey,
	}
	return cred, nil
}
