// Copyright © 2020 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

func UpdateConfigMap(client kubernetes.Interface, desiredCm *corev1.ConfigMap) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		currentCm, err := client.CoreV1().ConfigMaps(desiredCm.Namespace).Get(desiredCm.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if equality.Semantic.DeepEqual(desiredCm, currentCm) {
			return nil
		}
		_, err = client.CoreV1().ConfigMaps(desiredCm.Namespace).Update(desiredCm)
		return err
	})
}
