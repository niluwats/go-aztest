package config

import (
	"flag"
)

// AddFlags adds flags applicable to all services.
// Remember to call `flag.Parse()` in your main or TestMain.
func AddFlags() error {
	flag.StringVar(&Subscription_ID, "subscription", Subscription_ID, "Subscription for tests.")
	flag.StringVar(&Location_Default, "location", Location_Default, "Default location for tests.")
	flag.StringVar(&Cloud_Name, "cloud", Cloud_Name, "Name of Azure cloud.")
	flag.StringVar(&BaseGroup_Name, "baseGroupName", BaseGroupName(), "Specify prefix name of resource group for sample resources.")

	flag.BoolVar(&UseDevice_Flow, "useDeviceFlow", UseDevice_Flow, "Use device-flow grant type rather than client credentials.")
	flag.BoolVar(&Keep_Resources, "keepResources", Keep_Resources, "Keep resources created by samples.")

	return nil
}
