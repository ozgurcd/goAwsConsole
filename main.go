package main

import (
	"flag"
	"os"

	"github.com/ozgurcd/goAwsConsole/aws"
	"github.com/ozgurcd/goAwsConsole/models"
)

func main() {
	var profile string
	var rolename string

	profilePtr := flag.String("profile", "", "AWS profile to use")
	roleNamePtr := flag.String("role", "", "AWS role to assume")
	durationPtr := flag.Int("duration", 3600, "Duration of the assumed role")
	regionPtr := flag.String("region", "us-west-2", "AWS region to use")
	browserPtr := flag.String("browser", "Google Chrome", "Browser to use for opening the console")
	separateWindow := flag.Bool("separate-window", false, "Open the console in a separate window")
	profileDir := flag.String("profile-dir", "", "Directory to store Google Chrome profiles")

	flag.Parse()

	if *profilePtr != "" {
		profile = *profilePtr
	} else {
		profile = os.Getenv("AWS_PROFILE")
	}

	if *roleNamePtr != "" {
		rolename = *roleNamePtr
	} else {
		rolename = os.Getenv("AWS_ROLE")
	}

	if rolename == "" {
		rolename = "SREAccess"
	}

	duration := int32(*durationPtr)

	rConfig := models.RuntimeConfig{
		Profile:     profile,
		RoleName:    rolename,
		Duration:    duration,
		Region:      *regionPtr,
		Browser:     *browserPtr,
		SeparateWin: *separateWindow,
		ProfileDir:  *profileDir,
	}

	aws.InitAWS(profile, *regionPtr)
	aws.GetSTSCredentials(rConfig)
}
