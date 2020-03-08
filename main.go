package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.com/domac/crddemo/pkg/client/clientset/versioned"
	informers "github.com/domac/crddemo/pkg/client/informers/externalversions"
)

//程序启动参数
var (
	flagSet              = flag.NewFlagSet("crddemo", flag.ExitOnError)
	master               = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	kubeconfig           = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

//设置信号处理
func setupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}

func main() {
	flag.Parse()

	stopCh := setupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	mydemoClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	mydemoInformerFactory := informers.NewSharedInformerFactory(mydemoClient, time.Second*30)

	controller := NewController(kubeClient, mydemoClient,
		mydemoInformerFactory.Crddemo().V1().Mydemos())

	go mydemoInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
