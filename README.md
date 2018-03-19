# MCA Trade Lead Data Transformer

[![Build Status](https://travis-ci.org/GovWizely/lambda-mca-data.svg)](https://travis-ci.org/GovWizely/lambda-mca-data)
[![Go Report Card](https://goreportcard.com/badge/github.com/GovWizely/lambda-mca-data)](https://goreportcard.com/report/github.com/GovWizely/lambda-mca-data)
![Contributers](https://img.shields.io/github/contributors/GovWizely/lambda-mca-data.svg?maxAge=2592000) 

 **`lambda-mca-data`** is a scraper/transformer of [MCA Trade Lead data from dgMarket](http://www.dgmarket.com/tenders/ShowRssFeeds.do) 
                       that runs on [AWS Lambda](https://aws.amazon.com/lambda/) using the 
                       [Serverless Go](https://serverless.com/framework/docs/providers/aws/examples/hello-world/go/#hello-world-go-example) 
                       framework. It extracts country name and country code from the multi-value `<category>` attribute, and uploads
                       the resulting items as well-formed JSON to the `mca-data.json` file in the `trade-leads` S3 
                       bucket for the account you have specified via your AWS credentials. *The bucket must already be
                       created.* If there 
		               are no items in the source feed or the feed cannot be parsed, the lambda exits.

## Installation

### Prerequisites

Verify you have the latest Go (1.10+) and `dep`. For MacOS users, Homebrew makes it easy:

```bash
brew install go dep
```

Follow the [AWS Installation instructions for Serverless](https://serverless.com/framework/docs/providers/aws/guide/installation/) to install the Serverless framework.

Set your AWS Credentials one of [these ways](https://serverless.com/framework/docs/providers/aws/guide/credentials/).

### Build From Source

Be sure you are in your Go source directory:

```bash
cd ~/go/src
```

Clone the repo:
```bash
git clone https://github.com/GovWizely/lambda-mca-data.git
cd lambda-mca-data
```

Compile the executable:
```bash
make build
```

## Deploy to AWS

```bash
sls deploy -v
```

## Usage

### Invoke function

You can invoke the Lambda in AWS from the command line:

```bash
sls invoke -f mca -l
```

You can also set up triggers in AWS to invoke the Lambda, such as a daily cron.

## Contributing

Contributions and feedback are welcome! Proposals and pull requests will be considered and responded to. 

### Tests

To run the unit tests:

```bash
cd mca
go test -cover
```

## License
This project is licensed under the [MIT](https://github.com/GovWizely/lambda-mca-data/blob/master/LICENSE) license.
