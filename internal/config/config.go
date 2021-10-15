// package config manages loading configuration from environment and command-line params
package config

import (
	"bytes"
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/marstr/randname"
)

var (
	// these are our *global* config settings, to be shared by all packages.
	// each has corresponding public accessors below.
	// if anything requires a `Set` accessor, that indicates it perhaps
	// shouldn't be set here, because mutable vars shouldn't be global.
	Client_ID               string = "9fde7425-0ee7-4cdc-bef8-59ed0b5d927d"
	Client_Secret           string = "673fb7dc-49ad-4592-90df-e9ac733d62f0"
	Tenant_ID               string = "de4a51f3-b0b4-45d4-aa17-9da882dfe409"
	Subscription_ID         string = "4f4484af-dd57-4d42-8d95-b679a8581148"
	Location_Default        string = "westus"
	Authorization_ServerURL string
	Cloud_Name              string = "AzurePublicCloud"
	UseDevice_Flow          bool
	Keep_Resources          bool
	Group_Name              string = "az-testnilu" // deprecated, use baseGroupName instead
	BaseGroup_Name          string = "az-testnilu"
	User_Agent              string
	Environment_            *azure.Environment
)

// ClientID is the OAuth client ID.
func ClientID() string {
	return Client_ID
}

// ClientSecret is the OAuth client secret.
func ClientSecret() string {
	return Client_Secret
}

// TenantID is the AAD tenant to which this client belongs.
func TenantID() string {
	return Tenant_ID
}

// SubscriptionID is a target subscription for Azure resources.
func SubscriptionID() string {
	return Subscription_ID
}

// deprecated: use DefaultLocation() instead
// Location returns the Azure location to be utilized.
func Location() string {
	return Location_Default
}

// DefaultLocation() returns the default location wherein to create new resources.
// Some resource types are not available in all locations so another location might need
// to be chosen.
func DefaultLocation() string {
	return Location_Default
}

// AuthorizationServerURL is the OAuth authorization server URL.
// Q: Can this be gotten from the `azure.Environment` in `Environment()`?
func AuthorizationServerURL() string {
	return Authorization_ServerURL
}

// UseDeviceFlow() specifies if interactive auth should be used. Interactive
// auth uses the OAuth Device Flow grant type.
func UseDeviceFlow() bool {
	return UseDevice_Flow
}

// deprecated: do not use global group names
// utilize `BaseGroupName()` for a shared prefix
func GroupName() string {
	return Group_Name
}

// deprecated: we have to set this because we use a global for group names
// once that's fixed this should be removed
func SetGroupName(name string) {
	Group_Name = name
}

// BaseGroupName() returns a prefix for new groups.
func BaseGroupName() string {
	return BaseGroup_Name
}

// KeepResources() specifies whether to keep resources created by samples.
func KeepResources() bool {
	return Keep_Resources
}

// UserAgent() specifies a string to append to the agent identifier.
func UserAgent() string {
	if len(User_Agent) > 0 {
		return User_Agent
	}
	return "sdk-samples"
}

// Environment() returns an `azure.Environment{...}` for the current cloud.
func Environment() *azure.Environment {
	if Environment_ != nil {
		return Environment_
	}
	env, err := azure.EnvironmentFromName(Cloud_Name)
	if err != nil {
		// TODO: move to initialization of var
		panic(fmt.Sprintf(
			"invalid cloud name '%s' specified, cannot continue\n", Cloud_Name))
	}
	Environment_ = &env
	return Environment_
}

// GenerateGroupName leverages BaseGroupName() to return a more detailed name,
// helping to avoid collisions.  It appends each of the `affixes` to
// BaseGroupName() separated by dashes, and adds a 5-character random string.
func GenerateGroupName(affixes ...string) string {
	// go1.10+
	// import strings
	// var b strings.Builder
	// b.WriteString(BaseGroupName())
	b := bytes.NewBufferString(BaseGroupName())
	b.WriteRune('-')
	for _, affix := range affixes {
		b.WriteString(affix)
		b.WriteRune('-')
	}
	return randname.GenerateWithPrefix(b.String(), 5)
}

// AppendRandomSuffix will append a suffix of five random characters to the specified prefix.
func AppendRandomSuffix(prefix string) string {
	return randname.GenerateWithPrefix(prefix, 5)
}
