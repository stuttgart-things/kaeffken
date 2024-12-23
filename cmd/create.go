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
	files2Commit    []string
	githubPRAnswers = make(map[string]interface{})
	bracketFormat   = "curly"
	allValues       = make(map[string]interface{})
	renderOption    = "missingkey=error"
	brackets        = map[string]TemplateBracket{
		"curly":  TemplateBracket{"{{", "}}", `\{\{(.*?)\}\}`},
		"square": TemplateBracket{"[[", "]]", `\[\[(.*?)\]\]`},
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
		author, _ := cmd.LocalFlags().GetString("author")
		authorMail, _ := cmd.LocalFlags().GetString("mail")
		projectName, _ := cmd.LocalFlags().GetString("project")
		tmpDir, _ := cmd.LocalFlags().GetString("tmp")
		outputDir, _ := cmd.LocalFlags().GetString("output")

		// SET DEFAULTS
		if outputDir == "" {
			outputDir = tmpDir + "/" + time.Now().Format("20060102_150405")
		}

		if projectName == "unset" && runSurvey {
			projectName = "test-project"

			metaQuestions := map[string]modules.InputQuestion{
				"Project name?": {
					Question:  "Project name?",
					Default:   "",
					MinLength: 3,
					MaxLength: 18,
					Id:        "projectName",
					Type:      "string",
				},
			}

			projectAnswers, err := modules.AskInputQuestions(metaQuestions)
			if err != nil {
				log.Fatalf("ERROR ASKING META QUESTIONS: %v", err)
			}

			projectName = projectAnswers["projectName"].(string)
		}

		allValues["projectName"] = projectName
		fmt.Println("ALLL", allValues)

		// READ GIT PROFILE
		gitConfig := surveys.ReadGitProfile(profile)
		log.Info("ALL QUESTION-FILES: ", gitConfig.Questions)
		log.Info("ALL TEMPLATE-FILES ", gitConfig.Templates)
		log.Info("RUN INTERACTIVE SURVEY: ", runSurvey)
		log.Info("DEFAULT: GITHUB-REPO: ", gitConfig.GitRepo)
		log.Info("DEFAULT GITHUB-OWNER: ", gitConfig.GitOwner)
		log.Info("DEFAULT ROOTFOLDER: ", gitConfig.RootFolder)

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

			err = survey.Run()
			if err != nil {
				log.Fatalf("ERROR RUNNING SURVEY: %v", err)
			}

			// SET ANWERS TO ALL VALUES
			for _, question := range allQuestions {
				allValues[question.Name] = question.Default
			}

			log.Info("ALL VALUES: ", allValues)

			// fmt.Println(commitAnswers)
			// commit := sthingsBase.ConvertStringToBoolean(commitAnswers["commit"].(string))
			// fmt.Println(commit)
			// if commit {
			// 	fmt.Println("Committing files to branch")
			// 	// filesList := []string{"/tmp/bla:blaa"}
			// 	fmt.Println(token)
			// 	fmt.Println(gitOwner)
			// 	fmt.Println(gitOwner)
			// 	fmt.Println(gitRepo)
			// 	//
			// } else {
			// 	fmt.Println("Not committing files to branch")
			// }

		// IF NON-INTERACTIVE - USE THE DEFAULTS
		case false:
			allValues = defaults
			allValues["projectName"] = projectName
		}

		// RENDERING W/ ALL VALUES
		for _, template := range allTemplateData {

			// RENDER TEMPLATE
			renderedTemplate, err := sthingsBase.RenderTemplateInline(template.TemplateContent, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(renderedTemplate))

			// RENDER SUBFOLDER
			renderedSubFolder, err := sthingsBase.RenderTemplateInline(gitConfig.SubFolder, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(renderedSubFolder))

			// SET OUPUT FILE PATH
			outputFilePathLocal := outputDir + "/" + template.OutputFileName
			outputFilePathGit := gitConfig.RootFolder + "/" + string(renderedSubFolder) + "/" + template.OutputFileName

			if runSurvey {
				log.Info("CREATING RENDERED TEMPLATES ON DISK")

				outputQuestions := map[string]modules.InputQuestion{
					"Create file to?": {
						Question:  "CREATE RENDERED TEMPATE OF " + template.TemplateFileName + " TO?",
						Default:   outputFilePathLocal,
						MinLength: 5,
						MaxLength: 64,
						Id:        "outputFilePathLocal",
						Type:      "string",
					},
				}

				outputAnswers, err := modules.AskInputQuestions(outputQuestions)
				if err != nil {
					log.Fatalf("ERROR ASKING META QUESTIONS: %v", err)
				}

				outputFilePathLocal = outputAnswers["outputFilePathLocal"].(string)

			}
			// CREATING DIRS
			sthingsBase.CreateNestedDirectoryStructure(GetFolderPath(outputFilePathLocal), 0777)
			log.Info("CREATED DIR: ", GetFolderPath(outputFilePathLocal))

			// CREATING FILE ON DISK
			sthingsBase.WriteDataToFile(outputFilePathLocal, string(renderedTemplate))
			log.Info("RENDERED TEMPLATE WRITTEN TO: ", outputFilePathLocal)
			files2Commit = append(files2Commit, outputFilePathLocal+":"+outputFilePathGit)

		}

		fmt.Println("FILES TO COMMIT: ", files2Commit)

		// CREATE BRANCH AND PR ON GITHUB
		if runSurvey {
			// ASK FOR GITHUB BRANCHING FLOW
			githubPRAnswers = surveys.RunGitHubBranchingFlow(gitConfig, "test-project")
		} else {
			// USE DEFAULTS FROM PROFILE
			githubPRAnswers = surveys.ConfigToMap(gitConfig, "test-project")
		}

		// SET COMMIT MESSAGE
		allValues["commitMessage"] = projectName

		// CREATE BRANCH ON GITHUB
		modules.CreateBranchOnGitHub(token, gitOwner, author, authorMail, gitRepo, allValues["projectName"].(string), allValues["commitMessage"].(string), files2Commit)

		// CREATE PR ON GITHUB
		labels := []string{"infrastructre", "automation"}
		prSubject := "TEST PR"

		commitBranch := allValues["projectName"].(string)
		repoBranch := allValues["projectName"].(string)

		baseBranch := "main"
		prDescription := "PR DESCRIPTION"

		modules.CreatePullRequestOnGitHub(token, prSubject, gitOwner, gitOwner, commitBranch, gitRepo, gitRepo, repoBranch, baseBranch, prDescription, labels)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().Bool("survey", false, "run (prompted) survey. default: false")
	createCmd.Flags().Bool("branch", true, "create branch on github. default: true")
	createCmd.Flags().Bool("pr", true, "create pull request on github. default: true")
	createCmd.Flags().String("profile", "tests/vspherevm-workflow.yaml", "workflow profile")
	createCmd.Flags().String("project", "unset", "project name")
	createCmd.Flags().String("output", "", "output directory")
	createCmd.Flags().String("tmp", "/tmp/kaeffken", "tmp directory")
	createCmd.Flags().String("mail", "kaeffken@sthings.com", "author mail")
	createCmd.Flags().String("author", "kaeffken", "author name")
}

func GetFolderPath(filePath string) string {
	return filepath.Dir(filePath)
}
