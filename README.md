# aws-credentials-mutating-webhook
This Webhook has the target to inject AWS credentials into a pod, getting ARNs details from annotations.
The project has been created with [kubebuilder](https://book.kubebuilder.io) 
It uses Kubernetes and AWS Go libraries to make the webhook working.


## Description
This project was born to fulfill the request to be authorized to do something on AWS even though the request comes from outside AWS itself.
So that [IAM Roles Anywhere](https://docs.aws.amazon.com/sdkref/latest/guide/access-rolesanywhere.html) is one of the possible ways to aim it. 
In order to configure IAM Roles Anywhere, pleas read [this doc](https://docs.aws.amazon.com/rolesanywhere/latest/userguide/getting-started.html).
The Webhook server must start up using a Certificate issued by the same CA used to create the Tust Anchor.
**AWS required configurations**:
- Trust Anchor (CA, self-signed is ok as well)
- Profile:
    - Policy like this:
    ```
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Sid": "ListObjectsInBucket",
                    "Effect": "Allow",
                    "Action": [
                        "s3:ListBucket"
                    ],
                    "Resource": [
                        "arn:aws:s3:::onetag-aws-anywhere"
                    ]
                },
                {
                    "Sid": "AllObjectActions",
                    "Effect": "Allow",
                    "Action": "s3:*Object",
                    "Resource": [
                        "arn:aws:s3:::onetag-aws-anywhere/*"
                    ]
                }
            ]
        }
    ```
- Role like this:
    - Policy (AmazonS3FullAcess):
    ```
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Action": [
                        "s3:*",
                        "s3-object-lambda:*"
                    ],
                    "Resource": "*"
                }
            ]
        }
    ```

    - trust relationship:
    ```
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "rolesanywhere.amazonaws.com"
                    },
                    "Action": [
                        "sts:AssumeRole",
                        "sts:TagSession",
                        "sts:SetSourceIdentity"
                    ]
                }
            ]
        }
    ```

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running Webhook on the cluster
1. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/aws-credentials-mutating-webhook:tag
```

2. Deploy required resources to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/aws-credentials-mutating-webhook:tag
```

3. deploy example container (k8s job):
```sh
make example-deploy
```

### Undeploy Webhook
UnDeploy the webhook resources from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.


More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

