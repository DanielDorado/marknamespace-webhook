# Marknamespace Webhook

Marknamespace Webhook is a Kubernetes Mutating Webhook that labels and annotates namespaces
according to its name and a set of rules defined in a configuration file.

Useful in Openshift environments where the developers only have `oc new-project`
permissions, but the namespace needs to be labeled according to any convention.

# Configuration

The configuration is in the file `manifests/marknamespace-conf.yml`.

Here there are three parts:

1. `server`: Application configuration.
2. `labels`: Labels Creation rules.
3. `annotations`: Annotations Creation rules.

The rules to create labels and annotations are a list of regular exprexion cases,
`caseNamespace`. Each `caseNamespace` is evaluated with the namespace name until there
is a match, then, the evaluation stops, and the labels or annotations, in the `inject`
section, are patched to the new namespace. The name and value in the `inject` section are
go templates where the values are a slice with the `caseNamespace` subexpressions.

## Example

In this example, there are rules to create labels for a Kubernetes cluster with two namespace
name styles:

- `<region>-<area>-<team>-<environment>` for previous environments
- `<region>-<area>-<team>` for production environment

``` yaml
server:
  port: 8443
  TLS:
    certFile: /data/certs/tls.crt
    keyFile: /data/certs/tls.key

labels:
  - caseNamespace: "^([^-]+)-([^-]+)-([^-]+)-([^-]+)$"
    inject:
    - name: region
      value: "reg-{{index . 0}}"
    - name: area
      value: "{{index . 1}}"
    - name: team
      value: "{{index . 2}}"
    - name: environment
      value: "{{index . 3}}"
  - caseNamespace: "^([^-]+)-([^-]+)-([^-]+)$"
    inject:
    - name: region
      value: "reg-{{index . 0}}"
    - name: area
      value: "{{index . 1}}"
    - name: team
      value: "{{index . 2}}"
    - name: environment
      value: "production"

annotations:
  - caseNamespace: ".*test.*"
    inject:
     - name: "purpose"
       value: "internal-test"
```

If a namespace `europe-customer-sales` is created, these labels will be added:

- `region: reg-europe`
- `area: customer`
- `team: sales`
- `environment: production`

# Getting Started

The manifests, build scripts, and deployment scripts work well with
[RedHat CRC](https://crc.dev/crc/).
The CRC internal registry is used to upload the webhook image. A Makefile is provided
to do all the required operations:

Test and build the go server:

``` sh
make test
make build
```

Build the docker image and push it to the CRC registry:

``` sh
make docker-build
```

Build the Kubernetes manifests, creating the certificates if they do not exist, and
clean them:

``` sh
make manifest-build
make manifest-clean
```

Deploy to Kubernetes and remove the deployed objects:

``` sh
make k8s-deploy
make k8s-clean
```


# References

- [CRC - Getting Started Guide](https://crc.dev/crc/)
- [Webhook Example](https://github.com/kubernetes/kubernetes/tree/release-1.24/test/images/agnhost/webhook)
