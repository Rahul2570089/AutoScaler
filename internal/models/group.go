package models

type Group struct {
	Name                   string
	LaunchSpec             LaunchSpec
	MinCapacity            int
	MaxCapacity            int
	DesiredCapacity        int
	Zones                  []string
	HealthCheckType        string
	HealthCheckGracePeriod int64
	DefaultCooldown        int64
	TerminationPolicies    []string
	Version                int64
}

type LaunchSpec struct {
	ImageID        string
	InstanceType   string
	KeyName        string
	SecurityGroups []string
	SubnetIDs      []string
	IAMProfile     string
	UserData       string
	Tags           map[string]string
}
