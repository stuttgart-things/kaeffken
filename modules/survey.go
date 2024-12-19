/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"gopkg.in/yaml.v2"

	"math/rand"
)

// QUESTION STRUCT TO HOLD THE QUESTION DATA FROM YAML
type Question struct {
	Prompt          string                 `yaml:"prompt"`
	Name            string                 `yaml:"name"`
	Default         string                 `yaml:"default,omitempty"`
	DefaultFunction string                 `yaml:"default_function,omitempty"`
	DefaultParams   map[string]interface{} `yaml:"default_params,omitempty"`
	Options         []string               `yaml:"options"`
	Kind            string                 `yaml:"kind,omitempty"` // "function" instead of "text"
	MinLength       int                    `yaml:"minLength,omitempty"`
	MaxLength       int                    `yaml:"maxLength,omitempty"`
	Type            string                 `yaml:"type,omitempty"` // Updated field to match the YAML
}

// FUNCTION MAPPING
var defaultFunctions = map[string]func(params map[string]interface{}) string{

	"getDefaultFavoriteFood": func(params map[string]interface{}) string {
		if spiceLevel, ok := params["spiceLevel"].(string); ok && spiceLevel != "" {
			return fmt.Sprintf("spicy %s", spiceLevel)
		}
		return "steak"
	},
	"getDefaultDrink": func(params map[string]interface{}) string {
		if temp, ok := params["temperature"].(string); ok && temp != "" {
			return fmt.Sprintf("%s water", temp)
		}
		return "water"
	},
}

// LOAD QUESTIONS FROM YAML FILE
func LoadQuestionFile(filename string) ([]*Question, error) {
	var questions []*Question
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

// BUILD THE SURVEY FUNCTION WITH THE NEW RANDOM SETUP
func BuildSurvey(questions []*Question) (*huh.Form, map[string]interface{}, error) {
	var groupFields []*huh.Group
	answers := make(map[string]interface{}) // To hold question names and resolved default values

	// Create a new random source
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // New random generator

	// Iterate over each question to create the survey fields
	for _, question := range questions {
		var field huh.Field

		// Set up default values for options if applicable
		if question.Default == "" && len(question.Options) > 0 {
			question.Default = question.Options[r.Intn(len(question.Options))] // Random default selection
		}

		// Handle the different question kinds
		switch question.Kind {
		case "function": // Handle "function" kind
			if question.DefaultFunction != "" {
				if fn, ok := defaultFunctions[question.DefaultFunction]; ok {
					question.Default = fn(question.DefaultParams)
				} else {
					return nil, nil, fmt.Errorf("default function %s not found", question.DefaultFunction)
				}
			}

			field = huh.NewInput().
				Title(question.Prompt).
				Value(&question.Default).
				Validate(func(input string) error {
					if len(input) < question.MinLength {
						return fmt.Errorf("input too short, minimum length is %d", question.MinLength)
					}
					if len(input) > question.MaxLength {
						return fmt.Errorf("input too long, maximum length is %d", question.MaxLength)
					}
					return nil
				})

		case "ask": // Handle "ask" kind
			field = huh.NewInput().
				Title(question.Prompt).
				Value(&question.Default).
				Validate(func(input string) error {
					if len(input) < question.MinLength {
						return fmt.Errorf("input too short, minimum length is %d", question.MinLength)
					}
					if len(input) > question.MaxLength {
						return fmt.Errorf("input too long, maximum length is %d", question.MaxLength)
					}
					return nil
				})

			// Store a placeholder for user input
			answers[question.Name] = "" // Will be updated during survey run

		default: // Handle multiple choice select options or other fields
			options := make([]huh.Option[string], len(question.Options))
			for i, opt := range question.Options {
				options[i] = huh.NewOption(opt, opt)
			}

			field = huh.NewSelect[string]().
				Title(question.Prompt).
				Options(options...).
				Value(&question.Default)
		}

		// Determine the data type and store the value correctly in the answers map
		switch question.Type {
		case "boolean": // Store as boolean
			answers[question.Name] = question.Default == "Yes" // Convert Yes/No to true/false

		case "int": // Store as integer
			if intValue, err := strconv.Atoi(question.Default); err == nil {
				answers[question.Name] = intValue
			} else {
				return nil, nil, fmt.Errorf("invalid default value for int type: %s", question.Default)
			}

		default: // Default to string
			answers[question.Name] = question.Default
		}

		// Create a group and add the field to it
		group := huh.NewGroup(field)
		groupFields = append(groupFields, group)
	}

	// Create and return the form along with the answers map
	return huh.NewForm(groupFields...), answers, nil
}

// InputQuestion struct to hold question data for asking input
type InputQuestion struct {
	Question  string
	Default   string
	MinLength int
	MaxLength int
	Type      string
	Id        string
}

// questions := map[string]modules.InputQuestion{
// 	"Do you like Go?": {
// 		Question:  "Do you like Go?",
// 		Default:   "Yes",
// 		MinLength: 2,
// 		MaxLength: 3,
// 		Type:      "boolean",
// 	},
// 	"Enter your name": {
// 		Question:  "Enter your name",
// 		Default:   "",
// 		MinLength: 3,
// 		MaxLength: 20,
// 		Type:      "string",
// 	},
// 	"What's your age?": {
// 		Question:  "What's your age?",
// 		Default:   "25",
// 		MinLength: 1,
// 		MaxLength: 3,
// 		Type:      "int",
// 	},
// }

// askInputQuestions asks questions and returns a map with answers
func AskInputQuestions(questions map[string]InputQuestion) (map[string]interface{}, error) {
	answers := make(map[string]interface{}) // To hold question names and their answers

	for _, iq := range questions {
		var field huh.Field
		var answer string

		// Set default if available
		if iq.Default != "" {
			answer = iq.Default
		}

		// Create the appropriate form field based on the question type
		field = huh.NewInput().
			Title(iq.Question).
			Value(&answer).
			Validate(func(input string) error {
				if len(input) < iq.MinLength {
					return fmt.Errorf("input too short, minimum length is %d", iq.MinLength)
				}
				if len(input) > iq.MaxLength {
					return fmt.Errorf("input too long, maximum length is %d", iq.MaxLength)
				}
				return nil
			})

		// Create the group and form
		group := huh.NewGroup(field)
		form := huh.NewForm(group) // Pass group as argument

		// Run the survey and store the answer
		err := form.Run()
		if err != nil {
			return nil, fmt.Errorf("Error running survey: %v", err)
		}

		// Store the answer in the map based on its type
		switch iq.Type {
		case "boolean":
			if answer == "Yes" {
				answers[iq.Id] = true
			} else if answer == "No" {
				answers[iq.Id] = false
			} else {
				return nil, fmt.Errorf("Invalid answer for boolean question: %s", answer)
			}
		case "int":
			intValue, err := strconv.Atoi(answer)
			if err != nil {
				return nil, fmt.Errorf("Invalid answer for int question: %s", answer)
			}
			answers[iq.Id] = intValue
		default:
			answers[iq.Id] = answer // Default to string
		}
	}

	return answers, nil
}

// questions := map[string]interface{}{
// 	"Do you like Go?":             []string{"Yes", "No"},
// 	"What's your favorite color?": []string{"Red", "Blue", "Green", "Yellow"},
// }

func AskMultipleChoiceQuestions(questions map[string]interface{}) (map[string]interface{}, error) {
	// Create a map to store the answers
	answers := make(map[string]interface{})

	// Iterate over the questions map
	for questionName, option := range questions {
		var field huh.Field
		var answer interface{} // Use an interface{} to handle dynamic types

		// Check if the option is a slice (which means it's a set of possible values/options)
		switch opts := option.(type) {
		case []string:
			// If options are provided, create a select field
			options := make([]huh.Option[string], len(opts))
			for i, opt := range opts {
				options[i] = huh.NewOption(opt, opt)
			}

			// Create a select input, where answer is of type string
			var selectedOption string
			field = huh.NewSelect[string]().
				Title(questionName). // No colon here
				Options(options...).
				Value(&selectedOption)

			// Store the selected option
			answer = selectedOption

		default:
			return nil, fmt.Errorf("unsupported question type for %s", questionName)
		}

		// Create the group and add the field
		group := huh.NewGroup(field)
		// Run the survey
		err := huh.NewForm(group).Run()
		if err != nil {
			return nil, fmt.Errorf("error running survey for %s: %v", questionName, err)
		}

		// Store the user's answer in the answers map
		answers[questionName] = answer
	}

	return answers, nil
}
