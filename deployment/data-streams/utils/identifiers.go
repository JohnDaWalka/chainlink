package utils

import "fmt"

// DonIdentifier generates a unique identifier for a DON based on its ID and name.
func DonIdentifier(donId uint64, donName string) string {
	return fmt.Sprintf("don-%d-%s", donId, donName)
}
