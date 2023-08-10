package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/runnerm/simply-com-client"
	log "github.com/sirupsen/logrus"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"time"
)

var GroupName = os.Getenv("GROUP_NAME")
var LogLevel = os.Getenv("LOG_LEVEL")

func main() {
	time.Sleep(10 * time.Second)
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
	client     simplyComClient.SimplyClient
	kubeClient *kubernetes.Clientset
}

func (e *SimplyDnsSolver) Name() string {
	return "simply-dns-solver"
}
func (e *SimplyDnsSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	log.Info("Presenting challenge for: ", ch.ResolvedFQDN)
	err := loadCredentials(ch, e)
	if err != nil {
		log.Errorf("load credentials failed(check secret configuration): %v", err)
		return err
	}

	id, txtData, fetchErr := e.client.GetRecord(ch.ResolvedFQDN, ch.Key, "TXT")
	if fetchErr == nil && id != 0 && txtData != ch.Key {
		_, err := e.client.UpdateRecord(id, ch.ResolvedFQDN, ch.Key, "TXT")
		if err != nil {
			log.Errorf("presenting challenge failed: %v", err)
			return err
		}
		log.Debug("Challenge have been created with record id: ", id)
		return nil
	} else if fetchErr == nil && id != 0 && txtData == ch.Key {
		log.Debug("Challenge have been created with record id: ", id)
		return nil
	} else {
		id, err = e.client.AddRecord(ch.ResolvedFQDN, ch.Key, "TXT")
		if err != nil {
			log.Errorf("presenting challenge failed: %v", err)
			return err
		} else {
			log.Debug("Challenge have been created with record id: ", id)
		}
		return nil
	}
}

func (e *SimplyDnsSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	log.Info("Cleaning up challenge for: ", ch.ResolvedFQDN)
	err := loadCredentials(ch, e)
	if err != nil {
		log.Errorf("load credentials failed(check secret configuration): %v", err)
		return err
	}
	id, _, err := e.client.GetRecord(ch.ResolvedFQDN, ch.Key, "TXT")
	if err != nil {
		log.Errorf("error on fetching record: %v", err)
		return err
	}
	log.Debug("Record(", id, ") fetched for cleanup.")
	res := e.client.RemoveRecord(id, ch.ResolvedFQDN)
	if res == true {
		log.Debug("Record(", id, ") have been cleaned up.")
		return nil
	} else {
		log.Errorf("record(%d) have not been cleaned up", id)
		return err
	}

}

func (e *SimplyDnsSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	if LogLevel == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.Info("Initializing Simply dns solver")
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		log.Errorf("init failed with error: %v", err)
		return err
	}
	e.kubeClient = cl
	return nil
}

func loadConfig(cfgJSON *extapi.JSON) (SimplyDNSProviderConfig, error) {
	log.Debug("Loading config...")
	cfg := SimplyDNSProviderConfig{}

	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		log.Errorf("error decoding solver config: %v", err)
		return cfg, errors.New("error decoding solver config")
	}
	log.Debug("Config loaded successfully.")
	return cfg, nil
}

func stringFromSecretData(secretData *map[string][]byte, key string) (string, error) {
	data, ok := (*secretData)[key]
	if !ok {
		log.Errorf("key %q not found in secret data", key)
		return "", errors.New("key not found in secret data")
	}
	return string(data), nil
}

func loadCredentials(ch *v1alpha1.ChallengeRequest, e *SimplyDnsSolver) error {
	if e.client.Credentials.AccountName != "" && e.client.Credentials.ApiKey != "" {
		return nil
	}

	cfg, err := loadConfig(ch.Config)
	if err != nil {
		log.Errorf("error on loading config: %v", err)
		return err
	}
	if cfg.AccountName != "" && cfg.ApiKey != "" {
		log.Debug("Loading API credentials from config.")
		e.client = simplyComClient.CreateSimplyClient(cfg.AccountName, cfg.ApiKey)
		return nil
	} else {
		secretName := cfg.SecretRef
		log.Debug("Loading API credentials from secret: ", secretName)
		sec, err := e.kubeClient.CoreV1().Secrets(ch.ResourceNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			log.Errorf("error on loading secret from kubernetes api: %v", err)
			return err
		}

		accountName, err := stringFromSecretData(&sec.Data, "account-name")
		apiKey, err := stringFromSecretData(&sec.Data, "api-key")
		if err != nil {
			log.Errorf("error on reading secret: %v", err)
			return err
		}

		e.client = simplyComClient.CreateSimplyClient(accountName, apiKey)
		return nil
	}
}
