// Copyright (c) 2018-2019,2021 Tigera, Inc. All rights reserved.

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

package resourcemgr

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	client "github.com/projectcalico/calico/libcalico-go/lib/clientv3"
	"github.com/projectcalico/calico/libcalico-go/lib/options"
)

func init() {
	registerResource(
		api.NewGlobalNetworkSet(),
		newGlobalNetworkSetList(),
		false,
		[]string{"globalnetworkset", "globalnetworksets", "gnetsets"},
		[]string{"NAME"},
		[]string{"NAME", "NETS"},
		map[string]string{
			"NAME": "{{.ObjectMeta.Name}}",
			"NETS": "{{joinAndTruncate .Spec.Nets \",\" 80}}",
		},
		func(ctx context.Context, client client.Interface, resource ResourceObject) (ResourceObject, error) {
			r := resource.(*api.GlobalNetworkSet)
			return client.GlobalNetworkSets().Create(ctx, r, options.SetOptions{})
		},
		func(ctx context.Context, client client.Interface, resource ResourceObject) (ResourceObject, error) {
			r := resource.(*api.GlobalNetworkSet)
			return client.GlobalNetworkSets().Update(ctx, r, options.SetOptions{})
		},
		func(ctx context.Context, client client.Interface, resource ResourceObject) (ResourceObject, error) {
			r := resource.(*api.GlobalNetworkSet)
			return client.GlobalNetworkSets().Delete(ctx, r.Name, options.DeleteOptions{ResourceVersion: r.ResourceVersion})
		},
		func(ctx context.Context, client client.Interface, resource ResourceObject) (ResourceObject, error) {
			r := resource.(*api.GlobalNetworkSet)
			return client.GlobalNetworkSets().Get(ctx, r.Name, options.GetOptions{ResourceVersion: r.ResourceVersion})
		},
		func(ctx context.Context, client client.Interface, resource ResourceObject) (ResourceListObject, error) {
			r := resource.(*api.GlobalNetworkSet)
			return client.GlobalNetworkSets().List(ctx, options.ListOptions{ResourceVersion: r.ResourceVersion, Name: r.Name})
		},
	)
}

// newNetworkSetList creates a new (zeroed) NetworkSetList struct with the TypeMetadata initialised to the current
// version.
func newGlobalNetworkSetList() *api.GlobalNetworkSetList {
	return &api.GlobalNetworkSetList{
		TypeMeta: metav1.TypeMeta{
			Kind:       api.KindGlobalNetworkSetList,
			APIVersion: api.GroupVersionCurrent,
		},
	}
}
