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

	homerun "github.com/stuttgart-things/homerun-library"
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
	homerunToken    = os.Getenv("HOMERUN_TOKEN")
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
		tmpDir, _ := cmd.LocalFlags().GetString("tmp")
		outputDir, _ := cmd.LocalFlags().GetString("output")
		homerunAddr, _ := cmd.LocalFlags().GetString("homerun")

		if homerunToken == "" {
			log.Warn("ERROR: HOMERUN_TOKEN ENVIRONMENT VARIABLE IS NOT SET")
		}

		// SET DEFAULTS
		if outputDir == "" {
			outputDir = tmpDir + "/" + time.Now().Format("20060102_150405")
		}

		// if projectName == "unset" && runSurvey {
		// 	// ASK FOR PROJECT NAME
		// 	metaQuestions := map[string]modules.InputQuestion{
		// 		"Project name?": {
		// 			Question:  "Project name?",
		// 			Default:   "",
		// 			MinLength: 3,
		// 			MaxLength: 18,
		// 			Id:        "projectName",
		// 			Type:      "string",
		// 		},
		// 	}

		// 	projectAnswers, err := modules.AskInputQuestions(metaQuestions)
		// 	if err != nil {
		// 		log.Fatalf("ERROR ASKING META QUESTIONS: %v", err)
		// 	}

		// 	projectName = projectAnswers["projectName"].(string)
		// }

		// allValues["projectName"] = projectName
		// fmt.Println("ALLL", allValues)

		// READ GIT PROFILE
		gitConfig := surveys.ReadGitProfile(profile)
		log.Info("ALL QUESTION-FILES: ", gitConfig.Questions)
		log.Info("ALL TEMPLATE-FILES ", gitConfig.Templates)
		log.Info("RUN INTERACTIVE SURVEY: ", runSurvey)
		log.Info("DEFAULT: GITHUB-REPO: ", gitConfig.GitRepo)
		log.Info("DEFAULT GITHUB-OWNER: ", gitConfig.GitOwner)
		log.Info("DEFAULT ROOTFOLDER: ", gitConfig.RootFolder)

		// LOAD AND ASK PRE QUESTIONS
		preQuestions, _ := modules.LoadQuestionFile(profile)
		if len(preQuestions) > 0 {
			log.Info("PRE-QUESTIONS FOUND")
		} else {
			log.Info("NO PRE-QUESTIONS FOUND")
		}

		// GET PRE-SURVEY AND DEFAULTS
		preSurvey, preSurveyValues, err := modules.BuildSurvey(preQuestions)
		if err != nil {
			log.Fatalf("ERROR BUILDING SURVEY: %v", err)
		}

		// SET PRE-SURVEY VALUES TO ALL VALUES
		allValues = preSurveyValues

		if runSurvey {
			// SET ANWERS TO ALL VALUES
			err = preSurvey.Run()
			if err != nil {
				log.Fatalf("ERROR RUNNING SURVEY: %v", err)
			}

			// SET ANWERS TO ALL VALUES
			for _, question := range preQuestions {
				allValues[question.Name] = question.Default
			}

		}

		// LOAD ALL QUESTION FILES
		for _, questionFile := range gitConfig.Questions {

			// RENDER QUESTION FILE
			renderedQuestionFilePath, err := sthingsBase.RenderTemplateInline(questionFile, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
			if err != nil {
				log.Error("ERROR RENDERING QUESTION FILE: ", err)
			}
			log.Info("LOADING QUESTION FILE: ", string(renderedQuestionFilePath))

			questions, _ := modules.LoadQuestionFile(string(renderedQuestionFilePath))

			if len(questions) > 0 {
				log.Info("LOADED QUESTIONS FROM FILE: ", len(questions))
			} else {
				log.Warn("NO QUESTIONS FOUND IN FILE: ", string(renderedQuestionFilePath))
			}

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

		// IF NON-INTERACTIVE - USE THE DEFAULTS
		case false:
			// MERGE PRE-SURVEY VALUES AND DEFAULTS
			allValues = sthingsBase.MergeMaps(defaults, allValues)
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

		log.Info("GITHUB PR ANSWERS: ", githubPRAnswers)

		// SET COMMIT MESSAGE
		allValues["commitMessage"] = allValues["projectName"].(string)

		// CREATE BRANCH ON GITHUB
		modules.CreateBranchOnGitHub(token, gitOwner, author, authorMail, gitRepo, allValues["projectName"].(string), allValues["commitMessage"].(string), files2Commit)

		// CREATE PR ON GITHUB
		labels := []string{"infrastructre", "automation"}
		prSubject := "TEST PR"

		baseBranch := "main"
		prDescription := "PR DESCRIPTION"

		modules.CreatePullRequestOnGitHub(token, prSubject, gitOwner, gitOwner, gitConfig.GitBranch, gitRepo, gitRepo, gitConfig.GitBranch, baseBranch, prDescription, labels)

		// SEND NOTIFICATION TO HOMERUN
		if homerunToken != "" {

			message := homerun.Message{
				Title:           "kaeffken",
				Message:         "Memory usage is high",
				Severity:        "INFO",
				Author:          author,
				Timestamp:       time.Now().UTC().Format(time.RFC3339), // Generate current timestamp
				System:          "terraform",
				Tags:            "terraform,kaeffken",
				AssigneeAddress: authorMail,
				AssigneeName:    author,
				Artifacts:       "Admin",
				Url:             "Admin",
			}

			err, respCode := modules.SendMessageToHomerun(homerunAddr, homerunToken, message)
			if err != nil {
				log.Error("ERROR SENDING MESSAGE:", err)
			}

			if respCode != "200 OK" {
				log.Error("UNEXPECTED RESPONSE CODE:", err)
			} else {
				log.Info("NOTIFICATION SENT TO HOMERUN")
			}

		} else {
			log.Warn("ERROR: HOMERUN_TOKEN ENVIRONMENT VARIABLE IS NOT SET")
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().Bool("survey", false, "run (prompted) survey. default: false")
	createCmd.Flags().Bool("branch", true, "create branch on github. default: true")
	createCmd.Flags().Bool("pr", true, "create pull request on github. default: true")
	createCmd.Flags().String("profile", "tests/vspherevm-workflow.yaml", "workflow profile")
	createCmd.Flags().String("output", "", "output directory")
	createCmd.Flags().String("homerun", "https://homerun.homerun-dev.sthings-vsphere.labul.sva.de/generic", "homerun address")
	createCmd.Flags().String("tmp", "/tmp/kaeffken", "tmp directory")
	createCmd.Flags().String("mail", "kaeffken@sthings.com", "author mail")
	createCmd.Flags().String("author", "kaeffken", "author name")
}

func GetFolderPath(filePath string) string {
	return filepath.Dir(filePath)
}
