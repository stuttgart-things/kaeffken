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
	GitRepo       string   `yaml:"gitRepo"`
	GitOwner      string   `yaml:"gitOwner"`
	GitBranch     string   `yaml:"gitBranch"`
	CommitMessage string   `yaml:"commitMessage"`
	RootFolder    string   `yaml:"rootFolder"`
	SubFolder     string   `yaml:"subFolder"`
	Questions     []string `yaml:"questions"`
	Templates     []string `yaml:"templates"`
	Technology    string   `yaml:"technology"`
	PrTitle       string   `yaml:"prTitle"`
	PrDescription string   `yaml:"prDescription"`
	PrTags        []string `yaml:"prTags"`
	Aliases       []string `yaml:"aliases"`
}

func RunGitHubBranchingFlow(config Config, values map[string]interface{}) map[string]interface{} {

	defaults := ConfigToMap(config)
	// MERGE WITH VALUES
	allValues := sthingsBase.MergeMaps(defaults, values)

	branchDefault, err := sthingsBase.RenderTemplateInline(config.GitBranch, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
	if err != nil {
		fmt.Println(err)
	}

	commitMessageDefault, err := sthingsBase.RenderTemplateInline(config.CommitMessage, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
	if err != nil {
		fmt.Println(err)
	}

	PrTitleDefault, err := sthingsBase.RenderTemplateInline(config.CommitMessage, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
	if err != nil {
		fmt.Println(err)
	}

	PrDescriptionDefault, err := sthingsBase.RenderTemplateInline(config.CommitMessage, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
	if err != nil {
		fmt.Println(err)
	}

	githubBranchQuestions := map[string]modules.InputQuestion{
		"GITHUB REPO?": {
			Question:  "GITHUB REPO?",
			Default:   config.GitRepo,
			MinLength: 3,
			MaxLength: 18,
			Id:        "gitRepo",
			Type:      "string",
		},
		"GITHUB OWNER?": {
			Question:  "GITHUB OWNER?",
			Default:   config.GitOwner,
			MinLength: 3,
			MaxLength: 18,
			Id:        "gitOwner",
			Type:      "string",
		},
		"BRANCH?": {
			Question:  "BRANCH NAME?",
			Default:   string(branchDefault),
			MinLength: 3,
			MaxLength: 34,
			Id:        "gitBranch",
			Type:      "string",
		},
		"COMMIT MESSAGE?": {
			Question:  "COMMIT MESSAGE?",
			Default:   string(commitMessageDefault),
			MinLength: 3,
			MaxLength: 34,
			Id:        "commitMessage",
			Type:      "string",
		},
		"PULL REQUEST TITLE?": {
			Question:  "PULL REQUEST TITLE?",
			Default:   string(PrTitleDefault),
			MinLength: 3,
			MaxLength: 34,
			Id:        "prTitle",
			Type:      "string",
		},
		"PULL REQUEST DESCRIPTION?": {
			Question:  "PULL REQUEST DESCRIPTION?",
			Default:   string(PrDescriptionDefault),
			MinLength: 3,
			MaxLength: 34,
			Id:        "prDescription",
			Type:      "string",
		},
	}

	githubPRAnswers, err := modules.AskInputQuestions(githubBranchQuestions)
	if err != nil {
		log.Fatalf("ERROR ASKING META QUESTIONS: %v", err)
	}

	return githubPRAnswers
}

// CONFIGTOMAP CONVERTS A CONFIG STRUCT TO A MAP[STRING]INTERFACE{}
func ConfigToMap(cfg Config) map[string]interface{} {
	return map[string]interface{}{
		"gitRepo":       cfg.GitRepo,
		"gitOwner":      cfg.GitOwner,
		"gitBranch":     cfg.GitBranch,
		"commitMessage": cfg.CommitMessage,
		"rootFolder":    cfg.RootFolder,
		"subFolder":     cfg.SubFolder,
		"technology":    cfg.Technology,
		"prTitle":       cfg.PrTitle,
		"prDescription": cfg.PrDescription,
		"prTags":        cfg.PrTags,
		"aliases":       cfg.Aliases,
	}
}

func ReadGitProfile(filename string) (config Config) {

	if err := modules.ReadYAML(filename, &config); err != nil {
		fmt.Printf("ERROR READING YAML FILE: %v\n", err)
	}

	return
}
