//
// Copyright (c) 2017 Joey <majunjiev@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"fmt"
	"time"

	ovirtsdk4 "gopkg.in/imjoey/go-ovirt.v4"
)

func main() {
	inputRawURL := "https://10.1.111.229/ovirt-engine/api"

	conn, err := ovirtsdk4.NewConnectionBuilder().
		URL(inputRawURL).
		Username("admin@internal").
		Password("qwer1234").
		Insecure(true).
		Compress(true).
		Timeout(time.Second * 10).
		Build()
	if err != nil {
		fmt.Printf("Make connection failed, reason: %v\n", err)
		return
	}
	defer conn.Close()

	// To use `Must` methods, you should recover it if panics
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panics occurs, try the non-Must methods to find the reason")
		}
	}()

	// Get the service that manages the vms
	vmsService := conn.SystemService().VmsService()
	vm := vmsService.List().
		Search("name=myvm").
		MustSend().
		MustVms().
		Slice()[0]

	// Update the placement policy of the virtual machine so that it's pinned to the host
	vmService := vmsService.VmService(vm.MustId())
	vmService.Update().
		Vm(
			ovirtsdk4.NewVmBuilder().
				PlacementPolicy(
					ovirtsdk4.NewVmPlacementPolicyBuilder().
						HostsOfAny(
							*ovirtsdk4.NewHostBuilder().
								Name("myhost").
								MustBuild()).
						MustBuild()).
				MustBuild()).
		MustSend()
}