package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"goaztest/internal/config"
	"goaztest/internal/iam"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/go-autorest/autorest"

	"goaztest/internal/util"

	"goaztest/resources"

	"goaztest/network"

	"github.com/Azure/go-autorest/autorest/to"
)

func main() {
	// command := fmt.Sprintf("ssh-keygen -b 2048 -t rsa -f " + "test" + " -q -N \"\"\"\"")
	// var stdout bytes.Buffer
	// var stderr bytes.Buffer
	// cmd := exec.Command("bash", "-c", command)
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Fatal(err)
	// }
	// vm, err := CreateVM(context.TODO(), "vmtestN2", "IPtestN2", "vmuser", "vmpass", "test")
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Fatal(err)
	// }
	// fmt.Println(vm)
	Example_createVM()
}

func Example_createVM() {
	var groupName = config.GenerateGroupName("resgrp-nilu1")
	// TODO: remove and use local `groupName` only
	config.SetGroupName(groupName)

	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()
	// defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		util.LogAndPanic(err)
	}

	_, err = network.CreateVirtualNetworkAndSubnets(ctx, "niluvnet1", "nilusubnet1-1", "nilusubnet1-2")
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("created vnet and 2 subnets")

	_, err = network.CreateNetworkSecurityGroup(ctx, "nilunsg1")
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("created network security group")

	_, err = network.CreatePublicIP(ctx, "niluIp1")
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("created public IP")

	_, err = network.CreateNIC(ctx, "niluvnet1", "nilusubnet1-1", "nilunsg1", "niluIp1", "niluNIC1")
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("created nic")

	_, err = CreateVM(ctx, "nilutestVm1", "niluNIC1", "vmuser1", "vmPassword@1", "test.pub")
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("created VM")

	// set or change VM metadata
	_, err = UpdateVM(ctx, "nilutestVm1", map[string]*string{
		"runtime": to.StringPtr("go"),
		"cloud":   to.StringPtr("azure"),
	})
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("updated VM")

	// set or change system state
	_, err = StartVM(ctx, "nilutestVm1")
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("started VM")

	// _, err = RestartVM(ctx, "nilutestVm1")
	// if err != nil {
	// 	util.LogAndPanic(err)
	// }
	// util.PrintAndLog("restarted VM")

	// _, err = StopVM(ctx, "nilutestVm1")
	// if err != nil {
	// 	util.LogAndPanic(err)
	// }
	// util.PrintAndLog("stopped VM")

	// Output:
	// created vnet and 2 subnets
	// created network security group
	// created public IP
	// created nic
	// created VM
	// updated VM
	// started VM
	// restarted VM
	// stopped VM
}

func getVMClient() compute.VirtualMachinesClient {
	vmClient := compute.NewVirtualMachinesClient(config.SubscriptionID())
	a, _ := iam.GetResourceManagementAuthorizer()
	vmClient.Authorizer = a
	vmClient.AddToUserAgent(config.UserAgent())
	return vmClient
}

func getVMExtensionsClient() compute.VirtualMachineExtensionsClient {
	extClient := compute.NewVirtualMachineExtensionsClient(config.SubscriptionID())
	a, _ := iam.GetResourceManagementAuthorizer()
	extClient.Authorizer = a
	extClient.AddToUserAgent(config.UserAgent())
	return extClient
}

// CreateVM creates a new virtual machine with the specified name using the specified NIC.
// Username, password, and sshPublicKeyPath determine logon credentials.
func CreateVM(ctx context.Context, vmName, nicName, username, password, sshPublicKeyPath string) (vm compute.VirtualMachine, err error) {
	// see the network samples for how to create and get a NIC resource
	nic, _ := network.GetNic(ctx, nicName)
	fmt.Println(nic)

	var sshKeyData string
	if _, err = os.Stat(sshPublicKeyPath); err == nil {
		sshBytes, err := ioutil.ReadFile(sshPublicKeyPath)
		if err != nil {
			log.Fatalf("failed to read SSH key data: %v", err)
		}
		sshKeyData = string(sshBytes)
	} else {
		sshKeyData = ""
	}

	vmClient := getVMClient()
	future, err := vmClient.CreateOrUpdate(
		ctx,
		config.GroupName(),
		vmName,
		compute.VirtualMachine{
			Location: to.StringPtr(config.Location()),
			VirtualMachineProperties: &compute.VirtualMachineProperties{
				HardwareProfile: &compute.HardwareProfile{
					VMSize: compute.VirtualMachineSizeTypesBasicA0,
				},
				StorageProfile: &compute.StorageProfile{
					ImageReference: &compute.ImageReference{
						Publisher: to.StringPtr("Canonical"),
						Offer:     to.StringPtr("0001-com-ubuntu-server-focal"),
						Sku:       to.StringPtr("20_04-lts"),
						Version:   to.StringPtr("20.04.202010140"),
					},
				},
				OsProfile: &compute.OSProfile{
					ComputerName:  to.StringPtr(vmName),
					AdminUsername: to.StringPtr(username),
					AdminPassword: to.StringPtr(password),
					LinuxConfiguration: &compute.LinuxConfiguration{
						SSH: &compute.SSHConfiguration{
							PublicKeys: &[]compute.SSHPublicKey{
								{
									Path: to.StringPtr(
										fmt.Sprintf("/home/%s/.ssh/authorized_keys",
											username)),
									KeyData: to.StringPtr(sshKeyData),
								},
							},
						},
					},
				},
				NetworkProfile: &compute.NetworkProfile{
					NetworkInterfaces: &[]compute.NetworkInterfaceReference{
						{
							ID: nic.ID,
							NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
								Primary: to.BoolPtr(true),
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return vm, fmt.Errorf("cannot create vm: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return vm, fmt.Errorf("cannot get the vm create or update future response: %v", err)
	}

	return future.Result(vmClient)
}

// GetVM gets the specified VM info
func GetVM(ctx context.Context, vmName string) (compute.VirtualMachine, error) {
	vmClient := getVMClient()
	return vmClient.Get(ctx, config.GroupName(), vmName, compute.InstanceView)
}

// UpdateVM modifies the VM resource by getting it, updating it locally, and
// putting it back to the server.
func UpdateVM(ctx context.Context, vmName string, tags map[string]*string) (vm compute.VirtualMachine, err error) {

	// get the VM resource
	vm, err = GetVM(ctx, vmName)
	if err != nil {
		return
	}

	// update it
	vm.Tags = tags

	// PUT it back
	vmClient := getVMClient()
	future, err := vmClient.CreateOrUpdate(ctx, config.GroupName(), vmName, vm)
	if err != nil {
		return vm, fmt.Errorf("cannot update vm: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return vm, fmt.Errorf("cannot get the vm create or update future response: %v", err)
	}

	return future.Result(vmClient)
}

// DeallocateVM deallocates the selected VM
func DeallocateVM(ctx context.Context, vmName string) (osr autorest.Response, err error) {
	vmClient := getVMClient()
	future, err := vmClient.Deallocate(ctx, config.GroupName(), vmName)
	if err != nil {
		return osr, fmt.Errorf("cannot deallocate vm: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return osr, fmt.Errorf("cannot get the vm deallocate future response: %v", err)
	}

	return future.Result(vmClient)
}

// StartVM starts the selected VM
func StartVM(ctx context.Context, vmName string) (osr autorest.Response, err error) {
	vmClient := getVMClient()
	future, err := vmClient.Start(ctx, config.GroupName(), vmName)
	if err != nil {
		return osr, fmt.Errorf("cannot start vm: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return osr, fmt.Errorf("cannot get the vm start future response: %v", err)
	}

	return future.Result(vmClient)
}

// RestartVM restarts the selected VM
func RestartVM(ctx context.Context, vmName string) (osr autorest.Response, err error) {
	vmClient := getVMClient()
	future, err := vmClient.Restart(ctx, config.GroupName(), vmName)
	if err != nil {
		return osr, fmt.Errorf("cannot restart vm: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return osr, fmt.Errorf("cannot get the vm restart future response: %v", err)
	}

	return future.Result(vmClient)
}

// StopVM stops the selected VM
func StopVM(ctx context.Context, vmName string) (osr autorest.Response, err error) {
	vmClient := getVMClient()
	// skipShutdown parameter is optional, we are taking its default value here
	future, err := vmClient.PowerOff(ctx, config.GroupName(), vmName, nil)
	if err != nil {
		return osr, fmt.Errorf("cannot power off vm: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return osr, fmt.Errorf("cannot get the vm power off future response: %v", err)
	}

	return future.Result(vmClient)
}
