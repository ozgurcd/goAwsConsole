package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"os/user"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AwsCredentials struct {
	SessionId    string `json:"sessionId"`
	SessionKey   string `json:"sessionKey"`
	SessionToken string `json:"sessionToken"`
}

type AwsFederationResponse struct {
	SigninToken string `json:"SigninToken"`
}

const (
	consoleSessDuration  = time.Duration(30) * time.Minute
	awsFederationURL     = "https://signin.aws.amazon.com/federation"
	awsFederationURLTemp = "signin.aws.amazon.com/federation"
	awsTEMPConsoleUrl    = "https://console.aws.amazon.com/"
)

var (
	AwsConfig aws.Config
	Region    string
)

func InitAWS(profile string, region string) {
	var cfg aws.Config
	var err error

	Region = region

	options := []func(*config.LoadOptions) error{}

	if profile != "" {
		options = append(options, config.WithSharedConfigProfile(profile))
	}
	options = append(options, config.WithRegion(region))

	cfg, err = config.LoadDefaultConfig(context.TODO(), options...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	AwsConfig = cfg
}

func GetSTSCredentials(rolename string, duration int32, browser string, separateWindow bool, profileDir string) {
	stsClient := sts.NewFromConfig(AwsConfig)

	callerIdentity, err := stsClient.GetCallerIdentity(
		context.TODO(),
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		return
	}

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("unable to get current user, %v", err)
		currentUser = &user.User{Username: "unknown"}
	}

	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", *callerIdentity.Account, rolename)
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(currentUser.Username),
		DurationSeconds: &duration,
	}

	result, err := stsClient.AssumeRole(context.TODO(), input)
	if err != nil {
		log.Fatalf("unable to assume role, %v", err)
	}

	Credentials := AwsCredentials{
		SessionId:    *result.Credentials.AccessKeyId,
		SessionKey:   *result.Credentials.SecretAccessKey,
		SessionToken: *result.Credentials.SessionToken,
	}

	creds, err := json.Marshal(Credentials)
	if err != nil {
		log.Fatalf("unable to marshal credentials, %v", err)
		return
	}

	consoleUrl := fmt.Sprintf(
		"https://%s.%s?Action=getSigninToken&SessionDuration=%d&Session=%s",
		Region,
		awsFederationURLTemp,
		int64(consoleSessDuration.Seconds()),
		url.QueryEscape(string(creds)))

	resp, err := http.Post(consoleUrl, "application/x-www-form-urlencoded", nil)
	if err != nil {
		log.Fatalf("unable to get signin token, %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read response body, %v", err)
		return
	}

	var federationResponse AwsFederationResponse
	err = json.Unmarshal(body, &federationResponse)
	if err != nil {
		fmt.Printf("Decode AWS Federated response: %v", err)
		return
	}

	destinationURL := url.QueryEscape(
		fmt.Sprintf("%sconsole/home?region=%s", awsTEMPConsoleUrl, Region))

	loginURL := fmt.Sprintf("%s?Action=login&Issuer=goAwsConsole&Destination=%s&SigninToken=%s", awsFederationURL, destinationURL, federationResponse.SigninToken)

	var args []string

	switch runtime.GOOS {
	case "darwin":
		var uniqueEnv string
		if separateWindow {
			if profileDir == "" {
				currentIndex := time.Now().UnixNano() % 26
				endIndex := currentIndex + 6
				if endIndex > 26 {
					endIndex = endIndex % 26
					uniqueEnv = "abcdefghijklmnopqrstuvwxyz"[currentIndex:] + "abcdefghijklmnopqrstuvwxyz"[:endIndex]
				} else {
					uniqueEnv = "abcdefghijklmnopqrstuvwxyz"[currentIndex:endIndex]
				}
			} else {
				uniqueEnv = profileDir
			}
			args = []string{
				"open",
				"-na",
				browser,
				"--args",
				fmt.Sprintf("--profile-directory=\"%s\"", uniqueEnv),
				"--new-window"}

		} else {
			args = []string{
				"open",
				"-na",
				browser,
				"--args",
			}
		}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}

	cmd := exec.Command(args[0], append(args[1:], loginURL)...)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("unable to open browser, %v", err)
	}
}
