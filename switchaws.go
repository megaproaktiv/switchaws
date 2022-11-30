package switchaws

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	Regionfound bool
)


func ExpandTildeToUserHome(filePath string) string {
	if strings.HasPrefix(filePath, "~/") {
		return filepath.Join(os.Getenv("HOME"), filePath[2:])
	}
	return filePath
}


func ConvertCredentialsEntry(credFileVar string) (string, error) {
	switch credFileVar {
	case "region":
		if !Regionfound {
			Regionfound = true
			return "AWS_DEFAULT_REGION", nil
		} else {
			return "", fmt.Errorf("region already found")
		}
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
		return "", fmt.Errorf("unknown credentials file variable: %s", credFileVar)
	}
}

func ConvertAWSConfigEntry(credFileVar string) (string, error) {
	switch credFileVar {
	case "region":
		if !Regionfound {
			Regionfound = true
			return "AWS_DEFAULT_REGION", nil
		} else {
			return "", fmt.Errorf("region already found")
		}
	case "output":
		return "AWS_DEFAULT_OUTPUT", nil
	default:
		return "", fmt.Errorf("unknown credentials file variable: %s", credFileVar)
	}
}

func ConvertAdditionalConfigEntry(credFileVar string) (string, error) {
	switch credFileVar {
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
	default:
		return "", fmt.Errorf("unknown credentials file variable: %s", credFileVar)
	}
}
