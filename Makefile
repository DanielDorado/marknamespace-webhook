# Go build
.PHONY: build
build:
	mkdir -p ./bin
	go build -o ./bin/marknamespace

.PHONY: test
test:
	go test

.PHONY: docker-build
docker-build:
	scripts/docker-build.sh

.PHONY: manifest-build
manifest-build:
	scripts/k8s-build.sh

.PHONY: manifest-clean
manifest-clean:
	rm build -rf

.PHONY: k8s-deploy
k8s-deploy:
	oc get project/marknamespace-webhook > /dev/null || ( sleep 4 && oc new-project marknamespace-webhook )
	oc apply -f build

.PHONY:
 k8s-clean:
	( oc get MutatingWebhookConfiguration/marknamespace-webhook && oc delete ValidatingWebhookConfiguration/marknamespace-webhook ) || true
	oc delete project/marknamespace-webhook || true
