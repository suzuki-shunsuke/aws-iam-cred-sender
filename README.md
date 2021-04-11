# aws-iam-cred-sender

AWS Lambda Function to send an initial password to a new user via Slack DM

## Assumption

IAM User name and Slack's display name must be same.

## Overview

When an IAM User is created, the Lambda Function is triggered via CloudWatch Event.
The function searches the Slack User with IAM User name.
If the IAM User isn't found, the function notifies that the IAM User is created to a Slack channel.

![image](https://user-images.githubusercontent.com/13323303/114290928-3ba40200-9abe-11eb-8f9b-72b3680d4a1e.png)

If the IAM User is found, the function creates the IAM User's login profile and sends the initial password to the user via Slack DM.

![image](https://user-images.githubusercontent.com/13323303/114290993-bbca6780-9abe-11eb-9efe-ff2376400a96.png)

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
message_for_system_user: "AWS IAM User for system has been created: `{{.UserName}}`"
message: |
  Your AWS IAM User has been created.
  AWS Account ID: `{{.AWSAccountID}}`
  User Name: `{{.UserName}}`
  Initial password:

  {{.Password}}

  Please sign in AWS and change your password.

  https://{{AWSAccountID}}.signin.aws.amazon.com/console

  And please create your AWS Access Key if needed.
  https://console.aws.amazon.com/iam/home#/users/{{.UserName}}?section=security_credentials
initial_password_length: 32
slack:
  channel_id_for_system_user: ""
secret_version_id: ""
```

Required

```yaml
aws_account_id: "xxx"
secret_id: xxx
region: xxx
```

#### Template

The template is rendered with Go's [text/template](https://golang.org/pkg/text/template/).
And [sprig Function](http://masterminds.github.io/sprig/) can be used.

Template variables:

* AWSAccountID
* UserName
* Password

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
