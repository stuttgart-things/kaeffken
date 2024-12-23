/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package surveys

import (
	"fmt"
	"log"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/stuttgart-things/kaeffken/modules"
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
)

// CONFIG REPRESENTS THE STRUCTURE OF THE YAML FILE
type Config struct {
	GitRepo    string   `yaml:"gitRepo"`
	GitOwner   string   `yaml:"gitOwner"`
	GitBranch  string   `yaml:"gitBranch"`
	RootFolder string   `yaml:"rootFolder"`
	SubFolder  string   `yaml:"subFolder"`
	Questions  []string `yaml:"questions"`
	Templates  []string `yaml:"templates"`
	Technology string   `yaml:"technology"`
}

func RunGitHubBranchingFlow(config Config, projectName string) map[string]interface{} {

	configMap := ConfigToMap(config, projectName)

	branchDefault, err := sthingsBase.RenderTemplateInline(config.GitBranch, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, configMap)
	if err != nil {
		fmt.Println(err)
	}

	githubBranchQuestions := map[string]modules.InputQuestion{
		"GITHUB REPO?": {
			Question:  "GITHUB REPO?",
			Default:   config.GitRepo,
			MinLength: 3,
			MaxLength: 18,
			Id:        "girepo",
			Type:      "string",
		},
		"GITHUB OWNER?": {
			Question:  "GITHUB OWNER?",
			Default:   config.GitOwner,
			MinLength: 3,
			MaxLength: 18,
			Id:        "gitowner",
			Type:      "string",
		},
		"BRANCH?": {
			Question:  "BRANCH NAME?",
			Default:   string(branchDefault),
			MinLength: 3,
			MaxLength: 34,
			Id:        "branchName",
			Type:      "string",
		},
	}

	githubPRAnswers, err := modules.AskInputQuestions(githubBranchQuestions)
	if err != nil {
		log.Fatalf("ERROR ASKING META QUESTIONS: %v", err)
	}

	return githubPRAnswers
}

// ConfigToMap converts a Config struct to a map[string]interface{}
func ConfigToMap(cfg Config, projectName string) map[string]interface{} {
	return map[string]interface{}{
		"gitRepo":     cfg.GitRepo,
		"gitOwner":    cfg.GitOwner,
		"gitBranch":   cfg.GitBranch,
		"rootFolder":  cfg.RootFolder,
		"subFolder":   cfg.SubFolder,
		"technology":  cfg.Technology,
		"projectName": projectName,
	}
}

func ReadGitProfile(filename string) (config Config) {

	if err := modules.ReadYAML(filename, &config); err != nil {
		fmt.Printf("ERROR READING YAML FILE: %v\n", err)
	}

	return
}
