package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

var (
	credentialsFile string
	configFile string
	verbose         bool
	help            bool
)

func init() {
	flag.StringVar(&credentialsFile, "f", "~/.aws/credentials", "Path to AWS credentials file")
	flag.StringVar(&configFile, "c", "~/.aws/config", "Path to AWS credentials file")
	flag.BoolVar(&verbose, "v", false, "Verbose output for debugging")
	flag.BoolVar(&help, "h", false, "Print command usage")
	flag.Parse()

	if help {
		usage := `aws-profiles is a tool to manage multiple AWS profiles using the credentials file

Usage: 

aws-profiles [-f filepath] [-v] [-h] profile-name

  -f	Override the default credentials file location (~/.aws/credentials)
  -v	Turn on verbose logging for debugging
  -h	Print this message
`
		fmt.Println(usage)
		os.Exit(0)
	}

	// disable logging if verbose is false
	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
}


func main() {
	requestedProfile := os.Args[1]
	//log.Printf("Args1: %v", requestedProfile)
	//log.Printf("Credentialsfile: %s", expandTildeToUserHome(credentialsFile))
	cfgCredentials, err := ini.Load(expandTildeToUserHome(credentialsFile))
	if err != nil {
		fmt.Printf("Unable to read file: %v\n", err)
	}
	cfgConfig, err := ini.Load(expandTildeToUserHome(configFile))

	if err != nil {
		fmt.Printf("Unable to read file: %v\n", err)
	}

	for _, section := range cfgCredentials.Sections() {
		if section.Name() == requestedProfile {
			keyHash := section.KeysHash()
			for k, v := range keyHash {
				if k == "__name__" {
					continue
				}
				canonicalForm, err := convertCredentialsEntry(k)
				if err != nil {
					continue
				}
				fmt.Printf("export %s=%s\n", canonicalForm, v)
			}
		}
	}
	
	for _, section := range cfgConfig.Sections() {
		profile := "profile "+requestedProfile
		sectionName :=strings.TrimSpace(section.Name())
		if  strings.Compare(profile,sectionName) == 0 {
			keyHash := section.KeysHash()
			for k, v := range keyHash {
				canonicalForm, err := convertConfigEntry(k)
				if err != nil {
					continue
				}
				// Only Change dir if not in subdir
				if canonicalForm == "CHDIR" {
					currentWorkingDirectory, err := os.Getwd()
					if err != nil {
						log.Fatal(err)
					}
					if strings.HasPrefix(currentWorkingDirectory, v){

					}else{
						fmt.Printf("export %s=%s\n", canonicalForm, v)
					}
				}else{ 
					fmt.Printf("export %s=%s\n", canonicalForm, v)
				}
			}
		}
	}
	fmt.Printf("export %s=%s\n", "AWSUME_PROFILE", requestedProfile)
	fmt.Printf("export %s=%s\n", "AWS_DEFAULT_PROFILE", requestedProfile)
	// Treat mfa/awsumerole
	// if source_profile is there
	// get key/secret from the source profile
	// get arn from the source profile
	// ask user (in bash?)
	// aws sts get-session-token --serial-number arn-of-the-mfa-device --token-code code-from-token
	// export stuff
	// "SecretAccessKey": "secret-access-key",
	// "SessionToken": "temporary-session-token",
	// "Expiration": "expiration-date-time",
	// "AccessKeyId": "access-key-id"
}

func expandTildeToUserHome(filePath string) string {
	if strings.HasPrefix(filePath, "~/") {
		return filepath.Join(os.Getenv("HOME"), filePath[2:])
	}
	return filePath
}

func convertCredentialsEntry(credFileVar string) (string, error) {
	switch credFileVar {
	case "region":
		return "AWS_DEFAULT_REGION", nil
	case "output":
		return "AWS_DEFAULT_OUTPUT", nil
	case "aws_access_key_id":
		return "AWS_ACCESS_KEY_ID", nil
	case "aws_secret_access_key":
		return "AWS_SECRET_ACCESS_KEY", nil
	case "aws_sts_token":
		return "AWS_STS_TOKEN", nil
	case "aws_session_token":
		return "AWS_SESSION_TOKEN", nil
	default:
		return "", fmt.Errorf("Unknown credentials file variable: %s", credFileVar)
	}
}

func convertConfigEntry(credFileVar string) (string, error) {
	switch credFileVar {
	case "region":
		return "AWS_DEFAULT_REGION", nil
	case "workdir":
		return "CHDIR", nil
	case "itermbadge":
		return "ITERMBADGE", nil
	case "taskwarrior":
		return "TASKWARRIOR", nil
	case "source_profile":
		return "source_profile", nil
	case "url":
		return "AWS_URL", nil
	case "output":
		return "AWS_DEFAULT_OUTPUT", nil
	default:
		return "", fmt.Errorf("Unknown credentials file variable: %s", credFileVar)
	}
}
