ec2-vuls-config
===

ec2-vuls-config is useful cli to create config file for [Vuls](https://github.com/future-architect/vuls) in Amazon EC2.

## How to install and settings

### 1. Installation

```
$ go get -u github.com/ohsawa0515/ec2-vuls-config
```

### 2. Set AWS credentials
 
* Credential file (`$HOME/.aws/credentials`) 

```
[default]
aws_access_key_id = <YOUR_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>
```

* Environment variable

```
$ export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
$ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_ACCESS_KEY
```

### 3. Set AWS region

```
$ export AWS_REGION=us-east-1
```

### 4. Set the `Name` and `Vuls-Scan` tag to EC2 instance that you want to scan

e.g.

```
Name: web-server-1
Vuls-Scan: True
```

### 5. Prepare config.toml for Vuls scan

See [README of Vuls](https://github.com/future-architect/vuls/blob/master/README.md#step6-config) or [config.toml.sample](https://github.com/ohsawa0515/ec2-vuls-config/blob/master/config.toml.sample)

## Usage

### Execute

```
$ ec2-vuls-config
```

After execute, config.toml would be generated as follows.

```
[default]
port        = "22"
user        = "vuls"
keyPath     = "/opt/vuls/.ssh/id_rsa"

[servers]

### Generate by ec2-vuls-config ###
# Updated 2000-01-01T00:01:00+09:00

[servers.web-server-1]
host = "192.0.2.11"

### ec2-vuls-config end ###
```

### Options

#### --config

Specify the file path to the config.toml (Default: `$PWD/config.toml`). 

e.g.

```
$ ec2-vuls-config --config path/to/config.toml
```

#### --filters

Filtering EC2 instances like [describe-instances command](http://docs.aws.amazon.com/cli/latest/reference/ec2/describe-instances.html).  
Also, by default, filtering works that status is running, platform is linux and Vuls-Scan=True tag. 


e.g.

* To scan all instances with a Vuls-Scan=True tag (Default)

```
$ ec2-vuls-config --filters "Name=tag:Vuls-Scan,Values=True"
```

* To scan all instances with the web-server

```
$ ec2-vuls-config --filters "Name=tag:Name,Values=web-server"
```

## Contribution

1. Fork ([https://github.com/ohsawa0515/ec2-vuls-config/fork](https://github.com/ohsawa0515/ec2-vuls-config/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

## License

See [LICENSE](https://github.com/ohsawa0515/ec2-vuls-config/blob/master/LICENSE).
