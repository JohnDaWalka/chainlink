package addressbook

import (
	"errors"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/deployment"
)

func Test_UsageExample(t *testing.T) {
	//foo(nil)
}

func foo(ab AddressBook) {

	addresses, _ := ab.Addresses().Fetch()
	for _, address := range addresses {
		fmt.Println(address.Labels)
	}

	addresses, err := ab.Addresses().By().Chain(1).QualifierEquals("foo").Fetch()
	if err != nil {
		// do stuff
	}
	for _, address := range addresses {
		fmt.Println(address.Labels)
	}

	key := NewAddressKey(1, "blah", deployment.Version1_0_0, "")
	rec, err := ab.Addresses().By().Id(key)
	if errors.Is(err, ErrRecordNotFound) {
		fmt.Printf("Address could not be found for %v\n", key)
	} else if err != nil {
		// panic
	}
	fmt.Printf("Found record %v\n", rec)

	metadatas, _ := ab.Metadata().By().Chain(3).Fetch()
	for i, m := range metadatas {
		fmt.Printf("Metadata result %d: %v\n", i, m)
	}

	metadata, _ := ab.Metadata().By().Id(NewMetadataKey(3, "0xABCD0123"))
	fmt.Printf("Metadata: %v\n", metadata)

	metadataMap, _ := ab.Metadata().Associate(rec)
	metadata = metadataMap[rec.Key()]

	key = NewAddressKey(1, "blah", deployment.Version1_0_0, "")
	rec, err = ab.Addresses().By().Id(key)
	if errors.Is(err, ErrRecordNotFound) {
		fmt.Printf("Address could not be found for %v\n", key)
	} else if err != nil {
		// panic
	}
	fmt.Printf("Found record %v\n", rec)

}
