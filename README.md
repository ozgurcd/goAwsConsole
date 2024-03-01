# goAwsConsole

This application works as CLI and opens the AWS console in the default browser using the profile settings preconfigured.

Since it uses AWS STS to assume the role, it is necessary to pre-configure the role in the AWS before using it. By default, it uses the role 'SREAccess'. 


Usage:
```bash
goAwsConsole
```