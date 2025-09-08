package infra

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

func StartS3(testLogger zerolog.Logger, input *s3provider.Input, stageGen *stagegen.StageGen) (*s3provider.Output, error) {
	var s3ProviderOutput *s3provider.Output
	if input != nil {
		fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting MinIO")))
		var s3ProviderErr error
		s3ProviderOutput, s3ProviderErr = s3provider.NewMinioFactory().NewFrom(input)
		if s3ProviderErr != nil {
			return nil, errors.Wrap(s3ProviderErr, "minio provider creation failed")
		}
		testLogger.Debug().Msgf("S3Provider.Output value: %#v", s3ProviderOutput)
		fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("MinIO started in %.2f seconds", stageGen.Elapsed().Seconds())))
	}

	return s3ProviderOutput, nil
}
