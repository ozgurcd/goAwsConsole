# goAwsConsole

This application works as CLI and opens the AWS console in the default browser using the profile settings preconfigured.

Since it uses AWS STS to assume the role, it is necessary to pre-configure the role in the AWS before using it. By default, it uses the role 'SREAccess'. 


Usage:
```bash
goAwsConsole

  -browser string
        Browser to use for opening the console (default "Google Chrome")
  -duration int
        Duration of the assumed role (default 3600)
  -profile string
        AWS profile to use
  -profile-dir string
        Directory to store profiles, only valid for Google Chrome
  -region string
        AWS region to use (default "us-west-2")
  -role string
        AWS role to assume
  -separate-window
        Open the console in a separate window
```

