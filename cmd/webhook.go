package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podAnnotator struct {
	Client    client.Client
	decoder   *admission.Decoder
	wsOptions WebhookOptions
}

func (a *podAnnotator) Handle(ctx context.Context, req admission.Request) admission.Response {
	outputCreds := new(CredentialOutput)
	pod := corev1.Pod{}
	err := a.decoder.Decode(req, &pod)
	if err != nil {
		ErrorLogger.Panicf("cannot decode pod %s", ctx.Err().Error())
		ErrorLogger.Panicf("error %s", &req.Object)
		return admission.Errored(http.StatusBadRequest, err)
	}

	if shouldMutate(&pod, a) {
		namespace = pod.GetNamespace()
		credsError := getCredentials(outputCreds)
		if credsError != nil {
			return admission.Errored(http.StatusInternalServerError, credsError)
		}

		InfoLogger.Println(*outputCreds)

		secName, err := handleSecret(SecretHandler{
			K8sClient: a.Client,
			Ctx:       ctx,
			Data:      generateSecret(outputCreds),
		})

		if err != nil {
			ErrorLogger.Panic(err)
			return admission.Errored(http.StatusBadRequest, err)
		}

		mPod, err := injectAWSEnvs(&pod, outputCreds, secName)
		if err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}

		marshaledPod, err := json.Marshal(mPod)
		if err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}
		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
	} else {
		InfoLogger.Printf("Intrcepted pod %s.. Nothing to do, AWS annotations not present", pod.Name)
		return admission.Allowed("AWS annotations not present")
	}
}

func shouldMutate(pod *corev1.Pod, pa *podAnnotator) bool {
	annotations := pod.Annotations

	if annotations != nil {

		trustAnchorArnValue, existTAA := annotations[fmt.Sprintf("%s/%s", annotationPrefix, trustAnchorArnAnotationName)]
		profileArnValue, existPA := annotations[fmt.Sprintf("%s/%s", annotationPrefix, profileArnAnotationName)]
		roleArnValue, existRA := annotations[fmt.Sprintf("%s/%s", annotationPrefix, roleArnAnotationName)]
		if existTAA && existPA && existRA {
			InfoLogger.Println("trustAnchor Arn" + trustAnchorArnValue)
			InfoLogger.Println("profile Arn" + profileArnValue)
			InfoLogger.Println("role Arn" + roleArnValue)
			opts = CredentialProcess{
				CertificateId:     fmt.Sprintf("%s/%s", pa.wsOptions.certDir, pa.wsOptions.certFileName),
				PrivateKeyId:      fmt.Sprintf("%s/%s", pa.wsOptions.certDir, pa.wsOptions.pKeyFileName),
				TrustAnchorArnStr: trustAnchorArnValue,
				RoleArnStr:        roleArnValue,
				ProfileArnStr:     profileArnValue,
			}
			return true
		}
	}
	return false
}

func generateSecret(creds *CredentialOutput) []SecretDataValue {
	values := []SecretDataValue{}

	values = append(values, SecretDataValue{
		Name:  awsAccessKeyIdEnvName,
		Value: creds.AccessKeyId,
	})

	values = append(values, SecretDataValue{
		Name:  awsSecretAccessKeyEnvName,
		Value: creds.SecretAccessKey,
	})

	values = append(values, SecretDataValue{
		Name:  awsSessionTokenEnvName,
		Value: creds.SessionToken,
	})

	return values

}
