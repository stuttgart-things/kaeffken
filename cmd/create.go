/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stuttgart-things/kaeffken/modules"
	"github.com/stuttgart-things/kaeffken/surveys"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
)

type TemplateBracket struct {
	begin        string `mapstructure:"begin"`
	end          string `mapstructure:"end"`
	regexPattern string `mapstructure:"regex-pattern"`
}

type TemplateData struct {
	TemplateFileName string
	TemplateContent  string
	OutputFileName   string
	RenderedContent  string
}

var (
	allTemplateData []TemplateData
	allQuestions    []*modules.Question
	bracketFormat   = "curly"
	allValues       = make(map[string]interface{})
	renderOption    = "missingkey=error"
	brackets        = map[string]TemplateBracket{
		"curly":  TemplateBracket{"{{", "}}", `\{\{(.*?)\}\}`},
		"square": TemplateBracket{"[[", "]]", `\[\[(.*?)\]\]`},
	}
	metaQuestions = map[string]modules.InputQuestion{
		"Project name?": {
			Question:  "Project name?",
			Default:   "",
			MinLength: 3,
			MaxLength: 18,
			Id:        "projectName",
			Type:      "string",
		},
	}

	metaAnswers = map[string]interface{}{
		"projectName": "",
	}

	commitQuestions = map[string]map[string]interface{}{
		"commit": {
			"name":    "Commit rendered files to branch?",
			"options": []string{"true", "false"},
		},
		"pr": {
			"name":    "Create a pull request?",
			"options": []string{"true", "false"},
		},
	}
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create things",
	Long:  `Create things like rendered resource definitions for storing in gitops repositories.`,
	Run: func(cmd *cobra.Command, args []string) {

		// GET VARS
		runSurvey, _ := cmd.LocalFlags().GetBool("survey")
		profile, _ := cmd.LocalFlags().GetString("profile")
		projectName, _ := cmd.LocalFlags().GetString("project")
		tmpDir, _ := cmd.LocalFlags().GetString("tmp")
		outputDir, _ := cmd.LocalFlags().GetString("output")

		if outputDir == "" {
			outputDir = tmpDir + "/" + time.Now().Format("20060102_150405")
		}

		gitConfig := surveys.ReadGitProfile(profile)
		fmt.Println(gitConfig.GitBranch)
		fmt.Println(projectName)
		allValues["projectName"] = projectName

		// INFO OUTPUT
		log.Info("ALL QUESTION-FILES: ", gitConfig.Questions)
		log.Info("ALL TEMPLATE-FILES ", gitConfig.Templates)
		log.Info("RUN INTERACTIVE SURVEY: ", runSurvey)
		log.Info("DEFAULT: GITHUB-REPO: ", gitConfig.GitRepo)
		log.Info("DEFAULT GITHUB-OWNER: ", gitConfig.GitOwner)
		log.Info("DEFAULT ROOTFOLDER: ", gitConfig.RootFolder)

		// LOAD ALL TEMPLATES FILES
		// GET ALL NEEDED VARIABLES FROM ALL TEMPLATE FILES
		// IMPLEMENT HERE!

		// LOAD ALL QUESTION FILES
		for _, questionFile := range gitConfig.Questions {
			questions, _ := modules.LoadQuestionFile(questionFile)
			allQuestions = append(allQuestions, questions...)
		}

		// LOAD AND VERIFY ALL TEMPLATE FILES
		for _, template := range gitConfig.Templates {

			// SPLIT TEMPLATE PATH BY :
			templateFilePaths := strings.Split(template, ":")

			// VERIFY IF TEMPLATE FILE EXISTS
			templateExists, err := sthingsBase.VerifyIfFileOrDirExists(templateFilePaths[0], "file")
			if err != nil {
				log.Fatalf("ERROR VERIFYING TEMPLATE FILE: %v", err)
			}

			if templateExists {
				log.Info("LOCAL TEMPLATE FOUND : ", templateFilePaths[0])
				templateFile := sthingsBase.ReadFileToVariable(templateFilePaths[0])

				renderedOutputFileName, err := sthingsBase.RenderTemplateInline(templateFilePaths[1], renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
				if err != nil {
					fmt.Println(err)
				}

				// renderedTemplate, err := sthingsBase.RenderTemplateInline(templateFile, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
				// if err != nil {
				// 	fmt.Println(err)
				// }
				//fmt.Println(string(renderedTemplate))

				templateData := TemplateData{TemplateFileName: templateFilePaths[0], TemplateContent: templateFile, OutputFileName: string(renderedOutputFileName), RenderedContent: ""}
				allTemplateData = append(allTemplateData, templateData)

			} else {
				log.Error("LOCAL TEMPLATE NOT FOUND : ", templateFilePaths[0])
				os.Exit(3)
			}

		}

		log.Info("ALL GIVE TEMPLATES DO EXIST")

		fmt.Println(allTemplateData)

		// BUILD THE SURVEY + GET DEFAULTS, VALUES FROM FUNCTION CALLS AND RANDOM VALUES
		survey, defaults, err := modules.BuildSurvey(allQuestions)
		if err != nil {
			log.Fatalf("ERROR BUILDING SURVEY: %v", err)
		}

		log.Info("DEFAULTS: ", defaults)

		switch bool(runSurvey) {
		// IF INTERACTIVE - RUN THE SURVEY
		case true:

			// surveys.RunGitHubProjectFlow()

			// //ASK QUESTIONS AND GET THE ANSWERS
			// metaAnswers, err := modules.AskInputQuestions(metaQuestions)
			// if err != nil {
			// 	log.Fatalf("ERROR ASKING META QUESTIONS: %v", err)
			// }

			// // Run the prompts and get the answers
			// branchAnswers, err := modules.AskMultipleChoiceQuestions(branchQuestions)
			// if err != nil {
			// 	log.Fatalf("ERROR RUNNING PROMPTS: %v", err)
			// }

			//ASK QUESTIONS AND GET THE ANSWERS
			// branchingAnswers, err := modules.AskInputQuestions(branchQuestions)
			// if err != nil {
			// 	log.Fatalf("ERROR ASKING META QUESTIONS: %v", err)
			// }

			// fmt.Println(metaAnswers, branchAnswers, branchingAnswers)

			err = survey.Run()
			if err != nil {
				log.Fatalf("ERROR RUNNING SURVEY: %v", err)
			}

			// SET ANWERS TO ALL VALUES
			for _, question := range allQuestions {
				allValues[question.Name] = question.Default
			}

			log.Info("ALL VALUES: ", allValues)

			for _, template := range allTemplateData {

				//fmt.Println(template.TemplateFileName)
				fmt.Println(template.OutputFileName)

			}

			// FOR EACH FILE ASK FOR TARGET PATH - LOOP
			// DEFAULT IS TEMPLATE NAME WITHOUT .TPL
			// FOLDER IS DEFAULT FOLDER QUESTION

			// Run the prompts and get the answers
			commitAnswers, err := modules.AskMultipleChoiceQuestions(commitQuestions)
			if err != nil {
				log.Fatalf("ERROR RUNNING PROMPTS: %v", err)
			}

			surveys.RunGitHubBranchingFlow("test-project", gitConfig)

			fmt.Println(commitAnswers)
			commit := sthingsBase.ConvertStringToBoolean(commitAnswers["commit"].(string))
			fmt.Println(commit)
			if commit {
				fmt.Println("Committing files to branch")
				// filesList := []string{"/tmp/bla:blaa"}
				fmt.Println(token)
				fmt.Println(gitOwner)
				fmt.Println(gitOwner)
				fmt.Println(gitRepo)
				// modules.CreateBranchOnGitHub(token, gitOwner, gitOwner, "kaeffken@sthings.com", gitRepo, "test-commit", metaAnswers["projectName"].(string), filesList)
			} else {
				fmt.Println("Not committing files to branch")
			}

		// IF NON-INTERACTIVE - USE THE DEFAULTS
		case false:
			allValues = defaults
		}

		// RENDERING W/ ALL VALUES
		for _, template := range allTemplateData {

			renderedTemplate, err := sthingsBase.RenderTemplateInline(template.TemplateContent, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(renderedTemplate))

			outputFilePath := outputDir + "/" + template.OutputFileName

			if runSurvey {
				log.Info("CREATING RENDERED TEMPLATES ON DISK")

				outputQuestions := map[string]modules.InputQuestion{
					"Create file to?": {
						Question:  "CREATE RENDERED TEMPATE OF " + template.TemplateFileName + " TO?",
						Default:   outputFilePath,
						MinLength: 5,
						MaxLength: 64,
						Id:        "outputFilePath",
						Type:      "string",
					},
				}

				outputAnswers, err := modules.AskInputQuestions(outputQuestions)
				if err != nil {
					log.Fatalf("ERROR ASKING META QUESTIONS: %v", err)
				}

				outputFilePath = outputAnswers["outputFilePath"].(string)

			}
			// CREATING DIRS
			sthingsBase.CreateNestedDirectoryStructure(GetFolderPath(outputFilePath), 0777)
			log.Info("CREATED DIR: ", GetFolderPath(outputFilePath))

			// CREATING FILE ON DISK
			sthingsBase.WriteDataToFile(outputFilePath, string(renderedTemplate))
			log.Info("RENDERED TEMPLATE WRITTEN TO: ", outputFilePath)

			// templateData := TemplateData{TemplateFileName: templateData[0], TemplateContent: templateFile, OutputFileName: string(renderedOutputFileName), RenderedContent: ""}
			// allTemplateData = append(allTemplateData, templateData)

		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().Bool("survey", false, "run (prompted) survey. default: false")
	createCmd.Flags().String("profile", "tests/vspherevm-workflow.yaml", "workflow profile")
	createCmd.Flags().String("project", "unset", "project name")
	createCmd.Flags().String("output", "", "output directory")
	createCmd.Flags().String("tmp", "/tmp/kaeffken", "tmp directory")

}

func GetFolderPath(filePath string) string {
	return filepath.Dir(filePath)
}
