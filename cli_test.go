package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./gothree -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("gothree version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

func TestRun_awsKeyFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./gothree -aws-key", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_awsIdFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./gothree -aws-id", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_regionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./gothree -region", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_bucketFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./gothree -bucket", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_pathFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./gothree -path", " ")

	status := cli.Run(args)
	_ = status
}

func Test_replaceExt(t *testing.T) {
	type args struct {
		filePath string
		from     string
		to       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "last 1 day",
			args: args{
				filePath: "/path/to/example.1",
				from:     "1",
				to:       "20180101",
			},
			want: "/path/to/example.20180101",
		},
		{
			name: "gzip",
			args: args{
				filePath: "/path/to/example.20180101",
				from:     "20180101",
				to:       "20180101.gz",
			},
			want: "/path/to/example.20180101.gz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceExt(tt.args.filePath, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("replaceExt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveName(t *testing.T) {
	n := time.Now().Local()
	today := n.Format("20060102")
	want := "/path/to/example." + today + ".gz"

	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "last 1 day",
			args: args{
				filePath: "/path/to/example.1",
			},
			want: want,
		},
		{
			name: "last 1 day gzip",
			args: args{
				filePath: "/path/to/example.1.gz",
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := saveName(tt.args.filePath); got != tt.want {
				t.Errorf("saveName() = %v, want %v", got, tt.want)
			}
		})
	}
}
