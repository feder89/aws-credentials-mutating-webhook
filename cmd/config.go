package main

import "os"

func config() {

	prefix, ok := os.LookupEnv("AWS_ANNOTATION_PREFIX")
	if ok {
		annotationPrefix = prefix
	}

	anchorAnnotationName, ok := os.LookupEnv("AWS_TRUST_ANCHOR_ANNOTATION_NAME")
	if ok {
		trustAnchorArnAnotationName = anchorAnnotationName
	}

	profileAnnotationName, ok := os.LookupEnv("AWS_PROFILE_ANNOTATION_NAME")
	if ok {
		profileArnAnotationName = profileAnnotationName
	}

	roleAnnotationName, ok := os.LookupEnv("AWS_ROLE_ANNOTATION_NAME")
	if ok {
		roleArnAnotationName = roleAnnotationName
	}

}
