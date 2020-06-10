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

package util

import (
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
)

// SinkToDuckV1Beta1 converts a Destination from duckv1 to duckv1beta1
func SinkToDuckV1Beta1(destination *duckv1.Destination) *duckv1beta1.Destination {
	r := destination.Ref
	return &duckv1beta1.Destination{
		Ref: &corev1.ObjectReference{
			Kind:       r.Kind,
			Namespace:  r.Namespace,
			Name:       r.Name,
			APIVersion: r.APIVersion,
		},
		URI: destination.URI,
	}
}
