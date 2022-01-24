package soak

import (
	"context"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// ApplyKustomization builds and applies a Kustomization into the existing
// cluster.
func ApplyKustomization(basePath string, namespace string) error {
	k := krusty.MakeKustomizer(
		krusty.MakeDefaultOptions(),
	)
	fSys := filesys.MakeFsOnDisk()
	m, err := k.Run(fSys, basePath)
	if err != nil {
		return err
	}
	return createKustomizationResources(&m, namespace)
}

func findGVR(gvk *schema.GroupVersionKind, cfg *rest.Config) (*meta.RESTMapping, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
}

func createKustomizationResources(resourceMap *resmap.ResMap, namespace string) error {
	kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return err
	}

	client, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	for _, r := range (*resourceMap).Resources() {
		yaml, err := r.AsYAML()
		if err != nil {
			return err
		}

		gvk := r.GetGvk()

		decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDecoder()
		obj := &unstructured.Unstructured{}
		if err = runtime.DecodeInto(decoder, yaml, obj); err != nil {
			return err
		}

		mapping, err := findGVR(&schema.GroupVersionKind{
			Group:   gvk.Group,
			Version: gvk.Version,
			Kind:    gvk.Kind,
		}, kubeConfig)
		if err != nil {
			return err
		}

		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			dr = client.Resource(mapping.Resource).Namespace(namespace)
		} else {
			dr = client.Resource(mapping.Resource)
		}

		_, err = dr.Create(context.Background(), obj, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
