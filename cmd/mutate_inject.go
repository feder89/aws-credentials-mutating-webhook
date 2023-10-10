package main

import (
	corev1 "k8s.io/api/core/v1"
)

var (
	awsAccessKeyIdEnvName     string = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKeyEnvName string = "AWS_SECRET_ACCESS_KEY"
	awsSessionTokenEnvName    string = "AWS_SESSION_TOKEN"
)

func injectAWSEnvs(pod *corev1.Pod, credsOutput *CredentialOutput, secretName string) (*corev1.Pod, error) {
	mutatingPod := pod.DeepCopy()

	injectEnvVar(mutatingPod, envVarObjectRef(awsAccessKeyIdEnvName))
	injectEnvVar(mutatingPod, envVarObjectRef(awsSecretAccessKeyEnvName))
	injectEnvVar(mutatingPod, envVarObjectRef(awsSessionTokenEnvName))

	return mutatingPod, nil
}

func injectEnvVar(pod *corev1.Pod, envVar corev1.EnvVar) {
	for i, container := range pod.Spec.Containers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.Containers[i].Env = append(container.Env, envVar)
		}
	}
}

func HasEnvVar(container corev1.Container, checkEnvVar corev1.EnvVar) bool {
	for _, envVar := range container.Env {
		if envVar.Name == checkEnvVar.Name {
			return true
		}
	}
	return false
}

func envVarObjectRef(value string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: value,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key:                  value,
				LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
			},
		},
	}
}
