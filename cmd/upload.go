package cmd

import (
	"fmt"

	"github.com/akerl/madlibrarian/utils"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func uploadRunner(_ *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("invalid arguments provided. See --help for more info")
	}

	path := args[0]
	bucket := args[1]
	prefix := args[2]

	s, err := utils.NewStoryFromPath(path)
	if err != nil {
		return err
	}

	newStory, err := s.Upload(bucket, prefix)
	if err != nil {
		return err
	}

	output, err := yaml.Marshal(newStory)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(output))

	return nil
}

var uploadCmd = &cobra.Command{
	Use:   "upload PATH BUCKET PREFIX",
	Short: "upload a quote",
	RunE:  uploadRunner,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
