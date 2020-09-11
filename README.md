# MailService Popper

MailService Popper is the Lambda function that is triggered when the mail service SQS is pushed to. It will
build an email using the SQS message and send it via SES.

## Setup
Ensure the following environment variables are set, usually with an .env file.

* TEST_TO (comma separated emails)
* TEST_CC (comma separated emails)
* TEST_BCC (comma separated emails)
* TEST_REPLY_TO (comma separated emails)
* TEST_FROM (single email address)
* TEST_SUBJECT

* AWS_FROM_ARN
* AWS_REGION
* AWS_PROFILE
* GO_ENV (development or production)

## Testing
Currently, the only test file is email_test.go, this will test sending an email via SES.

## Deployment
Deployment will happen automatically via CodePipeline that is managed via a [CDK stack](https://github.com/2pisoftware/Cmfive-Mail-Service-CDK).