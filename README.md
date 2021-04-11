# aws-iam-cred-sender

AWS Lambda Function to send an initial password to a new user via Slack DM

## Assumption

IAM User name and Slack's display name must be same.

## Overview

When an IAM User is created, the Lambda Function is triggered via CloudWatch Event.
The function searches the Slack User with IAM User name.
If the IAM User isn't found, the function notifies that the IAM User is created to a Slack channel.
If the IAM User is found, the function creates the IAM User's login profile and sends the initial password to the user via Slack DM.

## Architecture

```
CloudWatch Event => Lambda Function = inital password => User (Slack DM)
```

## Configuration

### Envrionment Variable

`CONFIG`

Default Config

```yaml
log_level: info
region: ""
message_for_system_user: "AWS Account for system user has been created: `{{.UserName}}`"
message: |
  Your AWS Account has been created.
  Initial password:

  {{.Password}}

  Please sign in AWS and change your password.

  https://signin.aws.amazon.com/console

  Please create your AWS Access Key if needed.
  https://console.aws.amazon.com/iam/home#/users/{{.UserName}}?section=security_credentials
initial_password_length: 64
```

Required

```yaml
secret_id: xxx
version_id: xxx
slack:
  channel_id_for_system_user: xxx
```

### AWS Secrets Manager's Secret

* slack_bot_token

### Lambda Execution Role

* secretsmanager:GetSecretValue
* iam:CreateLoginProfile

### Slack App Permission

* chat:write (chat.postMessage)
* users:read (users.list)

## LICENSE

[MIT](LICENSE)
