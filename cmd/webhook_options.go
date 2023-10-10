package main

import (
	"os"
	"path/filepath"
)

type WebhookOptions struct {
	certDir      string
	certFileName string
	pKeyFileName string
}

func NewWebhookOptions(certificateDir string, certificateFileName string, privateKeyFileName string) WebhookOptions {
	wo := WebhookOptions{
		certDir:      certificateDir,
		certFileName: certificateFileName,
		pKeyFileName: privateKeyFileName,
	}
	setDefaults(&wo)
	return wo
}

func setDefaults(wo *WebhookOptions) {
	if len(wo.certDir) == 0 {
		wo.certDir = filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs")
	}

	if len(wo.certFileName) == 0 {
		wo.certFileName = "tls.crt"
	}

	if len(wo.pKeyFileName) == 0 {
		wo.pKeyFileName = "tls.key"
	}
}
