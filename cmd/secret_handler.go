package main

import (
	"fmt"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	secretName string
)

func handleSecret(sh SecretHandler) (string, error) {
	/*sec := &corev1.Secret{}
	object := types.NamespacedName{Namespace: namespace, Name: awsAnywhereSecretName}
	exists, err := checkSecretExist(sh, object, sec)
	if err != nil {
		ErrorLogger.Panicln(err)
		return err
	}

	if  exists {
		err := deleteSecret(sh, sec)
		if err != nil {
			ErrorLogger.Panicln(err)
			return err
		}
	}*/

	err := createSecret(sh)

	if err != nil {
		ErrorLogger.Panicln(err)
		return "", err
	}

	return secretName, nil

}

func checkSecretExist(sh SecretHandler, obj types.NamespacedName, secret *corev1.Secret) (bool, error) {
	err := sh.K8sClient.Get(sh.Ctx, obj, secret)
	if err != nil {
		ErrorLogger.Panicln(err)
		return false, err
	}

	if secret == nil {
		return false, nil
	}

	return true, nil
}

func createSecret(sh SecretHandler) error {
	secretData := make(map[string][]byte)

	secretName = fmt.Sprintf("%s-%s", awsAnywhereSecretName, uuid.NewString()[:6])

	for _, e := range sh.Data {
		secretData[e.Name] = []byte(e.Value)
		InfoLogger.Printf("name is %s value is %s", e.Name, secretData[e.Name])
	}
	InfoLogger.Printf("data is %v", secretData)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: secretData,
	}

	InfoLogger.Printf("data is %v", secret)

	err := sh.K8sClient.Create(sh.Ctx, secret)

	if err != nil {
		ErrorLogger.Panicln(err)
		return err
	}

	return nil
}

func deleteSecret(sh SecretHandler, obj client.Object) error {
	err := sh.K8sClient.Delete(sh.Ctx, obj)
	if err != nil {
		ErrorLogger.Panicln(err)
		return err
	}
	return nil
}
