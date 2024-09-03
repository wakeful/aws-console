package console

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type response struct {
	Token string `json:"SigninToken"`
}

var errInvalidCred = errors.New("invalid credentials")

func getAuthToken(ctx context.Context, payload string, region string) (string, error) {
	fURL := federationURL
	if strings.HasPrefix(region, "cn-") {
		fURL = federationCNURL
	}

	tokenURL, err := url.Parse(fURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	parameters := url.Values{}
	parameters.Add("Action", "getSigninToken")
	parameters.Add("SessionType", "json")
	parameters.Add("DurationSeconds", "3600")
	parameters.Add("Session", payload)
	tokenURL.RawQuery = parameters.Encode()

	var output response

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	const timeout = 5 * time.Second

	var (
		resp       *http.Response
		httpClient = &http.Client{
			Timeout: timeout,
		}
	)

	resp, err = httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting signin token: %w", err)
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return "", errInvalidCred
	}

	if errDecodeJSON := json.NewDecoder(resp.Body).Decode(&output); errDecodeJSON != nil {
		return "", fmt.Errorf("error decoding token: %w", errDecodeJSON)
	}

	return output.Token, nil
}

func buildPayload(ctx context.Context, sess aws.Config, policyARN string) (string, error) {
	token, err := sess.Credentials.Retrieve(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	type d struct {
		AccessKeyID     string `json:"sessionId"`
		SecretAccessKey string `json:"sessionKey"`
		SessionToken    string `json:"sessionToken"`
	}

	var data d

	if token.CanExpire {
		data = d{
			AccessKeyID:     token.AccessKeyID,
			SecretAccessKey: token.SecretAccessKey,
			SessionToken:    token.SessionToken,
		}
	} else {
		stsClient := sts.NewFromConfig(sess)

		const duration = 2520

		fedToken, errGetFedToken := stsClient.GetFederationToken(ctx, &sts.GetFederationTokenInput{
			Name:            aws.String("aws-console"),
			DurationSeconds: aws.Int32(duration),
			PolicyArns:      []types.PolicyDescriptorType{{Arn: aws.String(policyARN)}},
		})
		if errGetFedToken != nil {
			return "", fmt.Errorf("failed to get federation token for custom role: %w", errGetFedToken)
		}

		data = d{
			AccessKeyID:     *fedToken.Credentials.AccessKeyId,
			SecretAccessKey: *fedToken.Credentials.SecretAccessKey,
			SessionToken:    *fedToken.Credentials.SessionToken,
		}
	}

	payload, err := json.Marshal(&data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	return string(payload), nil
}
