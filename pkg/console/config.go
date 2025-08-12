// Copyright 2025 variHQ OÜ
// SPDX-License-Identifier: BSD-3-Clause

package console

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// GetAWSConfig loads AWS SDK config with the given or default region, returning the config, region, and any error.
func GetAWSConfig(ctx context.Context, region string) (*aws.Config, string, error) {
	sess, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("error loading default aws config: %w", err)
	}

	slog.Debug("region in config:", slog.String("region", sess.Region))

	if region != "" {
		sess.Region = region
		slog.Debug(
			"setting region to user selected one",
			slog.String("region", sess.Region),
		)
	}

	if sess.Region == "" {
		sess.Region = DefaultRegion
		slog.Debug(
			"no region found failing back to default",
			slog.String("region", sess.Region),
		)
	}

	slog.Debug("using region", "region", sess.Region)

	return &sess, sess.Region, nil
}
