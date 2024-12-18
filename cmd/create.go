/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/stuttgart-things/kaeffken/modules"

	"github.com/spf13/cobra"
)

var (
	allQuestions []*modules.Question
	allValues    = make(map[string]interface{})
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create things",
	Long:  `Create things like rendered resource definitions for storing in gitops repositories.`,
	Run: func(cmd *cobra.Command, args []string) {

		// GET VARS
		questionFiles, _ := cmd.Flags().GetStringSlice("questions")
		runSurvey, _ := cmd.LocalFlags().GetBool("survey")

		// INFO OUTPUT
		log.Info("ALL QUESTION-FILES: ", questionFiles)
		log.Info("RUN INTERACTIVE SURVEY: ", runSurvey)

		// LOAD ALL TEMPLATES FILES
		// GET ALL NEEDED VARIABLES FROM ALL TEMPLATE FILES
		// IMPLEMENT HERE!

		// LOAD ALL QUESTION FILES
		for _, questionFile := range questionFiles {
			questions, _ := modules.LoadQuestionFile(questionFile)
			allQuestions = append(allQuestions, questions...)
		}

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

			fmt.Println(allValues)
		// IF NON-INTERACTIVE - USE THE DEFAULTS
		case false:
			allValues = defaults
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().Bool("survey", false, "run (prompted) survey. default: false")
	createCmd.Flags().StringSlice("questions", []string{}, "question files")

}
