package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var hasuraUri string
var configPath string

var hasurafigCmd = &cobra.Command{
	Use:   "hasurafig",
	Short: "This command will take a URI and update the existing URI in the Hasura Config file.",
	Long: `This command takes the passed in URI and replaces the existing 'localhost' URI listed in the
Hasura config file that is located at the passed in path. Note, if it isn't set to the default localhost
path and has already been updated, this command will not process the edit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Editing the Hasura configuration file located at %v with \nthe new URI of %v.\n...\n\n",
			configPath, hasuraUri)

		if configExists(){
			text := getConfigContents()
			rootHasuraUri := strings.Trim(strings.Trim(hasuraUri, string('"')), "/")
			rootHasuraUri = strings.TrimRight(rootHasuraUri, ":8080")
			text = strings.Replace(text, "http://localhost", rootHasuraUri, 2)
			fmt.Printf("New config text:\n\n%v", text)

		} else {
			fmt.Printf("The configuration file doesn't appear to exist at %v.\n", configPath)
		}
	},
}

func getConfigContents() string {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	return string(content)
}

func configExists() bool {
	if _, err := os.Stat(configPath); err == nil {
		return true
	}
	return false
}

func init() {
	mutateCmd.AddCommand(hasurafigCmd)
	hasurafigCmd.Flags().StringVarP(&hasuraUri, "uri", "u", "", "The URI for the API to update the existing URI with in the configuration file.")
	hasurafigCmd.Flags().StringVarP(&configPath, "path", "p", "", "The path of the Hasura configuration file to edit.")
}
