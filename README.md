# goAwsConsole

This application works as CLI and opens the AWS console in the default browser using the profile settings preconfigured.

Since it uses AWS STS to assume the role, it is necessary to pre-configure the role in the AWS before using it. By default, it uses the role 'SREAccess'. 


Build:
```bash
make
```
will compile the application for your current OS and architecture. If you want to build for a different OS or architecture, the available options are: `linux`, `mac`, `macintel` and `windows`. For example, if you want to build for Linux, you can use the following command to build for a different OS or architecture:

```bash
make linux
```


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

Tested only on MacOS. Please let me know if you can test on Linux & Windows.
