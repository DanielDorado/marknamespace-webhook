package main

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	// "k8s.io/kubernetes/pkg/apis/admission"
)

const (
	patchLabels      = `{"op": "add", "path": "/metadata/labels/%s", "value": "%s"}`
	patchAnnotations = `{"op": "add", "path": "/metadata/annotations/%s", "value": "%s"}`
)

func (c *Config) mutateNamespace(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.Info("mutateNamespace")
	// check inputs
	namespaceResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
	if ar.Request.Resource != namespaceResource { // check resource
		err := fmt.Errorf("Wrong Resource received in mutateNamespace. Expected: %+v. Get: %+v.",
			namespaceResource, ar.Request.Resource)
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	namespaceKind := metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}
	if ar.Request.Kind != namespaceKind { // check kind
		err := fmt.Errorf("Wrong Kind received in mutateNamespace. Expected: %+v. Get: %+v.",
			namespaceKind, ar.Request.Kind)
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	// deserialize
	namespace := corev1.Namespace{}
	decoder := codecs.UniversalDeserializer()
	_, _, err := decoder.Decode(ar.Request.Object.Raw, nil, &namespace)
	if err != nil {
		err = fmt.Errorf("Decoding in mutateNamespace: %w", err)
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	// get current labels and annotations
	labels := namespace.ObjectMeta.Labels
	annotations := namespace.ObjectMeta.Annotations
	// calculate labels and annotations
	newLabels := getLabelsFromName(namespace.Name, c)
	newAnnotations := getAnnotationsFromName(namespace.Name, c)
	// remove repeated
	labelsToInject := removeRepetedInMap(newLabels, labels)
	annotationsToInject := removeRepetedInMap(newAnnotations, annotations)
	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
	if len(labelsToInject)+len(annotationsToInject) > 0 { // patch it
		patches := []string{}
		for k, v := range labelsToInject {
			patches = append(patches, fmt.Sprintf(patchLabels, k, v))
		}
		for k, v := range annotationsToInject {
			patches = append(patches, fmt.Sprintf(patchAnnotations, k, v))
		}
		reviewResponse.Patch = []byte("[\n" + strings.Join(patches, ",\n") + "\n]")
		patchType := v1.PatchTypeJSONPatch
		reviewResponse.PatchType = &patchType
		klog.Infof("mutateNamespace applying patch in %s:\n%s", namespace.Name, reviewResponse.Patch)
	} else {
		klog.Infof("No labels or annotations to patch when creating namespace: %s", namespace.Name)
	}
	return &reviewResponse
}

func removeRepetedInMap(newMap, originalMap map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range newMap {
		vOld, ok := originalMap[k]
		if !(ok && v == vOld) {
			result[k] = v
		}
	}
	return result
}

func getLabelsFromName(name string, config *Config) map[string]string {
	return processNamespace(name, config.Labels)
}

func getAnnotationsFromName(name string, config *Config) map[string]string {
	return processNamespace(name, config.Annotations)
}
