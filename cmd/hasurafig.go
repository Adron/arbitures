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
		if configExists(configPath) {
			text := getConfigContents(configPath)
			rootHasuraUri := strings.Trim(strings.Trim(hasuraUri, string('"')), "/")
			rootHasuraUri = strings.TrimRight(rootHasuraUri, ":8080")
			text = strings.Replace(text, "http://localhost", rootHasuraUri, 2)
			fmt.Printf("New config text:\n\n%vWriting contents.\n", text)
			renameConfigFile(configPath)
			writeNewConfigFile(configPath, text)
		} else {
			fmt.Printf("The configuration file doesn't appear to exist at %v.\n", configPath)
		}
	},
}

func writeNewConfigFile(configFile string, rewrittenConfiguration string) {
	file, err := os.Create(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := file.WriteString(rewrittenConfiguration)
	if err != nil {
		fmt.Println(err)
		file.Close()
		return
	}
	fmt.Println(l, "Configuration file written successfully.")
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func renameConfigFile(existingConfigFile string) {
	oldConfigFile := existingConfigFile + ".old"
	oldConfigErr := os.Remove(oldConfigFile)
	if oldConfigErr != nil {
		fmt.Println(oldConfigErr)
	}
	err := os.Rename(existingConfigFile, oldConfigFile)
	if err != nil {
		fmt.Println(err)
	}
}

func getConfigContents(configFile string) string {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func configExists(configFile string) bool {
	if _, err := os.Stat(configFile); err == nil {
		return true
	}
	return false
}

func init() {
	mutateCmd.AddCommand(hasurafigCmd)
	hasurafigCmd.Flags().StringVarP(&hasuraUri, "uri", "u", "", "The URI for the API to update the existing URI with in the configuration file.")
	hasurafigCmd.Flags().StringVarP(&configPath, "path", "p", "", "The path of the Hasura configuration file to edit.")
}
