package crd

import (
	"errors"
	"io/ioutil"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensiontypesv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	scheme                          *runtime.Scheme
	ErrorObjectNotCRD               = errors.New("object is not a custom resource definition")
	ErrorUnsupportedCRDGroupVersion = errors.New("unsupported CRD group version")
)

func init() {
	scheme = runtime.NewScheme()
	err := apiextensionv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
}

func GetCRDFromFile(filePath string) (*apiextensionv1.CustomResourceDefinition, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	decode := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode
	obj, _, err := decode(b, nil, nil)
	if err != nil {
		return nil, err
	}

	switch o := obj.(type) {
	case *apiextensionv1.CustomResourceDefinition:
		return o, nil
	case *apiextensionv1beta1.CustomResourceDefinition:
		return nil, ErrorUnsupportedCRDGroupVersion
	default:
		return nil, ErrorObjectNotCRD
	}
}

func NewClient(kubeconfig string) (apiextensiontypesv1.CustomResourceDefinitionInterface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := apiextension.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset.ApiextensionsV1().CustomResourceDefinitions(), nil
}
