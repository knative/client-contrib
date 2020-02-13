package cmd

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type AdminParams struct {
	KubeCfgPath  string
	ClientConfig clientcmd.ClientConfig
	ClientSet    *kubernetes.Clientset
}
