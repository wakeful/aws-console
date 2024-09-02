package console

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAWSConfig(ctx context.Context) (*aws.Config, error) {
	sess, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading default aws config: %w", err)
	}

	if sess.Region == "" {
		sess.Region = DefaultRegion
		slog.Debug("setting default aws region", "region", sess.Region)
	}

	return &sess, nil
}
