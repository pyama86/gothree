# gothree
[![Build Status](https://travis-ci.org/pyama86/gothree.svg?branch=master)](https://travis-ci.org/pyama86/gothree)

## Description


## Usage

```bash
Usage of gothree:
  -access-key-id string
        Please specify access key id of aws
  -bucket string
        Please specify bucket of aws s3
  -path string
        Please specify path of aws s3 (default "/")
  -region string
        Please specify region of aws (default "ap-northeast-1")
  -secret-access-key string
        Please specify secret access key of aws
  -version
        Print version information and quit.
```

### gothree is used in combination with logrotate

- logrotate config
```
/var/log/syslog
/var/log/auth
{
        rotate 7
        daily
        missingok
        notifempty
        delaycompress
        compress
        postrotate
          source /root/.aws
          /usr/local/bin/gothree $1
        endscript
}
```

```
$ cat /root/.aws
export AWS_ACCESS_KEY_ID=***********
export AWS_SECRET_ACCESS_KEY=***********
export AWS_REGION=your region
export AWS_BUCKET=your buket name
```

## Install

To install, use `go get`:
```bash
$ go get -d github.com/pyama86/gothree
```

## Download

[We prepare binareis](https://github.com/pyama86/gothree/releases)

## Contribution

1. Fork ([https://github.com/pyama86/gothree/fork](https://github.com/pyama86/gothree/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pyama86](https://github.com/pyama86)
