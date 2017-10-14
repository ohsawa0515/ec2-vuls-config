ec2-vuls-config
===

ec2-vuls-config is useful command line tool to create config file for [Vuls](https://github.com/future-architect/vuls) in Amazon EC2.
By specifying the EC2 tag, you select the scan target Automatically and rewrite the config file.

# Installation

## Step1. Set the `Name` and `vuls:scan` tag to EC2 instances that you want to scan

```console
Name : web-server-1
vuls:scan : true
```

## Step2. Installation

* Binary

Download from [releases page](https://github.com/ohsawa0515/ec2-vuls-config/releases).

* Go get

```console
$ go get -u github.com/ohsawa0515/ec2-vuls-config
```

## Step3. Set AWS credentials
 
* Credential file (`$HOME/.aws/credentials`) 

```console
[default]
aws_access_key_id = <YOUR_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>
```

* Environment variable

```console
$ export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
$ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_ACCESS_KEY
```

## Step4. Set AWS region

```console
$ export AWS_REGION=us-east-1
```


## Step5. Prepare config.toml for Vuls scan

See [vuls#configuration](https://github.com/future-architect/vuls#configuration) or [config.toml.sample](https://github.com/ohsawa0515/ec2-vuls-config/blob/master/config.toml.sample)

## Step6. Execute

By default, it is filtered under the following conditions.

- Status of EC2 instance is running
- Linux (will not select Windows)
- `vuls:scan` tag is set to `true`

```console
$ ec2-vuls-config
```

After execute, config.toml would be rewrites as follows.

```toml
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

# Tags

It can be reflected in config by setting a tag such as `vuls:user`, `vuls:port` and so on.

`<...>` is the name of tag.

```toml
[servers]

[servers.<Name>]
host = "<<Private IP address of instance>>"
port = "<vuls:port>"
user = "<vuls:user>"
keyPath = "<vuls:keyPath>"

# Set value of tag as comma-separated.
cpeNames = [
"<vuls:cpeNames>",
]

# Set value of tag as comma-separated.
ignoreCves = [
"<vuls:ignoreCves>",
]

# Example

# `vuls:user` => vuls
# `vuls:port` => 22
# `vuls:keyPath` => /opt/vuls/.ssh/id_rsa
# `vuls:cpeNames` => cpe:/a:rubyonrails:ruby_on_rails:4.2.7.1,cpe:/a:rubyonrails:ruby_on_rails:4.2.8,cpe:/a:rubyonrails:ruby_on_rails:5.0.1
# `vuls:ignoreCves` => CVE-2014-2913,CVE-2016-6314

[servers.web-server-1]
host = "192.0.2.11"
user = "vuls"
port = "22"
keyPath = "/opt/vuls/.ssh/id_rsa"
cpeNames = [
"cpe:/a:rubyonrails:ruby_on_rails:4.2.7.1",
"cpe:/a:rubyonrails:ruby_on_rails:4.2.8",
"cpe:/a:rubyonrails:ruby_on_rails:5.0.1",
]
ignoreCves = [
"CVE-2014-2913",
"CVE-2016-6314",
]
```

# Command line options

## --config (-c)

Specify the file path to the config.toml to be read.
By default, `$PWD/config.toml`.

e.g.

```console
$ ec2-vuls-config --config /path/to/config.toml
```

## --filters (-f)

In addition to the default condition, it is used for further filter.
This option like [describe-instances command](http://docs.aws.amazon.com/cli/latest/reference/ec2/describe-instances.html).
Specify Name and Value and separate with a space.

e.g.

* To scan all instances with name of `web-server`

```console
$ ec2-vuls-config --filters "Name=tag:Name,Values=web-server"
```

* To scan all instances with name of `app-server` and instance type `c3.large`

```console
$ ec2-vuls-config --filters "Name=tag:Name,Values=app-server Name=instance-type,Values=r3.large"
```

## --out (-o)

Specify the path of the config file to be written.
By default, `$PWD/config.toml`.

e.g.

```console
$ ec2-vuls-config --out /path/to/config.toml
```


## --print (-p)

Echo the standard output instead of write into specified config file.

# Contribution

1. Fork ([https://github.com/ohsawa0515/ec2-vuls-config/fork](https://github.com/ohsawa0515/ec2-vuls-config/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

# License

See [LICENSE](https://github.com/ohsawa0515/ec2-vuls-config/blob/master/LICENSE).
