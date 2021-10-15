package config

import (
	"log"
	"os"
	"strconv"
)

// ParseEnvironment loads a sibling `.env` file then looks through all environment
// variables to set global configuration.
func ParseEnvironment() error {
	os.Setenv("AZURE_GROUP_NAME", "az-testnilu")
	os.Setenv("AZURE_BASE_GROUP_NAME", "az-testnilu")
	os.Setenv("AZURE_LOCATION_DEFAULT", "westus")
	os.Setenv("AZURE_USE_DEVICEFLOW", "true")
	os.Setenv("AZURE_SAMPLES_KEEP_RESOURCES", "0")
	os.Setenv("AZURE_CLIENT_ID", "9fde7425-0ee7-4cdc-bef8-59ed0b5d927d")
	os.Setenv("AZURE_CLIENT_SECRET", "673fb7dc-49ad-4592-90df-e9ac733d62f0")
	os.Setenv("AZURE_TENANT_ID", "de4a51f3-b0b4-45d4-aa17-9da882dfe409")
	os.Setenv("AZURE_SUBSCRIPTION_ID", "4f4484af-dd57-4d42-8d95-b679a8581148")
	// AZURE_GROUP_NAME and `config.GroupName()` are deprecated.
	// Use AZURE_BASE_GROUP_NAME and `config.GenerateGroupName()` instead.
	Group_Name = os.Getenv("AZURE_GROUP_NAME")
	BaseGroup_Name = os.Getenv("AZURE_BASE_GROUP_NAME")

	Location_Default = os.Getenv("AZURE_LOCATION_DEFAULT")

	var err error
	UseDevice_Flow, err = strconv.ParseBool(os.Getenv("AZURE_USE_DEVICEFLOW"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_USE_DEVICEFLOW, disabling\n")
		UseDevice_Flow = false
	}
	Keep_Resources, err = strconv.ParseBool(os.Getenv("AZURE_SAMPLES_KEEP_RESOURCES"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_SAMPLES_KEEP_RESOURCES, discarding\n")
		Keep_Resources = false
	}

	// these must be provided by environment
	// clientID
	// clientID = os.Getenv("AZURE_CLIENT_ID")
	Client_ID = os.Getenv("AZURE_CLIENT_ID")

	// clientSecret
	Client_Secret = os.Getenv("AZURE_CLIENT_SECRET")

	// tenantID (AAD)
	Tenant_ID = os.Getenv("AZURE_TENANT_ID")

	// subscriptionID (ARM)
	Subscription_ID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	// Cloud_Name=""
	return nil
}
