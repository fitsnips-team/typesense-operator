package controller

import (
	"context"
	"fmt"
	tsv1alpha1 "github.com/akyriako/typesense-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func (r *TypesenseClusterReconciler) ReconcileConfigMap(ctx context.Context, ts tsv1alpha1.TypesenseCluster) ([]string, error) {
	configMapName := fmt.Sprintf("%s-nodeslist", ts.Name)
	configMapExists := true
	configMapObjectKey := client.ObjectKey{Namespace: ts.Namespace, Name: configMapName}

	var nodes = &v1.ConfigMap{}
	if err := r.Get(ctx, configMapObjectKey, nodes); err != nil {
		if apierrors.IsNotFound(err) {
			configMapExists = false
		} else {
			r.logger.Error(err, fmt.Sprintf("unable to fetch config map: %s", configMapName))
		}
	}

	if !configMapExists {
		r.logger.Info("creating config map", "configmap", configMapObjectKey.Name)

		_, err := r.createConfigMap(ctx, configMapObjectKey, &ts)
		if err != nil {
			r.logger.Error(err, "creating config map failed", "configmap", configMapObjectKey.Name)
			return nil, err
		}
	} else {
		_, err := r.updateConfigMap(ctx, &ts, nodes)
		if err != nil {
			return nil, err
		}
	}

	nodesList := strings.Replace(nodes.Data["nodes"], fmt.Sprintf(":%d:%d", ts.Spec.PeeringPort, ts.Spec.ApiPort), "", 1)
	nodesSlice := strings.Split(nodesList, ",")
	return nodesSlice, nil
}

const nodeNameLenLimit = 64

func (r *TypesenseClusterReconciler) createConfigMap(ctx context.Context, key client.ObjectKey, ts *tsv1alpha1.TypesenseCluster) (*v1.ConfigMap, error) {
	nodes := make([]string, ts.Spec.Replicas)
	for i := 0; i < int(ts.Spec.Replicas); i++ {
		nodeName := fmt.Sprintf("%s-sts-%d.%s-sts-svc.%s.svc.cluster.local:%d:%d", ts.Name, i, ts.Name, ts.Namespace, ts.Spec.PeeringPort, ts.Spec.ApiPort)
		if len(nodeName) > nodeNameLenLimit {
			return nil, fmt.Errorf("raft error: node name should not exceed %d characters: %s", nodeNameLenLimit, nodeName)
		}

		nodes[i] = fmt.Sprintf("%s:%d:%d", nodeName, ts.Spec.PeeringPort, ts.Spec.ApiPort)
	}

	cm := &v1.ConfigMap{
		ObjectMeta: getObjectMeta(ts, &key.Name),
		Data: map[string]string{
			"nodes": strings.Join(nodes, ","),
		},
	}

	err := ctrl.SetControllerReference(ts, cm, r.Scheme)
	if err != nil {
		return nil, err
	}

	err = r.Create(ctx, cm)
	if err != nil {
		return nil, err
	}

	return cm, nil
}

func (r *TypesenseClusterReconciler) updateConfigMap(ctx context.Context, ts *tsv1alpha1.TypesenseCluster, cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	nodes := make([]string, 0)
	pods, err := r.getPods(ctx, ts)
	if err != nil {
		return nil, err
	}

	desired := cm.DeepCopy()

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if container.Name == "typesense" && strings.TrimSpace(pod.Status.PodIP) != "" && pod.Status.ContainerStatuses[0].Ready {
				nodes = append(nodes, fmt.Sprintf("%s:%d:%d", pod.Status.PodIP, ts.Spec.PeeringPort, ts.Spec.ApiPort))
			}
		}
	}

	availableNodes := len(nodes)
	if availableNodes == 0 {
		r.logger.Info("empty quorum configuration")
		return nil, fmt.Errorf("empty quorum configuration")
	}

	desired.Data = map[string]string{
		"nodes": strings.Join(nodes, ","),
	}

	r.logger.Info("quorum configuration", "nodes", availableNodes, "nodes", nodes)

	if cm.Data["nodes"] != desired.Data["nodes"] {
		r.logger.Info("updating quorum configuration")

		err := r.Update(ctx, desired)
		if err != nil {
			r.logger.Error(err, "updating quorum configuration failed")
			return nil, err
		}
	}

	return desired, nil
}

func (r *TypesenseClusterReconciler) getPods(ctx context.Context, ts *tsv1alpha1.TypesenseCluster) (*v1.PodList, error) {
	listOptions := []client.ListOption{
		client.InNamespace(ts.Namespace),
		client.MatchingLabels(getLabels(ts)),
	}

	pods := &v1.PodList{}
	err := r.List(ctx, pods, listOptions...)
	if err != nil {
		r.logger.Error(err, "failed to list quorum pods")
		return nil, err
	}

	if len(pods.Items) == 0 {
		r.logger.Info("no pods found in quorum")
		return nil, fmt.Errorf("no pods found in quorum")
	}

	return pods, nil
}
