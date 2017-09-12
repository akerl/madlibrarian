package cmd

import (
	"fmt"

	"github.com/akerl/madlibrarian/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func uploadRunner(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("Invalid arguments provided. See --help for more info")
	}

	path := args[0]
	bucket := args[1]
	prefix := args[2] // TODO: Add some kind of UUID here

	s, err := utils.NewStoryFromPath(path)
	if err != nil {
		return err
	}

	_, err = s.Generate()
	if err != nil {
		return err
	}

	if s.Meta.Type != "local" {
		return fmt.Errorf("Upload only makes sense for local stories")
	}

	var chunks []string
	funcMap, err := s.TypeObj.Funcs(&s)
	if err != nil {
		return err
	}
	for name := range funcMap {
		chunks = append(chunks, name)
	}

	newStory := utils.Story{
		Meta: utils.Metadata{
			Type:     "s3",
			Template: s.Meta.Template,
		},
		Data: map[string]interface{}{
			"s3": map[string]string{
				"bucket": bucket,
				"prefix": prefix,
			},
			"chunks": chunks,
		},
	}

	err = utils.Upload(s, bucket, prefix)
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
