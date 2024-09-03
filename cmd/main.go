package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/pkg/browser"
	"github.com/wakeful/aws-console/pkg/console"
)

func main() {
	policy := flag.String("policy", "arn:aws:iam::aws:policy/AdministratorAccess", "fall back policy ARN")
	region := flag.String("region", "", "AWS Region")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	ctx := context.Background()

	level := slog.LevelInfo
	if *debug {
		level = slog.LevelDebug
	}

	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource:   false,
			Level:       level,
			ReplaceAttr: nil,
		})))

	sess, cRegion, awsErr := console.GetAWSConfig(ctx, *region)
	if awsErr != nil {
		slog.Error("missing aws credentials", slog.String("error", awsErr.Error()))
		os.Exit(1)
	}

	consoleURL, awsErr := console.GetSignInURL(ctx, *sess, cRegion, *policy)
	if awsErr != nil {
		slog.Error("failed to construct signIn URL", slog.String("error", awsErr.Error()))
		os.Exit(1)
	}

	if err := browser.OpenURL(consoleURL); err != nil {
		slog.Error("failed to open browser", slog.String("error", err.Error()))

		_, _ = fmt.Fprintf(os.Stdout, "Please open the following URL in your browser: %s\n", consoleURL)
	}
}
