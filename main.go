package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/runnerm/simply-dns-webhook/client"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
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
	SecretRef   string `json:"secretName"`
	AccountName string `json:"accountName"`
	ApiKey      string `json:"apiKey"`
}

type SimplyDnsSolver struct {
	client     client.SimplyClient
	kubeClient *kubernetes.Clientset
}

func (e *SimplyDnsSolver) Name() string {
	return "simply-dns-solver"
}
func (e *SimplyDnsSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	fmt.Println("Challenge being presented for: ", ch.ResolvedFQDN)
	cred, err := loadCredentials(ch, e)
	if err != nil {
		_ = fmt.Errorf("load credentials failed(check secret configuration): %v", err)
		return err
	}

	id, txtData, fetchErr := e.client.GetTxtRecord(ch.ResolvedFQDN, cred)
	if fetchErr == nil && id != 0 && txtData != ch.Key {
		_, err := e.client.UpdateTXTRecord(id, ch.ResolvedFQDN, ch.Key, cred)
		if err != nil {
			_ = fmt.Errorf("presenting challenge failed: %v", err)
			return err
		}
		fmt.Println("Challenge have been created with record id: ", id)
		return nil
	} else if fetchErr == nil && id != 0 && txtData == ch.Key {
		fmt.Println("Challenge have been created with record id: ", id)
		return nil
	} else {
		id, err = e.client.AddTxtRecord(ch.ResolvedFQDN, ch.Key, cred)
		if err != nil {
			_ = fmt.Errorf("presenting challenge failed: %v", err)
			return err
		} else {
			fmt.Println("Challenge have been created with record id: ", id)
		}
		return nil
	}
}

func (e *SimplyDnsSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	fmt.Println("Challenge being cleaned up...")
	cred, err := loadCredentials(ch, e)
	if err != nil {
		_ = fmt.Errorf("load credentials failed: %v", err)
		return err
	}
	id, err := e.client.GetExactTxtRecord(ch.Key, ch.ResolvedFQDN, cred)
	if err != nil {
		_ = fmt.Errorf("error on fetching record: %v", err)
		return err
	}
	fmt.Println("Record(", id, ") fetched for cleanup.")
	res := e.client.RemoveTxtRecord(id, ch.DNSName, cred)
	if res == true {
		fmt.Println("Record(", id, ") have been cleaned up.")
		return nil
	} else {
		_ = fmt.Errorf("record(%d) have no tbeen cleaned up", id)
		return err
	}

}

func (e *SimplyDnsSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	fmt.Println("Initializing...")
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		_ = fmt.Errorf("init failed with error: %v", err)
		return err
	}
	e.kubeClient = cl
	return nil
}

func loadConfig(cfgJSON *extapi.JSON) (SimplyDNSProviderConfig, error) {
	fmt.Println("Loading config...")
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

func loadCredentials(ch *v1alpha1.ChallengeRequest, e *SimplyDnsSolver) (client.Credentials, error) {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		_ = fmt.Errorf("error on reading config: %v", err)
		return client.Credentials{}, err
	}
	if cfg.AccountName != "" && cfg.ApiKey != "" {
		fmt.Println("Loading API credentials from config.")
		cred := client.Credentials{
			AccountName: cfg.AccountName, ApiKey: cfg.ApiKey,
		}
		return cred, nil
	} else {
		secretName := cfg.SecretRef
		fmt.Println("Loading API credentials, secret reference:")
		fmt.Println(secretName)
		sec, err := e.kubeClient.CoreV1().Secrets(ch.ResourceNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			_ = fmt.Errorf("error on loading secret from kubernetes api: %v", err)
			return client.Credentials{}, err
		}

		accountName, err := stringFromSecretData(&sec.Data, "account-name")
		apiKey, err := stringFromSecretData(&sec.Data, "api-key")
		if err != nil {
			_ = fmt.Errorf("error on reading secret: %v", err)
			return client.Credentials{}, err
		}

		cred := client.Credentials{
			AccountName: accountName, ApiKey: apiKey,
		}
		return cred, nil
	}
}
