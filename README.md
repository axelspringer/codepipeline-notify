# CodePipeline Notify

[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![Build Status](https://travis-ci.org/katallaxie/vue-preboot.svg?branch=master)](https://travis-ci.org/katallaxie/vue-preboot)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

CodePipeline Notify is a Lambda Function which pushes CloudWatch CodePipeline Events to various services (e.g. Slack). It is configured by the System Manager Parameter Store and uses a DynamoDB for the CodePipeline configuration.

## Getting Started

The `template.yml` is a CloudFormation ChangeSet to CloudFormation Recipe available available in [vodka-tf](https://github.com/axelspringer/vodka-tf).

## Environment

We use the System Manager Parameter Store to configure the Lambda Function. The SSM path is configured via an environment variable.

### `SSM_PATH`

This configures the Parameter Store path from which the necessary configuration options are retrieved.

## Parameters (SSM)

### `dynamodb-tablename`

## License
[MIT](/LICENSE)
