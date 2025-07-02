package config

type JobTransferRule struct {
	// From is the CSA public key (hex-encoded) of the source job distributor
	From string
	// To is the CSA public key (hex-encoded) of the target job distributor
	To string
}

type JobDistributor interface {
	DisplayName() string
	AllowedJobTransfers() []JobTransferRule
}
