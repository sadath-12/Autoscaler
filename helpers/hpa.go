package helpers

import (
	"context"
	"fmt"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	hpaName   = "mongodb-hpa"
	namespace = "default"
)

func CalculateReplicas(metricValue float64) (int32, int32) {
	// Adjust the replicas based on the metric value.
	var minReplicas, maxReplicas int32
	if metricValue > 0.8 {
		minReplicas = 3
		maxReplicas = 10
	} else if metricValue > 0.6 {
		minReplicas = 2
		maxReplicas = 5
	} else {
		minReplicas = 1
		maxReplicas = 3
	}
	return minReplicas, maxReplicas
}

func UpdateHPA(clientset *kubernetes.Clientset, hpa *autoscalingv2beta2.HorizontalPodAutoscaler, minReplicas int32, maxReplicas int32) error {
	// Update the HPA minimum and maximum replica values based on the metric value.
	if hpa.Spec.MinReplicas != &minReplicas || hpa.Spec.MaxReplicas != maxReplicas {
		hpa.Spec.MinReplicas = &minReplicas
		hpa.Spec.MaxReplicas = maxReplicas
		_, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Update(context.Background(), hpa, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Updated HPA %s to %d/%d replicas\n", hpaName, minReplicas, maxReplicas)
	}
	return nil
}
