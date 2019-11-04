package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"regexp"
)

var FunctionSourceMapping = map[string]string{"ReceiveVersion": "receive-version.zip", "LookupVersion": "lookup-version.zip"}
var FunctionNameRegexMapping = map[string]*regexp.Regexp{
	"ReceiveVersion": regexp.MustCompile("ReceiveVersion"),
	"LookupVersion":  regexp.MustCompile("LookupVersion"),
}

func LinkupLambdaTargets(actualLambdaFuncs []*string) map[string]string {
	rtn := make(map[string]string, len(FunctionSourceMapping))
	for _, actualLambdaName := range actualLambdaFuncs {
		for key, rx := range FunctionNameRegexMapping {
			if rx.MatchString(*actualLambdaName) {
				rtn[*actualLambdaName] = FunctionSourceMapping[key]
			}
		}
	}
	return rtn
}

func UploadContent(s3Client *s3.S3, fileName string, basePath string, destBucket string, destPath string) (*string, error) {
	fp, openErr := os.Open(basePath + "/" + fileName)
	if openErr != nil {
		log.Printf("Could not open %s/%s: %s", basePath, fileName, openErr)
		return nil, openErr
	}

	outPath := aws.String(destPath + "/" + fileName)
	req := s3.PutObjectInput{
		Body:   fp,
		Bucket: aws.String(destBucket),
		Key:    outPath,
	}

	log.Printf("Uploading %s/%s to s3://%s/%s...", basePath, fileName, destBucket, *outPath)
	_, err := s3Client.PutObject(&req)
	if err != nil {
		log.Printf("Could not upload to s3://%s/%s: %s", destBucket, *outPath, err)
		return nil, err
	}

	log.Printf("Upload completed")
	return outPath, nil
}

func main() {
	var stackName = flag.String("stack", "", "Stack name to query")
	var fromPath = flag.String("from", "lambdas/deployables", "Location containing compiled and zipped deployables")
	flag.Parse()

	if stackName == nil || *stackName == "" {
		println("You must specify a stack name in the --stack argument")
		os.Exit(1)
	}

	mgr := NewLambdaMgr()
	lookupErr := mgr.LookupStack(stackName)

	if lookupErr != nil {
		log.Fatal(lookupErr)
	}

	log.Printf("Got stack %s:\n", mgr.Description())
	lambdas := mgr.LambdaFunctions()

	log.Printf("Found lambdas:\n")
	for _, entry := range lambdas {
		fmt.Printf("\t%s\n", *entry)
	}

	spew.Dump(mgr.LambdaCodeLocations)

	log.Printf("Lambda upload bucket is %s\n", *mgr.DeployablesBucket())

	app, stack, stage := mgr.AppStackStage()
	uploadBasePath := fmt.Sprintf("%s/%s/%s", app, stack, stage)
	log.Printf("Uploads will be to s3://%s/%s/{something}.zip\n", *mgr.DeployablesBucket(), uploadBasePath)

	targetsToUpload := LinkupLambdaTargets(lambdas)
	spew.Dump(targetsToUpload)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	s3Client := s3.New(sess)

	for funcName, deployableZipName := range targetsToUpload {
		uploadedPath, err := UploadContent(s3Client, deployableZipName, *fromPath, *mgr.DeployablesBucket(), uploadBasePath)
		if err == nil {
			updateErr := mgr.UpdateLambdaCode(&funcName, mgr.DeployablesBucket(), uploadedPath)
			if updateErr == nil {
				log.Printf("Successfully updated %s", funcName)
			} else {
				log.Printf("Could not update %s to code from %s: %s", funcName, uploadedPath, updateErr)
			}
		}
	}
}
