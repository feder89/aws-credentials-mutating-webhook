package main

import (
	helper "github.com/aws/rolesanywhere-credential-helper/aws_signing_helper"
)

var (
	debug              bool = true
	credentialsOptions helper.CredentialsOpts
)

func getCredentials(credsOutput *CredentialOutput) error {
	err := PopulateCredentialsOptions(&opts)
	if err != nil {
		ErrorLogger.Panicln(err)
		return err
	}

	helper.Debug = debug

	signer, signingAlgorithm, err := helper.GetSigner(&credentialsOptions)
	if err != nil {
		ErrorLogger.Panicln(err)
		return err
	}
	defer signer.Close()
	credentialProcessOutput, err := helper.GenerateCredentials(&credentialsOptions, signer, signingAlgorithm)
	if err != nil {
		ErrorLogger.Panicln(err)
		return err
	}

	credsOutput.AccessKeyId = credentialProcessOutput.AccessKeyId
	credsOutput.SecretAccessKey = credentialProcessOutput.SecretAccessKey
	credsOutput.SessionToken = credentialProcessOutput.SessionToken

	return nil
}

func PopulateCredentialsOptions(credentialProcess *CredentialProcess) error {

	credentialsOptions = helper.CredentialsOpts{
		PrivateKeyId:      credentialProcess.PrivateKeyId,
		CertificateId:     credentialProcess.CertificateId,
		RoleArn:           credentialProcess.RoleArnStr,
		ProfileArnStr:     credentialProcess.ProfileArnStr,
		TrustAnchorArnStr: credentialProcess.TrustAnchorArnStr,
		Debug:             debug,
		SessionDuration:   3600,
	}

	return nil
}
