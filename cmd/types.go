package main

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CredentialProcess struct {
	RoleArnStr        string
	ProfileArnStr     string
	TrustAnchorArnStr string
	CertificateId     string
	PrivateKeyId      string
}

type CredentialOutput struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

type SecretDataValue struct {
	Name  string
	Value string
}

type SecretHandler struct {
	K8sClient client.Client
	Ctx       context.Context
	Data      []SecretDataValue
}
