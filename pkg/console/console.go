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
	DefaultRegion          = "eu-west-1"
	federationURL   string = "https://signin.aws.amazon.com/federation"
	federationCNURL string = "https://signin.amazonaws.cn/federation"
)

var errMissingAuthToken = errors.New("missing auth token")

func fmtURL(token, region string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", errMissingAuthToken
	}

	var (
		consoleURL = "console.aws.amazon.com"
		fURL       = federationURL
	)

	if strings.HasPrefix(region, "cn-") {
		consoleURL = "console.amazonaws.cn"
		fURL = federationCNURL
	}

	userURL, err := url.Parse(fURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse federation URL: %w", err)
	}

	slog.Debug("using console", slog.String("url", consoleURL))

	parameters := url.Values{}
	parameters.Add("Action", "login")
	parameters.Add("Issuer", "aws-console")
	parameters.Add("Destination", fmt.Sprintf("https://%s.%s/", region, consoleURL))
	parameters.Add("SigninToken", token)

	userURL.RawQuery = parameters.Encode()

	return userURL.String(), nil
}

func GetSignInURL(ctx context.Context, sess aws.Config, region, policyARN string) (string, error) {
	payload, err := buildPayload(ctx, sess, policyARN)
	if err != nil {
		return "", fmt.Errorf("failed to build payload: %w", err)
	}

	token, err := getAuthToken(ctx, payload, region)
	if err != nil {
		return "", fmt.Errorf("failed to get auth token: %w", err)
	}

	sigURL, err := fmtURL(token, region)
	if err != nil {
		return "", fmt.Errorf("failed to fmt-url: %w", err)
	}

	return sigURL, nil
}
