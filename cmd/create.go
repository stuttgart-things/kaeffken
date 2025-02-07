/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	homerun "github.com/stuttgart-things/homerun-library"
	"github.com/stuttgart-things/kaeffken/modules"
	"github.com/stuttgart-things/kaeffken/surveys"
	"gopkg.in/yaml.v2"

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
	secretsToOutput = make(map[string]interface{})
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

		// GET VARS FROM FLAGS
		runSurvey, _ := cmd.LocalFlags().GetBool("survey")
		profile, _ := cmd.LocalFlags().GetString("profile")
		author, _ := cmd.LocalFlags().GetString("author")
		authorMail, _ := cmd.LocalFlags().GetString("mail")
		tmpDir, _ := cmd.LocalFlags().GetString("tmp")
		outputDir, _ := cmd.LocalFlags().GetString("output")
		homerunAddr, _ := cmd.LocalFlags().GetString("homerun")

		// SET DEFAULTS
		if homerunToken == "" {
			log.Warn("ERROR: HOMERUN_TOKEN ENVIRONMENT VARIABLE IS NOT SET")
		}
		if outputDir == "" {
			outputDir = tmpDir + "/" + time.Now().Format("20060102_150405")
		}

		// READ + OUTPUT GIT PROFILE
		gitConfig := surveys.ReadGitProfile(profile)
		log.Info("ALL QUESTION-FILES: ", gitConfig.Questions)
		log.Info("ALL TEMPLATE-FILES ", gitConfig.Templates)
		log.Info("RUN INTERACTIVE SURVEY: ", runSurvey)
		log.Info("DEFAULT: GITHUB-REPO: ", gitConfig.GitRepo)
		log.Info("DEFAULT GITHUB-OWNER: ", gitConfig.GitOwner)
		log.Info("DEFAULT ROOTFOLDER: ", gitConfig.RootFolder)
		log.Info("PULL REQUEST TAGS: ", gitConfig.PrTags)
		log.Info("ALIASES: ", gitConfig.Aliases)
		log.Info("SECRET-ALIASES: ", gitConfig.SecretAliases)
		log.Info("SECRET-FILES: ", gitConfig.SecretFiles)
		log.Info("VALUE-FILES: ", gitConfig.Values)

		// LOAD AND ASK PRE-QUESTIONS
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

		// READ VALUES (IF DEFINED)
		if len(gitConfig.Values) > 0 {

			// LOOP OVER VALUES
			for _, valueFile := range gitConfig.Values {

				// RENDER VALUE FILE
				renderedValueFilePath, err := sthingsBase.RenderTemplateInline(valueFile, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
				if err != nil {
					log.Error("ERROR RENDERING VALUE FILE: ", err)
				}

				// VERIFY IF VALUE FILE EXISTS
				valueFileExists, err := sthingsBase.VerifyIfFileOrDirExists(string(renderedValueFilePath), "file")
				if err != nil {
					log.Fatalf("ERROR VERIFYING VALUE FILE: %v", err)
				}

				if valueFileExists {
					log.Info("VALUE FILE FOUND: ", string(renderedValueFilePath))
					valueFileContent := sthingsBase.ReadFileToVariable(string(renderedValueFilePath))

					// UNMARSHAL YAML
					valueFileValues := make(map[string]interface{})
					err = yaml.Unmarshal([]byte(valueFileContent), &valueFileValues)
					if err != nil {
						log.Error("ERROR UNMARSHALING YAML: ", err)
					}

					// MERGE VALUE FILE VALUES TO ALL VALUES
					allValues = sthingsBase.MergeMaps(valueFileValues, allValues)

				} else {
					log.Error("VALUE FILE NOT FOUND: ", string(renderedValueFilePath))
					os.Exit(3)
				}

			}

			fmt.Println("VALUES: ", gitConfig.Values)
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

			renderedTemplateFileName, err := sthingsBase.RenderTemplateInline(templateFilePaths[0], renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
			if err != nil {
				fmt.Println(err)
			}

			// VERIFY IF TEMPLATE FILE EXISTS
			templateExists, err := sthingsBase.VerifyIfFileOrDirExists(string(renderedTemplateFileName), "file")
			if err != nil {
				log.Fatalf("ERROR VERIFYING TEMPLATE FILE: %v", err)
			}

			if templateExists {
				log.Info("LOCAL TEMPLATE FOUND : ", string(renderedTemplateFileName))
				templateFile := sthingsBase.ReadFileToVariable(string(renderedTemplateFileName))

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

		// RENDER ALIASES
		if len(gitConfig.Aliases) > 0 {
			allValues = modules.RenderAliases(gitConfig.Aliases, allValues)
		} else {
			log.Info("NO ALIASES FOUND")
		}

		// RENDER SECRET ALIASES + SECRET FILES
		if len(gitConfig.SecretAliases) > 0 {
			allValues = modules.RenderAliases(gitConfig.SecretAliases, allValues)
		} else {
			log.Info("NO SECRET ALIASES FOUND")
		}

		//RENDER TEMPLATES W/ ALL VALUES
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
					"SAVE FILE TO?": {
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
			sthingsBase.CreateNestedDirectoryStructure(modules.GetFolderPath(outputFilePathLocal), 0777)
			log.Info("CREATED DIR: ", modules.GetFolderPath(outputFilePathLocal))

			// CREATING FILE ON DISK
			sthingsBase.WriteDataToFile(outputFilePathLocal, string(renderedTemplate))
			log.Info("RENDERED TEMPLATE WRITTEN TO: ", outputFilePathLocal)
			files2Commit = append(files2Commit, outputFilePathLocal+":"+outputFilePathGit)

		}

		// CREATE SECRET FILE
		if len(gitConfig.SecretFiles) > 0 {
			allSecretsFromAllSecretsFile := modules.GetAllSecretsFromSopsDecyptedFiles(gitConfig.SecretFiles, allValues)
			fmt.Println("ALL SECRETS", allSecretsFromAllSecretsFile)

			// ADD ONLY TO BE CREATED SECRETS TO

			// LOOP OVER SECRET ALIASES
			for i, secretAlias := range gitConfig.SecretAliases {

				secretOutputName, err := sthingsBase.RenderTemplateInline(secretAlias, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
				if err != nil {
					log.Error("ERROR RENDERING QUESTION FILE: ", err)
				}

				// SPLIT SECRET ALIAS BY :
				secretOutputNameSplit := strings.Split(string(secretOutputName), ":")
				sourceKey := secretOutputNameSplit[0]
				targetKey := secretOutputNameSplit[1]

				fmt.Println(i, sourceKey)
				fmt.Println(i, targetKey)

				// CHECK IF SECRET ALIAS IS IN ALL SECRETS
				if _, ok := allSecretsFromAllSecretsFile[sourceKey]; ok {
					log.Info("SECRET ALIAS FOUND: ", sourceKey)
					secretsToOutput[targetKey] = allSecretsFromAllSecretsFile[sourceKey]

				} else {
					log.Warn("SECRET ALIAS NOT FOUND: ", sourceKey)
				}
			}

			log.Info("SECRETS TO OUTPUT: ", secretsToOutput)

			// CONVERT MAP TO YAML
			encryptedSecrets, err := yaml.Marshal(secretsToOutput)
			if err != nil {
				fmt.Printf("ERROR CONVERTING TO YAML: %v\n", err)
				return
			}

			// PRINT THE YAML
			fmt.Println(string(encryptedSecrets))

			// READ AGE KEY FROM ENV
			// CHECH IF AGE KEY IS SET IN ENV
			if os.Getenv("AGE") == "" {
				log.Error("ERROR: AGE KEY ENVIRONMENT VARIABLE IS NOT SET - SKIPPED CREATING SECRET")
			} else {
				ageKey := os.Getenv("AGE")
				encryptedSecret := sthingsCli.EncryptStore(ageKey, string(encryptedSecrets))
				fmt.Println(encryptedSecret)

				// RENDER OUTPUT FILE PATH
				secretOutputName, err := sthingsBase.RenderTemplateInline(gitConfig.SecretFileOutputName, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
				if err != nil {
					log.Error("ERROR RENDERING QUESTION FILE: ", err)
				}

				// RENDER SUBFOLDER
				renderedSubFolder, err := sthingsBase.RenderTemplateInline(gitConfig.SubFolder, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
				if err != nil {
					fmt.Println(err)
				}

				// WRITE ENCRYPTED SECRET TO FILE
				secretOutputPath := outputDir + "/" + string(secretOutputName)
				outputFilePathGit := gitConfig.RootFolder + "/" + string(renderedSubFolder) + "/" + string(secretOutputName)

				sthingsBase.WriteDataToFile(secretOutputPath, encryptedSecret)
				log.Info("SECRET OUTPUT: ", secretOutputPath)
				files2Commit = append(files2Commit, secretOutputPath+":"+outputFilePathGit)
			}

		}

		fmt.Println("FILES TO COMMIT: ", files2Commit)

		// CREATE BRANCH AND PR ON GITHUB
		if runSurvey {
			// ASK FOR GITHUB BRANCHING FLOW
			githubPRAnswers = surveys.RunGitHubBranchingFlow(gitConfig, allValues)
		} else {
			// USE DEFAULTS FROM PROFILE
			githubPRAnswers = surveys.ConfigToMap(gitConfig)
			//MERGE
			githubPRAnswers = sthingsBase.MergeMaps(githubPRAnswers, allValues)

		}

		log.Info("GIT-REPOSITORY: ", githubPRAnswers["gitRepo"].(string))
		log.Info("GIT-OWNER: ", githubPRAnswers["gitOwner"].(string))
		log.Info("GIT-BRANCH: ", githubPRAnswers["gitBranch"].(string))
		log.Info("COMMIT-MESSAGE: ", githubPRAnswers["commitMessage"].(string))
		log.Info("PULL REQUEST TITLE: ", githubPRAnswers["prTitle"].(string))
		log.Info("PULL REQUEST DESCRIPTION: ", githubPRAnswers["prDescription"].(string))

		// CREATE BRANCH ON GITHUB
		modules.CreateBranchOnGitHub(
			token,
			githubPRAnswers["gitOwner"].(string),
			author,
			authorMail,
			githubPRAnswers["gitRepo"].(string),
			githubPRAnswers["gitBranch"].(string),
			githubPRAnswers["commitMessage"].(string),
			files2Commit,
		)

		// CREATE PR ON GITHUB
		modules.CreatePullRequestOnGitHub(
			token,
			githubPRAnswers["prTitle"].(string),
			githubPRAnswers["gitOwner"].(string),
			githubPRAnswers["gitOwner"].(string),
			githubPRAnswers["gitBranch"].(string),
			githubPRAnswers["gitRepo"].(string),
			githubPRAnswers["gitRepo"].(string),
			githubPRAnswers["gitBranch"].(string),
			"main",
			githubPRAnswers["prDescription"].(string),
			gitConfig.PrTags,
		)

		// SEND NOTIFICATION TO HOMERUN
		if homerunToken != "" {

			message := homerun.Message{
				Title:           "[PR][KAEFFKEN] - " + githubPRAnswers["prTitle"].(string) + "] created",
				Message:         "REPOSITORY: https://github.com/" + githubPRAnswers["gitOwner"].(string) + "/" + githubPRAnswers["gitRepo"].(string) + " BRANCH: " + githubPRAnswers["gitBranch"].(string),
				Severity:        "SUCCESS",
				Author:          author,
				Timestamp:       time.Now().UTC().Format(time.RFC3339), // Generate current timestamp
				System:          "kaeffken",
				Tags:            strings.Join(gitConfig.PrTags, ",") + ",kaeffken",
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
