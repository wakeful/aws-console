package console

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAWSConfig(ctx context.Context, region string) (*aws.Config, string, error) {
	sess, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("error loading default aws config: %w", err)
	}

	slog.Debug("region in config:", slog.String("region", sess.Region))

	if region != "" {
		sess.Region = region
		slog.Debug("setting region to user selected one", "region", slog.String("region", sess.Region))
	}

	if sess.Region == "" {
		sess.Region = DefaultRegion
		slog.Debug("no region found failing back to default", "region", slog.String("region", sess.Region))
	}

	slog.Debug("using region", "region", sess.Region)

	return &sess, sess.Region, nil
}
