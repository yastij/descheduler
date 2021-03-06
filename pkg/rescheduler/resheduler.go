/*
Copyright 2017 The Kubernetes Authors.

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

package rescheduler

import (
	"fmt"

	"github.com/aveshagarwal/rescheduler/cmd/rescheduler/app/options"
	"github.com/aveshagarwal/rescheduler/pkg/rescheduler/client"
	eutils "github.com/aveshagarwal/rescheduler/pkg/rescheduler/evictions/utils"
	nodeutil "github.com/aveshagarwal/rescheduler/pkg/rescheduler/node"
	"github.com/aveshagarwal/rescheduler/pkg/rescheduler/strategies"
)

func Run(rs *options.ReschedulerServer) error {

	rsclient, err := client.CreateClient(rs.KubeconfigFile)
	if err != nil {
		return err
	}
	rs.Client = rsclient

	reschedulerPolicy, err := LoadPolicyConfig(rs.PolicyConfigFile)
	if err != nil {
		return err
	}
	if reschedulerPolicy == nil {
		return fmt.Errorf("\nreschedulerPolicy is nil\n")

	}
	evictionPolicyGroupVersion, err := eutils.SupportEviction(rs.Client)
	if err != nil || len(evictionPolicyGroupVersion) == 0 {
		return err
	}

	stopChannel := make(chan struct{})
	nodes, err := nodeutil.ReadyNodes(rs.Client, stopChannel)
	if err != nil {
		return err
	}

	strategies.RemoveDuplicatePods(rs.Client, reschedulerPolicy.Strategies["RemoveDuplicates"], evictionPolicyGroupVersion, nodes)
	strategies.LowNodeUtilization(rs.Client, reschedulerPolicy.Strategies["LowNodeUtilization"], evictionPolicyGroupVersion, nodes)

	return nil
}
