package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubeconfigPath = "run-k8s-job-kubeconfig"
)

var (
	errNoAuth = errors.New("you must provide either 'kubeconfig-file' or both 'cluster-url' and 'cluster-token'")
)

type ActionInput struct {
	kubeconfigFile string
	jobFile        string
	jobName        string
	namespace      string
	clusterURL     string
	clusterToken   string
	caFile         string
}

func BuildK8sConfig(input ActionInput) (*rest.Config, error) {
	if len(input.jobFile) == 0 {
		return nil, errors.New("'jobfile' is a required input but was empty")
	}

	if len(input.kubeconfigFile) == 0 {
		return buildConfigWithSecondaryAuth(input)
	}

	data, err := base64.StdEncoding.DecodeString(input.kubeconfigFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode kubeconfig file")
	}

	err = ioutil.WriteFile(kubeconfigPath, data, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode kubeconfig file")
	}
	defer os.Remove(kubeconfigPath)

	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func buildConfigWithSecondaryAuth(input ActionInput) (*rest.Config, error) {
	if len(input.clusterURL) == 0 || len(input.clusterToken) == 0 {
		return nil, errors.Wrap(errNoAuth, "missing input for cluster authentication")
	}

	config, err := clientcmd.BuildConfigFromFlags(input.clusterURL, "")
	if err != nil {
		return nil, err
	}

	config.BearerToken = input.clusterToken
	config.CAFile = input.caFile

	return config, nil
}
