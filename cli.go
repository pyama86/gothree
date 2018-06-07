package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
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
		awsKey string
		awsID  string
		region string
		bucket string
		path   string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&awsKey, "secret-access-key", "", "Please specify secret access key of aws")
	flags.StringVar(&awsID, "access-key-id", "", "Please specify access key id of aws")
	flags.StringVar(&region, "region", "ap-northeast-1", "Please specify region of aws")
	flags.StringVar(&bucket, "bucket", "", "Please specify bucket of aws s3")
	flags.StringVar(&path, "path", "", "Please specify path of aws s3")
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

	s, err := NewSthree(awsID, awsKey, region, bucket)

	if err != nil {
		logrus.Fatal(err)
	}

	for _, n := range flags.Args() {
		s.Put(n)
	}
	return ExitCodeOK
}

type sthree struct {
	SecretAccessKey string `validate:"required"`
	AccessKeyID     string `validate:"required"`
	Region          string `validate:"required"`
	Bucket          string `validate:"required"`
}

func assignEnv(v *string, key string) {
	if *v == "" && os.Getenv(key) != "" {
		*v = os.Getenv(key)
	}
}

func NewSthree(id, key, region, bucket string) (*sthree, error) {
	assignEnv(&id, "AWS_ACCESS_KEY_ID")
	assignEnv(&key, "AWS_SECRET_ACCESS_KEY")
	assignEnv(&region, "AWS_REGION")

	s := &sthree{
		AccessKeyID:     id,
		SecretAccessKey: key,
		Region:          region,
		Bucket:          bucket,
	}
	config := &validator.Config{TagName: "validate"}
	validate := validator.New(config)
	err := validate.Struct(s)
	if err != nil {
		return nil, err
	}
	return s, nil

}

func replaceExt(filePath, from, to string) string {
	rep := regexp.MustCompile(`\.` + from + `$`)
	return rep.ReplaceAllString(filePath, "."+to)
}

func saveName(filePath string) string {
	t := time.Now().Local()
	today := t.Format("20060102")

	name := replaceExt(filePath, "1.gz", today+".gz")
	return replaceExt(name, "1", today+".gz")
}

func (s *sthree) Put(filePath string) error {
	var reader io.Reader
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
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(saveName(filePath)),
		Body:   reader,
	})

	if err != nil {
		return err
	}
	return nil
}
