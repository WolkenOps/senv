# senv

To help you to populate your environment variables using AWS SSM Parameter Store. 

## Installation

### Using go

Execute the following command

```shell
$ get -u github.com/luisc09/senv
```

### Using bash 

Execute the following command

```shell
$ curl -L https://github.com/luisc09/senv/releases/download/latest/senv > /usr/local/bin/senv && chmod +x /usr/local/bin/senv
```

### Usage

1. Set your credentials through the AWS CLI. 
1. To format the output as dotenv execute:
    ```shell
    $ senv --paths /dev/myservice/environment/, /dev/global/environment
    ```
1. To format the output with `export` as prefix:
    ```shell
    $ senv --paths /dev/myservice/environment/, /dev/global/environment --export
    ```
1. To save the output to a file just redirect it:
    ```
    $ senv --paths /dev/myservice/environment/, /dev/global/environment --export > .env
    ```

### Docker

1. To build and run it locally execute: 
    ```shell
    $ docker build -t senv -f .docker/Dockerfile .
    $ docker run -v ~/.aws/:/root/.aws/ senv --paths  /dev/myservice/environment/
    ```
    You can additionally any AWS environment variable (in case you do not use `default`)
    
    ```shell
    $ docker run -v ~/.aws/:/root/.aws/ -e AWS_PROFILE=myprofile -e AWS_REGION=us-east-1 senv --paths /dev/myservice/environment/
    ```
1. To use the one in the Docker Hub:
    ```shell
    $ docker run -v ~/.aws/:/root/.aws/ luisc09/senv --paths  /dev/myservice/environment/
    ```
Notes:
The default for `AWS_REGION` is `us-east-1`, to override it use the `--env`/`-e` option. 

To use in an EC2 make sure to override the `AWS_REGION` to the one in which the EC2 is. 