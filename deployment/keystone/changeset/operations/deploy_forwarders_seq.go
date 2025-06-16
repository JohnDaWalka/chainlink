package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"golang.org/x/sync/errgroup"
)

type DeployKeystoneForwardersSequenceDeps struct {
}

type DeployKeystoneForwardersInput struct {
	Targets []uint64 // The target chains for the Keystone Forwarders
}

type DeployKeystoneForwardersOutput struct {
	AddressBook deployment.AddressBook // The address book containing the deployed Keystone Forwarders
}

var DeployKeystoneForwardersSequence = operations.NewSequence[DeployKeystoneForwardersInput, DeployKeystoneForwardersOutput, DeployKeystoneForwardersSequenceDeps](
	"deploy-keystone-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Forwarders",
	func(b operations.Bundle, deps DeployKeystoneForwardersSequenceDeps, input DeployKeystoneForwardersInput) (DeployKeystoneForwardersOutput, error) {
		ab := deployment.NewMemoryAddressBook()
		contractErrGroup := &errgroup.Group{}
		for _, target := range input.Targets {
			fmt.Println(target)
			contractErrGroup.Go(func() error {
				// For each target, we would deploy the Keystone Forwarder.
				// This is a placeholder for the actual deployment logic.
				// TODO: we would pass here the target as an input to the operation.
				_, err := operations.ExecuteOperation(b, DeployKeystoneForwarderOp, DeployForwarderOpDeps{}, DeployForwarderOpInput{})
				if err != nil {
					return err
				}
				//err = ab.Save(target, r.Output.Address.String(), r.Output.Tv)
				//if err != nil {
				//	return pkgerrors.Wrapf(err, "failed to save Keystone Forwarder address for target %d", target)
				//}

				return nil
			})
		}
		if err := contractErrGroup.Wait(); err != nil {
			return DeployKeystoneForwardersOutput{AddressBook: ab}, pkgerrors.Wrap(err, "failed to deploy Keystone contracts")
		}
		return DeployKeystoneForwardersOutput{AddressBook: ab}, nil
	},
)
