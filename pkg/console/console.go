package console

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	DefaultRegion        = "eu-west-1"
	federationURL string = "https://signin.aws.amazon.com/federation"
)

var errMissingAuthToken = errors.New("missing auth token")

func fmtURL(token, targetRegion string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", errMissingAuthToken
	}

	userURL, err := url.Parse(federationURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse federation URL: %w", err)
	}

	region := DefaultRegion
	if targetRegion != "" {
		region = targetRegion
	}

	slog.Debug("using region", slog.String("region", region))

	parameters := url.Values{}
	parameters.Add("Action", "login")
	parameters.Add("Issuer", "awslogin")
	parameters.Add("Destination", fmt.Sprintf("https://%s.console.aws.amazon.com/", region))
	parameters.Add("SigninToken", token)

	userURL.RawQuery = parameters.Encode()

	return userURL.String(), nil
}

func GetSignInURL(ctx context.Context, sess aws.Config, region, policyARN string) (string, error) {
	payload, err := buildPayload(ctx, sess, policyARN)
	if err != nil {
		return "", fmt.Errorf("failed to build payload: %w", err)
	}

	token, err := getAuthToken(ctx, payload)
	if err != nil {
		return "", fmt.Errorf("failed to get auth token: %w", err)
	}

	sigURL, err := fmtURL(token, region)
	if err != nil {
		return "", fmt.Errorf("failed to fmt-url: %w", err)
	}

	return sigURL, nil
}
