package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	"golang.org/x/net/context"

	helper "github.com/scaler/helpers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	prometheusURL = "http://127.0.0.1:9090"
	cpuMetric     = "process_cpu_seconds_total"
	memoryMetric  = "mongodb_ss_mem_resident"
	ioReadMetric  = "mongodb_sys_disks_sda_read_time_ms"
	ioWriteMetric = "mongodb_sys_disks_sda_write_time_ms"
	mongodbTarget = "mongo-exporter-prometheus-mongodb-exporter"
	hpaName       = "mongodb-hpa"
	namespace     = "default"
	waitInterval  = time.Minute
)

func main() {
	// Create a Prometheus API client.
	cfg := api.Config{Address: prometheusURL}
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	api := v1.NewAPI(client)

	// Create a Kubernetes client.
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", "/home/sadath/.kube/config")
		if err != nil {
			panic(err)
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	for {

		cpuValue, err := helper.QueryPrometheus(api, cpuMetric)
		if err != nil {
			panic(err)
		}

		memoryValue, err := helper.QueryPrometheus(api, memoryMetric)
		if err != nil {
			panic(err)
		}
		

		// Query Prometheus for the average I/O read operations.
		ioReadValue, err := helper.QueryPrometheus(api, ioReadMetric)
		if err != nil {
			panic(err)
		}

		// // Query Prometheus for the average I/O write operations.
		ioWriteValue, err := helper.QueryPrometheus(api, ioWriteMetric)
		if err != nil {
			panic(err)
		}

		// Calculate the average CPU, memory, and I/O usage.
		
		// cpuUsage := cpuValue / float64(time.Minute.Seconds())
		// memoryUsage := memoryValue / (1024 * 1024)
		// ioReadOps := ioReadValue / float64(time.Minute.Seconds())
		// ioWriteOps := ioWriteValue / float64(time.Minute.Seconds())

		fmt.Printf(" CPU usage: %.2f\n", cpuValue)
		fmt.Printf(" memory usage: %.2f MB\n", memoryValue)
		fmt.Printf(" I/O read operations: %.2f per second\n", ioReadValue)
		fmt.Printf(" I/O write operations: %.2f per second\n", ioWriteValue)

		avg := (cpuValue + memoryValue + ioReadValue + ioWriteValue)/4
		fmt.Print("average of all is",avg)

		// // Get the HPA configuration for MongoDB.


		hpa, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Get(context.Background(), hpaName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("HPA %s not found, skipping update\n", hpaName)
				continue
			}
			panic(err)
		}

		// // Update the HPA minimum and maximum replica values based on the metric value.
		// // Update the HPA if needed based on the metric value.
		minReplicas, maxReplicas := helper.CalculateReplicas(avg)
		err = helper.UpdateHPA(clientset, hpa, minReplicas, maxReplicas)
		if err != nil {
			panic(err)
		}

		// Wait for the configured interval before updating the HPA again.
		time.Sleep(2 * time.Second)
	}
}
