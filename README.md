# MailService Popper

MailService Popper is the Lambda function that is triggered when the mail service SQS is pushed to. It will
build an email using the SQS message and send it via SES.

## Setup
Make a copy of .env.example called .env and fill it in.

## Testing
Currently, the only text file is email_test.go, this will test sending an email via SES.

## Deployment
Run ./scripts/deploy.sh from the root of the project with the command line flag of either 'dev' or 'prod' to make a deployment.
The Lambda function must be already created with the name 'MailService_Dev' or 'MailService_Prod' and have permission to upload
Lambda function code via the AWS CLI.