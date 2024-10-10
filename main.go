package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input_file> <output_file>")
		return
	}

	inputFileName := os.Args[1]
	outputFileName := os.Args[2]

	inputData, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}
	inputText := string(inputData)

	tokens := tokenize(inputText)
	tokens = processHexBin(tokens)
	tokens = processCaseTransform(tokens)
	tokens = processQuotes(tokens)
	tokens = processPunctuation(tokens)
	tokens = processArticles(tokens)
	outputText := reconstructText(tokens)

	err = ioutil.WriteFile(outputFileName, []byte(outputText), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}

	fmt.Println("Processing complete. Output written to", outputFileName)
}

func processData(input string) string {
	return input
}

func tokenize(text string) []string {
	pattern := `(\w+(?:'\w+)?)|(\([^)]+\))|[.,!?;:]+|'\s*|\s+'|\s+`
	regex := regexp.MustCompile(pattern)
	tokens := regex.FindAllString(text, -1)
	return tokens
}

func processHexBin(tokens []string) []string {
	i := 0
	for i < len(tokens) {
		token := tokens[i]
		if token == "(hex)" || token == "(bin)" {
			// Find the previous non-space token
			j := i - 1
			for j >= 0 && strings.TrimSpace(tokens[j]) == "" {
				j--
			}
			if j >= 0 {
				prevToken := tokens[j]
				base := 16
				if token == "(bin)" {
					base = 2
				}
				number, err := strconv.ParseUint(prevToken, base, 64)
				if err == nil {
					tokens[j] = strconv.FormatUint(number, 10)
					tokens = append(tokens[:i], tokens[i+1:]...)
					i--
				} else {
					fmt.Printf("Error parsing number '%s' with base %d: %v\n", prevToken, base, err)
				}
			} else {
				fmt.Printf("No valid token before '%s' to convert.\n", token)
			}
		}
		i++
	}
	return tokens
}

func processCaseTransform(tokens []string) []string {
	i := 0
	for i < len(tokens) {
		token := tokens[i]

		// Check if the token is a transformation tag
		if strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")") {
			content := token[1 : len(token)-1]
			parts := strings.Split(content, ",")
			action := strings.TrimSpace(parts[0])
			count := 1
			if len(parts) > 1 {
				count, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}

			// Apply transformation to the specified number of previous non-space tokens
			transformed := 0
			j := i - 1
			for transformed < count && j >= 0 {
				if strings.TrimSpace(tokens[j]) != "" && isWord(tokens[j]) {
					switch action {
					case "up":
						tokens[j] = strings.ToUpper(tokens[j])
					case "low":
						tokens[j] = strings.ToLower(tokens[j])
					case "cap":
						tokens[j] = strings.Title(strings.ToLower(tokens[j]))
					}
					transformed++
				}
				j--
			}

			// Remove the tag token
			tokens = append(tokens[:i], tokens[i+1:]...)
			i-- // Adjust index after removing element
		}
		i++
	}
	return tokens
}

func processPunctuation(tokens []string) []string {
	result := []string{}
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if isPunctuation(token) {

			if len(result) > 0 && result[len(result)-1] == " " {
				result = result[:len(result)-1]
			}
			result = append(result, token)

			if i+1 < len(tokens) && tokens[i+1] != " " {
				result = append(result, " ")
			}
		} else {
			result = append(result, token)
		}
	}
	return result
}

func processQuotes(tokens []string) []string {
	inQuote := false
	result := []string{}
	i := 0
	for i < len(tokens) {
		token := tokens[i]
		trimmedToken := strings.TrimSpace(token)

		if trimmedToken == "'" {
			// Toggle the inQuote flag
			inQuote = !inQuote

			if inQuote {
				// Opening quote
				if len(result) > 0 && result[len(result)-1] == " " {
					result = result[:len(result)-1] // Remove space before opening quote
				}
				result = append(result, "'")
				// Skip space after opening quote
				if i+1 < len(tokens) && tokens[i+1] == " " {
					i++
				}
			} else {
				// Closing quote
				if len(result) > 0 && result[len(result)-1] == " " {
					result = result[:len(result)-1] // Remove space before closing quote
				}
				result = append(result, "'")
				// Skip space after closing quote
				if i+1 < len(tokens) && tokens[i+1] == " " {
					i++
				}
			}
		} else {
			result = append(result, token)
		}
		i++
	}
	return result
}

func processArticles(tokens []string) []string {
	for i := 0; i < len(tokens); i++ {
		if strings.ToLower(tokens[i]) == "a" {
			// Find the next non-space token
			j := i + 1
			for j < len(tokens) && strings.TrimSpace(tokens[j]) == "" {
				j++
			}
			if j < len(tokens) && isWord(tokens[j]) {
				nextWord := strings.ToLower(tokens[j])
				if startsWithVowelOrH(nextWord) {
					tokens[i] = "an"
				}
			}
		}
	}
	return tokens
}

func reconstructText(tokens []string) string {
	var builder strings.Builder
	for i, token := range tokens {
		// Determine if a space is needed
		needSpace := i > 0 &&
			token != " " &&
			tokens[i-1] != " " &&
			!isPunctuation(token) &&
			!isPunctuation(tokens[i-1]) &&
			token != "'" &&
			tokens[i-1] != "'"

		if needSpace {
			builder.WriteString(" ")
		}
		builder.WriteString(token)
	}
	// Trim any leading or trailing whitespace
	return strings.TrimSpace(builder.String())
}

func isWord(token string) bool {
	return regexp.MustCompile(`^\w+(?:'\w+)?$`).MatchString(token)
}

func startsWithVowelOrH(word string) bool {
	return strings.HasPrefix(word, "a") || strings.HasPrefix(word, "e") ||
		strings.HasPrefix(word, "i") || strings.HasPrefix(word, "o") ||
		strings.HasPrefix(word, "u") || strings.HasPrefix(word, "h")
}

func isPunctuation(token string) bool {
	return regexp.MustCompile(`^[.,!?;:]+$`).MatchString(token)
}
