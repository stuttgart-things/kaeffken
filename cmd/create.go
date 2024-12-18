/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/stuttgart-things/kaeffken/modules"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
)

type TemplateBracket struct {
	begin        string `mapstructure:"begin"`
	end          string `mapstructure:"end"`
	regexPattern string `mapstructure:"regex-pattern"`
}

var (
	allQuestions  []*modules.Question
	allTemplates  []string
	bracketFormat = "curly"
	allValues     = make(map[string]interface{})
	renderOption  = "missingkey=error"
	brackets      = map[string]TemplateBracket{
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
		questionFiles, _ := cmd.Flags().GetStringSlice("questions")
		templateFiles, _ := cmd.Flags().GetStringSlice("templates")
		runSurvey, _ := cmd.LocalFlags().GetBool("survey")

		// INFO OUTPUT
		log.Info("ALL QUESTION-FILES: ", questionFiles)
		log.Info("ALL TEMPLATE-FILES ", templateFiles)
		log.Info("RUN INTERACTIVE SURVEY: ", runSurvey)

		// LOAD ALL TEMPLATES FILES
		// GET ALL NEEDED VARIABLES FROM ALL TEMPLATE FILES
		// IMPLEMENT HERE!

		// LOAD ALL QUESTION FILES
		for _, questionFile := range questionFiles {
			questions, _ := modules.LoadQuestionFile(questionFile)
			allQuestions = append(allQuestions, questions...)
		}

		// LOAD ALL TEMPLATE FILES
		for _, templatePath := range templateFiles {
			templateExists, err := sthingsBase.VerifyIfFileOrDirExists(templatePath, "file")
			if err != nil {
				log.Fatalf("Error verifying template file: %v", err)
			}

			if templateExists {
				log.Info("LOCAL TEMPLATE FOUND : ", templatePath)
				templateFile := sthingsBase.ReadFileToVariable(templatePath)
				allTemplates = append(allTemplates, templateFile)

			} else {
				log.Error("LOCAL TEMPLATE NOT FOUND : ", templatePath)
				os.Exit(3)
			}
		}
		log.Info("ALL TEMPLATES LOADED")

		// BUILD THE SURVEY + GET DEFAULTS, VALUES FROM FUNCTION CALLS AND RANDOM VALUES
		survey, defaults, err := modules.BuildSurvey(allQuestions)
		if err != nil {
			log.Fatalf("Error building survey: %v", err)
		}

		log.Info("DEFAULTS: ", defaults)

		switch bool(runSurvey) {
		// IF INTERACTIVE - RUN THE SURVEY
		case true:
			err = survey.Run()
			if err != nil {
				log.Fatalf("Error running survey: %v", err)
			}

			// SET ANWERS TO ALL VALUES
			for _, question := range allQuestions {
				allValues[question.Name] = question.Default
			}

			log.Info("ALL VALUES: ", allValues)

			renderedTemplate, err := sthingsBase.RenderTemplateInline(allTemplates[0], renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(renderedTemplate))

		// IF NON-INTERACTIVE - USE THE DEFAULTS
		case false:
			allValues = defaults
		}

		// Example questions for input
		questions := map[string]modules.InputQuestion{
			"Do you like Go?": {
				Question:  "Do you like Go?",
				Default:   "Yes",
				MinLength: 2,
				MaxLength: 3,
				Type:      "boolean",
			},
			"Enter your name": {
				Question:  "Enter your name",
				Default:   "",
				MinLength: 3,
				MaxLength: 20,
				Type:      "string",
			},
			"What's your age?": {
				Question:  "What's your age?",
				Default:   "25",
				MinLength: 1,
				MaxLength: 3,
				Type:      "int",
			},
		}

		// Ask questions and get the answers
		answers, err := modules.AskInputQuestions(questions)
		if err != nil {
			log.Fatalf("Error asking questions: %v", err)
		}

		// Output the answers
		fmt.Println("\nYour answers:")
		for question, answer := range answers {
			// Convert interface{} to string for output
			switch v := answer.(type) {
			case string:
				fmt.Printf("%s: %s\n", question, v)
			case bool:
				fmt.Printf("%s: %t\n", question, v)
			case int:
				fmt.Printf("%s: %d\n", question, v)
			default:
				fmt.Printf("%s: %v\n", question, v) // Default for other types
			}
		}

		questions2 := map[string]interface{}{
			"Do you like Go?":             []string{"Yes", "No"},
			"What's your favorite color?": []string{"Red", "Blue", "Green", "Yellow"},
		}

		// Run the prompts and get the answers
		answers2, err := modules.AskMultipleChoiceQuestions(questions2)
		if err != nil {
			log.Fatalf("Error running prompts: %v", err)
		}

		// Output the answers
		fmt.Println("Your answers:")
		for question, answer := range answers2 {
			// Type assert the answer to string to print it
			switch v := answer.(type) {
			case string:
				fmt.Printf("%s: %s\n", question, v)
			default:
				fmt.Printf("%s: %v\n", question, answer)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().Bool("survey", false, "run (prompted) survey. default: false")
	createCmd.Flags().StringSlice("questions", []string{}, "question files")
	createCmd.Flags().StringSlice("templates", []string{}, "template files")
}
