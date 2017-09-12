package cmd

import (
	"fmt"
	"strings"

	"github.com/akerl/madlibrarian/utils"

	"github.com/spf13/cobra"
)

func generateRunner(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("No config path provided. See --help for more info")
	}

	path := args[0]
	var s utils.Story
	var err error

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		s, err = utils.NewStoryFromURL(path)
	} else {
		s, err = utils.NewStoryFromFile(path)
	}
	if err != nil {
		return err
	}

	quote, err := s.Generate()
	if err != nil {
		return err
	}
	fmt.Println(quote)
	return nil
}

var generateCmd = &cobra.Command{
	Use:   "generate PATH",
	Short: "generate a quote",
	RunE:  generateRunner,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
