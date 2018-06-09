package main

import (
	"bytes"
	"fmt"
	"os"
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

func Test_saveName(t *testing.T) {
	n := time.Now().Local()
	os.Setenv("GOTHREE_PREFIX", "example")
	today := n.Format("20060102")
	want := "example-example." + today + ".gz"

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
			name: "date-ext",
			args: args{
				filePath: fmt.Sprintf("/path/to/example.%s", today),
			},
			want: want,
		},
		{
			name: "gz",
			args: args{
				filePath: "/path/to/example.gz",
			},
			want: want,
		},
		{
			name: "date-ext.gz",
			args: args{
				filePath: fmt.Sprintf("/path/to/example.%s.gz", today),
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

func Test_lotateFileName(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lotateFileName(tt.args.filePath); got != tt.want {
				t.Errorf("lotateFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
