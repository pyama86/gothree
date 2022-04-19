package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"

	validator "gopkg.in/go-playground/validator.v8"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		awsKey      string
		awsID       string
		region      string
		bucket      string
		path        string
		concurrency int

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&awsKey, "secret-access-key", "", "Please specify secret access key of aws")
	flags.StringVar(&awsID, "access-key-id", "", "Please specify access key id of aws")
	flags.StringVar(&region, "region", "ap-northeast-1", "Please specify region of aws")
	flags.StringVar(&bucket, "bucket", "", "Please specify bucket of aws s3")
	flags.StringVar(&path, "path", "/", "Please specify path of aws s3")
	flags.IntVar(&concurrency, "concurrency", 5, "Upload concurrency")
	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	s, err := newSthree(awsID, awsKey, region, bucket, path, concurrency)

	if err != nil {
		logrus.Fatal(err)
	}

	for _, n := range flags.Args() {
		if err := s.Put(n); err != nil {
			logrus.Error(err)
		}

	}
	return ExitCodeOK
}

type sthree struct {
	SecretAccessKey string `validate:"required"`
	AccessKeyID     string `validate:"required"`
	Region          string `validate:"required"`
	Bucket          string `validate:"required"`
	Path            string `validate:"required"`
	Concurrency     int    `validate:"required"`
}

func assignEnv(v *string, key string) {
	if *v == "" && os.Getenv(key) != "" {
		*v = os.Getenv(key)
	}
}

func newSthree(id, key, region, bucket, path string, concurrency int) (*sthree, error) {
	assignEnv(&id, "AWS_ACCESS_KEY_ID")
	assignEnv(&key, "AWS_SECRET_ACCESS_KEY")
	assignEnv(&region, "AWS_REGION")
	assignEnv(&bucket, "AWS_BUCKET")

	s := &sthree{
		AccessKeyID:     id,
		SecretAccessKey: key,
		Region:          region,
		Bucket:          bucket,
		Path:            path,
		Concurrency:     concurrency,
	}
	config := &validator.Config{TagName: "validate"}
	validate := validator.New(config)
	err := validate.Struct(s)
	if err != nil {
		return nil, err
	}
	return s, nil

}

func today() string {
	t := time.Now().Local()
	return t.Format("20060102")
}

func saveName(filePath string) string {
	rep := regexp.MustCompile(fmt.Sprintf(`(\.1$|\.gz$|[-\.]%s)`, today()))
	filePath = rep.ReplaceAllString(filePath, "")

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Error(err)
		return filepath.Base(fmt.Sprintf("%s.%s.gz", filePath, today()))
	}

	if os.Getenv("GOTHREE_PREFIX") != "" {
		hostname = os.Getenv("GOTHREE_PREFIX")
	}
	return fmt.Sprintf("%s-%s", hostname, filepath.Base(fmt.Sprintf("%s.%s.gz", filePath, today())))
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func lotateFileName(filePath string) string {
	patterns := []string{
		fmt.Sprintf("%s.1", filePath),
		fmt.Sprintf("%s-%s", filePath, today()),
		fmt.Sprintf("%s.%s", filePath, today()),
	}

	for _, p := range patterns {
		if exists(p) {
			return p
		}
	}
	return filePath
}

func (s *sthree) Put(filePath string) error {
	var reader io.Reader

	filePath = lotateFileName(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if !strings.HasSuffix(filePath, ".gz") {
		r, w := io.Pipe()
		reader = r
		go func() {
			gw := gzip.NewWriter(w)
			io.Copy(gw, file)
			file.Close()
			gw.Close()
			w.Close()
		}()

	} else {
		reader = file
	}

	st, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s.AccessKeyID, s.SecretAccessKey, ""),
		Region:      &s.Region,
	})

	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(st)
	uploader.Concurrency = s.Concurrency
	logrus.Infof("start upload file:%s save: %s", filePath, saveName(filePath))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filepath.Join(s.Path, saveName(filePath))),
		Body:   reader,
	})

	if err != nil {
		return err
	}
	logrus.Infof("upload success: %s", result.Location)
	return nil
}
