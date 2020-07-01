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

package client

import (
	"k8s.io/client-go/rest"

	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
	"github.com/maximilien/kn-source-pkg/pkg/types/typesfakes"
	client_testing "k8s.io/client-go/testing"
	"knative.dev/client-contrib/plugins/source-kafka/pkg/types"
	"knative.dev/eventing-contrib/kafka/source/pkg/client/clientset/versioned/typed/sources/v1alpha1/fake"
)

// NewFakeKafkaSourceClient is to create a fake KafkaSourceClient to test
func NewFakeKafkaSourceClient(ns string) types.KafkaSourceClient {
	kafkaParams := NewFakeKafkaSourceParams()
	knFakeSourceClient := &typesfakes.FakeKnSourceClient{}
	knFakeSourceClient.KnSourceParamsReturns(kafkaParams.KnSourceParams)
	knFakeSourceClient.NamespaceReturns(ns)
	knFakeSourceClient.RestConfigReturns(&rest.Config{})

	fakeClientTest := fake.FakeSourcesV1alpha1{Fake: &client_testing.Fake{}}

	return &kafkaSourceClient{
		namespace:         ns,
		kafkaSourceParams: kafkaParams,
		client:            &fakeClientTest,
		knSourceClient:    knFakeSourceClient,
	}
}

// NewFakeKafkaSourceParams is to create a fake KafkaSourceParams to test
func NewFakeKafkaSourceParams() *types.KafkaSourceParams {
	return &types.KafkaSourceParams{
		KnSourceParams: &sourcetypes.KnSourceParams{},
	}
}
