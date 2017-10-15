ec2-vuls-config
===

ec2-vuls-config は [Vuls](https://github.com/future-architect/vuls)スキャンのために、Amazon EC2インスタンスの情報を収集して設定ファイルを生成するのに役立つコマンドラインツールです。  
EC2タグを指定することで、自動的にスキャン対象を選定し、設定ファイルを書き換えます。

# Installation

## Step1. スキャンしたいEC2インスタンスに`Name`タグと`vuls:scan`タグとその値を付与する

```console
Name : web-server-1
vuls:scan : true
```

## Step2. インストール

* Binary

[releases page](https://github.com/ohsawa0515/ec2-vuls-config/releases)からダウンロードできます。

* Go get

```console
$ go get -u github.com/ohsawa0515/ec2-vuls-config
$ go get -u github.com/golang/dep/...
$ dep ensure
```

## Step3. AWSクレデンシャルを設定

IAMポリシー例:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:DescribeInstances"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

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

## Step4. AWSリージョンを設定

```console
$ export AWS_REGION=us-east-1
```


## Step5. 設定ファイル(config.toml)を用意する

設定ファイルについては、[vuls#configuration](https://github.com/future-architect/vuls#configuration) か [config.toml.sample](https://github.com/ohsawa0515/ec2-vuls-config/blob/master/config.toml.sample) をご参照ください。

## Step6. 実行

デフォルトで以下のフィルタ条件が適用されています。

- EC2インスタンスのステータスがRunning
- Linux (Windowsは選択されない)
- `vuls:scan` タグの値は `true` のみ

```console
$ ec2-vuls-config
```

実行後, 設定ファイル(config.toml)は以下のように追記されています。

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

`vuls:user`、` vuls:port`などのEC2タグを設定することで、設定ファイルにに反映させることができます。

`<...>` はタグ名です。

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

読み込む設定ファイルのファイルパスを指定します。デフォルト: `$PWD/config.toml`

```console
$ ec2-vuls-config --config /path/to/config.toml
```

## --filters (-f)

デフォルトの条件に加えて、さらにフィルタリングしたい場合に使用します。フィルタリングは[describe-instances コマンド](http://docs.aws.amazon.com/cli/latest/reference/ec2/describe-instances.html)のように指定できます。
`Name`タグと`Value`タグのセットで指定し、スペース区切りで複数指定可能。

* `web-server`というNameタグのインスタンスをスキャンしたい場合

```console
$ ec2-vuls-config --filters "Name=tag:Name,Values=web-server"
```

* `app-server`というNameタグがついている、かつインスタンスタイプが`c3.large`のインスタンスをスキャンしたい場合

```console
$ ec2-vuls-config --filters "Name=tag:Name,Values=app-server Name=instance-type,Values=r3.large"
```

## --out (-o)

設定ファイルの出力先を指定します。デフォルト: `$PWD/config.toml`

```console
$ ec2-vuls-config --out /path/to/config.toml
```


## --print (-p)

設定ファイルに書き込む代わりに標準出力します。

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
