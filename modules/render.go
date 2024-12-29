/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package modules

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/getsops/sops/v3/decrypt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

type TemplateBracket struct {
	begin        string `mapstructure:"begin"`
	end          string `mapstructure:"end"`
	regexPattern string `mapstructure:"regex-pattern"`
}

var (
	renderOption = "missingkey=error"
	brackets     = map[string]TemplateBracket{
		"curly":  TemplateBracket{"{{", "}}", `\{\{(.*?)\}\}`},
		"square": TemplateBracket{"[[", "]]", `\[\[(.*?)\]\]`},
	}
	bracketFormat = "curly"
	fileFormat    = "yaml"
)

func RenderTemplate[T any](tmpl string, data T) (string, error) {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RenderAliases(aliases []string, allValues map[string]interface{}) map[string]interface{} {

	fmt.Println("ALL VALUES: ", allValues)

	for _, alias := range aliases {

		// SPLIT ALIAS KEY/VALUE BY :
		aliasValues := strings.Split(alias, ":")

		// RENDER KEY
		aliasKey, err := sthingsBase.RenderTemplateInline(aliasValues[0], renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
		if err != nil {
			fmt.Println(err)
		}

		// RENDER VALUE
		aliasValue, err := sthingsBase.RenderTemplateInline(aliasValues[1], renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
		if err != nil {
			fmt.Println(err)
		}

		// ASSIGN ALIAS TO ALL VALUES
		key := string(strings.TrimSpace(string(aliasKey)))
		value := string(strings.TrimSpace(string(aliasValue)))

		allValues[string(key)] = string(value)
		log.Info("ALIAS ADDED: ", key, ":", string(value))
	}

	return allValues
}

func GetAllSecretsFromSopsDecyptedFiles(secretFiles []string, allValues map[string]interface{}) (allSecretsFromAllSecretsFile map[string]interface{}) {

	allSecretsFromAllSecretsFile = make(map[string]interface{})

	// GET ALL SECRETS TO A VALUES MAP
	for _, secretFile := range secretFiles {

		// SPLIT SECRET FILE PATH BY :
		secretFilePaths := strings.Split(secretFile, ":")

		// RENDER SOURCE AND TARGET NAMES
		sourcePath, err := sthingsBase.RenderTemplateInline(secretFilePaths[0], renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
		if err != nil {
			fmt.Println(err)
		}

		// DECRYPT SOURCE FILE
		decryptedFile, err := decrypt.File(string(sourcePath), fileFormat)
		if err != nil {
			log.Error("FAILED TO DECRYPT: ", err)
		}

		// EXTRACT VALUES
		allSecretsFromAllSecretsFile = CreateSecretsMap(decryptedFile, nil)
	}

	return
}
