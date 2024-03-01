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
	"github.com/ozgurcd/goAwsConsole/models"
)

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

func GetSTSCredentials(config models.RuntimeConfig) {
	stsClient := sts.NewFromConfig(AwsConfig)

	callerIdentity, err := stsClient.GetCallerIdentity(
		context.TODO(),
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		log.Fatalf("unable to get caller identity, %v", err)
		return
	}

	currentUser, err := user.Current()
	if err != nil {
		currentUser = &user.User{Username: "unknown"}
	}

	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", *callerIdentity.Account, config.RoleName)
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(currentUser.Username),
		DurationSeconds: &config.Duration,
	}

	result, err := stsClient.AssumeRole(context.TODO(), input)
	if err != nil {
		log.Fatalf("unable to assume role, %v", err)
	}

	Credentials := models.AwsCredentials{
		SessionId:    *result.Credentials.AccessKeyId,
		SessionKey:   *result.Credentials.SecretAccessKey,
		SessionToken: *result.Credentials.SessionToken,
	}

	creds, err := json.Marshal(Credentials)
	if err != nil {
		log.Fatalf("unable to marshal credentials, %v", err)
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
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read response body, %v", err)
	}

	var federationResponse models.AwsFederationResponse
	err = json.Unmarshal(body, &federationResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling federation response, %v", err)
	}

	destinationURL := url.QueryEscape(
		fmt.Sprintf("%sconsole/home?region=%s", awsTEMPConsoleUrl, Region))

	loginURL := fmt.Sprintf("%s?Action=login&Issuer=goAwsConsole&Destination=%s&SigninToken=%s", awsFederationURL, destinationURL, federationResponse.SigninToken)

	var args []string

	switch runtime.GOOS {
	case "darwin":
		var uniqueEnv string
		if config.SeparateWin {
			if config.ProfileDir == "" {
				currentIndex := time.Now().UnixNano() % 26
				endIndex := currentIndex + 6
				if endIndex > 26 {
					endIndex = endIndex % 26
					uniqueEnv = "abcdefghijklmnopqrstuvwxyz"[currentIndex:] + "abcdefghijklmnopqrstuvwxyz"[:endIndex]
				} else {
					uniqueEnv = "abcdefghijklmnopqrstuvwxyz"[currentIndex:endIndex]
				}
			} else {
				uniqueEnv = config.ProfileDir
			}
			args = []string{
				"open",
				"-na",
				config.Browser,
				"--args",
				fmt.Sprintf("--profile-directory=\"%s\"", uniqueEnv),
				"--new-window"}

		} else {
			args = []string{
				"open",
				"-na",
				config.Browser,
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
