# Cmfive Mail Service Queue Trigger

## Introduction
The Cmfive Mail Service Queue Trigger is the Lambda function that is triggered when the mail service SQS is pushed to. It will
build an email using the SQS message and send it via SES.

## Testing
Ensure the following environment variables are set before running tests, usually with an .env file.

* TEST_TO (comma separated emails)
* TEST_CC (comma separated emails)
* TEST_BCC (comma separated emails)
* TEST_REPLY_TO (comma separated emails)
* TEST_FROM (single email address)
* TEST_SUBJECT
* AWS_FROM_ARN
* AWS_REGION
* AWS_PROFILE
* ENVIRONMENT (development or production)

## Deployment
Deployment will happen automatically via CodePipeline that is managed via a [CDK stack](https://github.com/2pisoftware/Cmfive-Mail-Service-CDK).