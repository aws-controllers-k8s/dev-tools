package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensiontypesv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = apiextensionv1.AddToScheme(scheme)
	return scheme
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := apiextension.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	err = ListCRDs(clientset, "examplecrd.yaml")
	if err != nil {
		panic(err)
	}

}

func CreateCRD(clientset apiextension.Interface, file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	s := getScheme()
	decode := serializer.NewCodecFactory(s).UniversalDeserializer().Decode
	obj, _, err := decode(b, nil, nil)
	if err != nil {
		return err
	}

	fmt.Printf("%++v\n\n", obj.GetObjectKind())

	crd := obj.(*apiextensionv1.CustomResourceDefinition) // This fails

	/* 	switch o := obj.(type) {
	   	case *rbacv1.ClusterRoleBindingList:
	   		fmt.Println("correct found") // Never happens
	   	default:
	   		fmt.Println("default case")
	   		_ = o
	   	} */

	fmt.Println(crd.Name)
	_, err = clientset.ApiextensionsV1().CustomResourceDefinitions().Create(context.Background(), crd, metav1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func ListCRDs(clientset apiextension.Interface, f string) error {
	l, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil
	}

	for _, i := range l.Items {
		fmt.Println(i)
	}
	return err
}

func NewClient(kubeconfigPath string) apiextensiontypesv1.CustomResourceDefinitionInterface {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err)
	}

	clientset, err := apiextension.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset.ApiextensionsV1().CustomResourceDefinitions()
}
