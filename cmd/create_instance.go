/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"context"
	"errors"
	"log"

	"github.com/sitewhere/swctl/internal"
	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// createInstanceCmd represents the instance command
var (
	namespace         = ""
	createInstanceCmd = &cobra.Command{
		Use:   "instance",
		Short: "Create SiteWhere Instance",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires one argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			if namespace == "" {
				namespace = name
			}

			instance := alpha3.SiteWhereInstance{
				Name:      name,
				Namespace: namespace}

			createSiteWhereInstance(&instance)
		},
	}
)

func init() {
	createInstanceCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of the instance.")
	createCmd.AddCommand(createInstanceCmd)
}

func createSiteWhereInstance(instance *alpha3.SiteWhereInstance) {
	var err error

	config, err := internal.GetKubeConfigFromKubeconfig()
	if err != nil {
		log.Printf("Error getting Kubernetes Config: %v", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error getting Kubernetes Client: %v", err)
		return
	}

	createNamespaceIfNotExist(instance.Namespace, clientset)
}

func createNamespaceIfNotExist(namespace string, clientset *kubernetes.Clientset) (*v1.Namespace, error) {
	var err error
	var ns *v1.Namespace

	ns, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	if ns != nil {
		return ns, nil
	}

	ns = &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"name": namespace,
			},
		},
	}

	result, err := clientset.CoreV1().Namespaces().Create(context.TODO(),
		ns,
		metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return result, err
}
