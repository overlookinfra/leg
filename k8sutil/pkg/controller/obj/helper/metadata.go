package helper

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/norm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Annotate(target metav1.Object, name, value string) bool {
	annotations := target.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	} else if candidate, ok := annotations[name]; ok && candidate == value {
		return false
	}

	annotations[name] = value
	target.SetAnnotations(annotations)
	return true
}

func Label(target metav1.Object, name, value string) bool {
	labels := target.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	} else if candidate, ok := labels[name]; ok && candidate == value {
		return false
	}

	labels[name] = value
	target.SetLabels(labels)
	return true
}

func CopyLabelsAndAnnotations(target, src metav1.Object) {
	for name, value := range src.GetAnnotations() {
		Annotate(target, name, value)
	}

	for name, value := range src.GetLabels() {
		Label(target, name, value)
	}
}

func SuffixObjectKey(key client.ObjectKey, suffix string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: key.Namespace,
		Name:      norm.MetaNameSuffixed(key.Name, "-"+suffix),
	}
}
