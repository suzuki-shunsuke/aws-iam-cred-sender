# aws-iam-cred-sender

AWS Lambda Function to send an initial password to a new user

## Overview

## Architecture

```
CloudWatch Event => Lambda Function = inital password => User
```

## Configuration

### Envrionment Variable

`CONFIG`

```yaml
log_level: info
secret_id: xxx
version_id: xxx
region: xxx
slack:
  channel_id_for_system_user: xxx
message_for_system_user:
message:
initial_password_length: 64
```

### AWS Secrets Manager's Secret

* slack_bot_token
* slack_user_token

### Lambda Execution Role

* secretsmanager:GetSecretValue
* iam:CreateLoginProfile

### Slack App Permission

* chat:write (chat.postMessage)
* users:read (users.list)

## LICENSE

[MIT](LICENSE)
