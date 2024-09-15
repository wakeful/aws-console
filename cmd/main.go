package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/wakeful/aws-console/pkg/console"
)

func main() {
	policy := flag.String("policy", "", "assume policy arn, e.q. arn:aws:iam::aws:policy/AdministratorAccess")
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

	if err := openConsole(ctx, region, policy); err != nil {
		slog.Error("missing aws credentials", slog.String("error", err.Error()))

		os.Exit(1)
	}
}

func openConsole(ctx context.Context, region *string, policy *string) error {
	sess, cRegion, awsErr := console.GetAWSConfig(ctx, *region)
	if awsErr != nil {
		return fmt.Errorf("missing aws credentials: %w", awsErr)
	}

	consoleURL, awsErr := console.GetSignInURL(ctx, *sess, cRegion, *policy)
	if awsErr != nil {
		return fmt.Errorf("failed to construct signIn URL: %w", awsErr)
	}

	_ = browser.OpenURL("https://signin.aws.amazon.com/oauth?Action=logout")

	const timeout = 2

	time.Sleep(timeout * time.Second)

	if err := browser.OpenURL(consoleURL); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Please open the following URL in your browser: %s\n", consoleURL)

		return fmt.Errorf("please open the following URL in your browser: %w", err)
	}

	return nil
}
