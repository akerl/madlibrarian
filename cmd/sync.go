package cmd

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type metaData struct {
	Template string
	Bucket   string
	Path     string
}

type quoteList []string

type chunkList map[string]quoteList

type configFile struct {
	Meta   metaData
	Chunks chunkList
}

func getIndex(max int) (int, error) {
	maxBig := int64(max)
	n, err := rand.Int(rand.Reader, big.NewInt(maxBig))
	if err != nil {
		return 0, err
	}
	i := int(n.Int64())
	return i - 1, nil
}

func randomizer(ql quoteList) func() (string, error) {
	return func() (string, error) {
		index, err := getIndex(len(ql))
		if err != nil {
			return "", err
		}
		return ql[index], nil
	}
}

func syncRunner(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("No file provided. See --help for more info")
	}

	filePath := args[0]

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	cf := configFile{}
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return err
	}

	funcMap := template.FuncMap{}
	for chunkName, ql := range cf.Chunks {
		funcMap[chunkName] = randomizer(ql)
	}

	tmpl, err := template.New(cf.Meta.Path).Funcs(funcMap).Parse(cf.Meta.Template)
	if err != nil {
		return err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, "")
	if err != nil {
		return err
	}

	fmt.Println(result.String())

	return nil
}

var syncCmd = &cobra.Command{
	Use:   "sync FILE",
	Short: "Sync quotes file to S3",
	RunE:  syncRunner,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
