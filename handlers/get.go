package handlers

import (
		"fmt"
	"os"
	"strings"
		"gopkg.in/alecthomas/kingpin.v2"
)

type GetHandler struct {
	SubCommand *kingpin.CmdClause
	Arguments  GetHandlerArguments
}

type GetHandlerArguments struct {
	CredentialsFilePath   *string
	ConfigFilePath   *string
}

func NewGetHandler(app *kingpin.Application) GetHandler {
	subCommand := app.Command("get", "get current AWS profile (that is set to default profile)")

	credentialsFilePath := subCommand.Flag("credentials-path", "Path to AWS Credentials file").Default("~/.aws/credentials").String()
	configFilePath := subCommand.Flag("config-path", "Path to AWS Config file").Default("~/.aws/config").String()

	return GetHandler{
		SubCommand: subCommand,
		Arguments:   GetHandlerArguments{
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath: configFilePath,
		},
	}
}

func (handler GetHandler) Handle() {
	configFile, err := ReadFile(*handler.Arguments.ConfigFilePath)
	if err != nil {
		fmt.Printf("Fail to read AWS config file: %v", err)
		os.Exit(1)
	}

	configDefaultSection, err := configFile.GetSection("default")
	if err == nil &&
		configDefaultSection.HasKey("role_arn") &&
		configDefaultSection.HasKey("source_profile") {

		defaultRoleArn := configDefaultSection.Key("role_arn").Value()
		defaultSourceProfile := configDefaultSection.Key("source_profile").Value()

		for _, section := range configFile.Sections() {
			if strings.Compare(section.Name(), "default") != 0 &&
				section.Haskey("role_arn") &&
				section.HasKey("source_profile") &&
				strings.Compare(section.Key("role_arn").Value(), defaultRoleArn) == 0 &&
				strings.Compare(section.Key("source_profile").Value(), defaultSourceProfile) == 0 {
				fmt.Printf("assuming %s\n", section.Name())
				os.Exit(0)
			}
		}
	}

	credentialsFile, err := ReadFile(*handler.Arguments.CredentialsFilePath)
	if err != nil {
		fmt.Printf("Fail to read AWS credentials file: %v", err)
		os.Exit(1)
	}

	credentialsDefaultSection, err := credentialsFile.GetSection("default")
	if err == nil &&
		credentialsDefaultSection.HasKey("aws_access_key_id") {

		defaultAWSAccessKeyId := credentialsDefaultSection.Key("aws_access_key_id").Value()

		for _, section := range credentialsFile.Sections() {
			if strings.Compare(section.Name(), "default") != 0 &&
				section.HasKey("aws_access_key_id") &&
				strings.Compare(section.Key("aws_access_key_id").Value(), defaultAWSAccessKeyId) == 0 {
				fmt.Printf("%s\n", section.Name())
				os.Exit(0)
			}
		}
	}
}