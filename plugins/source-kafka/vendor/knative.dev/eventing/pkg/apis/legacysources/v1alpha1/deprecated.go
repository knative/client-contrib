/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

const (
	// StatusConditionTypeDeprecated is the status.conditions.type used to provide deprecation
	// warnings.
	StatusConditionTypeDeprecated = "Deprecated"
)

type Deprecated struct{}

// MarkDeprecated adds a warning condition that this object's spec is using deprecated fields
// and will stop working in the future.
func (d *Deprecated) MarkDeprecated(s *duckv1.Status, reason, msg string) {
	dc := apis.Condition{
		Type:               StatusConditionTypeDeprecated,
		Reason:             reason,
		Status:             corev1.ConditionTrue,
		Severity:           apis.ConditionSeverityWarning,
		Message:            msg,
		LastTransitionTime: apis.VolatileTime{Inner: metav1.NewTime(time.Now())},
	}
	for i, c := range s.Conditions {
		if c.Type == dc.Type {
			s.Conditions[i] = dc
			return
		}
	}
	s.Conditions = append(s.Conditions, dc)
}
