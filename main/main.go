package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"switchaws"
	"github.com/go-ini/ini"
)

var (
	credentialsFile string
	configFile      string
	profile         string
	start           bool
	verbose         bool
	help            bool
)


func init() {

	flag.StringVar(&credentialsFile, "f", "~/.aws/credentials", "Path to AWS credentials file")
	flag.StringVar(&profile, "p", "", "profile to switch")
	flag.StringVar(&configFile, "c", "~/.aws/config", "Path to AWS credentials file")
	flag.BoolVar(&start, "s", false, "switch on start behaviour")
	flag.BoolVar(&verbose, "v", false, "Verbose output for debugging")
	flag.BoolVar(&help, "h", false, "Print command usage")
	flag.Parse()

	if help {
		usage := `switchaws is a tool to manage multiple AWS profiles using the credentials file

Usage: 

aws-profiles [-f filepath] [-v] [-h] -p profile-name
  -s    Turn on start behaviour defined in ~/.aws/config: 
    - Description: configname / environment variable
    - Change to project dir: chdir / CHDIR
    - Set iterm badge: itermbadge / ITERMBADGE
    - switch taskwarrior task: taskwarrior / TASKWARRIOR
    - open url: url / AWS_URL
    The environment variables are handled in switchawswrapper.sh
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
	requestedProfile := profile

	// if called with no parameter, just switch
	firstArg := os.Args[1]
	if !strings.HasPrefix(firstArg, "-") {
		requestedProfile = firstArg
	}

	//log.Printf("Args1: %v", requestedProfile)
	//log.Printf("Credentialsfile: %s", expandTildeToUserHome(credentialsFile))

	cfgCredentials, err := ini.Load(switchaws.ExpandTildeToUserHome(credentialsFile))
	if err != nil {
		fmt.Printf("Unable to read file: %v\n", err)
	}
	cfgConfig, err := ini.Load(switchaws.ExpandTildeToUserHome(configFile))

	if err != nil {
		fmt.Printf("Unable to read file: %v\n", err)
	}
	switchaws.Regionfound = false
	configFound := false
	credentialsFound := false
	for _, section := range cfgCredentials.Sections() {
		if section.Name() == requestedProfile {
			credentialsFound = true
			keyHash := section.KeysHash()
			for k, v := range keyHash {
				if k == "__name__" {
					continue
				}
				canonicalForm, err := switchaws.ConvertCredentialsEntry(k)
				if err != nil {
					continue
				}
				fmt.Printf("export %s=%s\n", canonicalForm, v)
			}
		}
	}

	for _, section := range cfgConfig.Sections() {
		profile := "profile " + requestedProfile
		sectionName := strings.TrimSpace(section.Name())
		if strings.Compare(profile, sectionName) == 0 {
			configFound = true
			keyHash := section.KeysHash()
			for k, v := range keyHash {

				canonicalForm, err := switchaws.ConvertAWSConfigEntry(k)
				if err == nil {
					fmt.Printf("export %s=%s\n", canonicalForm, v)
				}

				if start {
					canonicalForm, err = switchaws.ConvertAdditionalConfigEntry(k)
					if err == nil {
						if canonicalForm == "CHDIR" {
							currentWorkingDirectory, err := os.Getwd()
							if err != nil {
								log.Fatal(err)
							}
							if !strings.HasPrefix(currentWorkingDirectory, v) {
								fmt.Printf("export %s=%s\n", canonicalForm, v)
							}
						}

						if canonicalForm != "CHDIR" {
							fmt.Printf("export %s=%s\n", canonicalForm, v)
						}
					}

				}

			}
		}
	}
	fmt.Printf("export %s=%s\n", "AWSUME_PROFILE", requestedProfile)
	fmt.Printf("export %s=%s\n", "AWS_DEFAULT_PROFILE", requestedProfile)
	if !configFound && !credentialsFound {
		fmt.Printf("export %s=%s\n", "MESSAGE", "error_no_credentials_found")
	}
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


