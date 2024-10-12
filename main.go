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
	pattern := `(\w+(?:-\w+)*(?:'\w+)?)|(\([^)]+\))|[.,!?;:]+|'|\s+`
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
				}
				// Remove the tag and any spaces before/after
				tokens = removeTokenAndSpaces(tokens, i)
				i = j // Reset i to the position before the tag
			} else {
				// Remove the tag if no valid token before it
				tokens = removeTokenAndSpaces(tokens, i)
				i--
			}
		}
		i++
	}
	return tokens
}

func processCaseTransform(tokens []string) []string {
	i := 0
	inQuote := false
	for i < len(tokens) {
		token := tokens[i]

		// Toggle inQuote flag when encountering a quote
		if token == "'" {
			inQuote = !inQuote
		}

		if strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")") {
			content := token[1 : len(token)-1]
			parts := strings.Split(content, ",")
			action := strings.TrimSpace(parts[0])
			count := 1
			if len(parts) > 1 {
				count, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}
			// Apply transformation only if not inside quotes
			if !inQuote {
				transformed := 0
				j := i - 1
				tempInQuote := inQuote
				for transformed < count && j >= 0 {
					if tokens[j] == "'" {
						tempInQuote = !tempInQuote
					}
					if !tempInQuote && strings.TrimSpace(tokens[j]) != "" && isWord(tokens[j]) {
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
			}
			// Remove the tag regardless of whether we're inside quotes
			tokens = removeTokenAndSpaces(tokens, i)
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
			// Remove space before punctuation
			if len(result) > 0 && result[len(result)-1] == " " {
				result = result[:len(result)-1]
			}
			result = append(result, token)
			// Add space after punctuation if necessary
			if i+1 < len(tokens) && tokens[i+1] != " " && !isPunctuation(tokens[i+1]) {
				if token == "," || token == ":" || token == ";" {
					// Add space after commas, colons, and semicolons
					result = append(result, " ")
				}
				// Do not add space after periods, exclamation marks, question marks
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
		cap := false
		if tokens[i] == "A" {
			cap = true
		}
		if strings.ToLower(tokens[i]) == "a" {
			// Find the next non-space token
			j := i + 1
			for j < len(tokens) && strings.TrimSpace(tokens[j]) == "" {
				j++
			}
			if j < len(tokens) && isWord(tokens[j]) {
				nextWord := strings.ToLower(tokens[j])
				if startsWithVowelOrH(nextWord) {
					if !cap {
						tokens[i] = "an"
					} else {
						tokens[i] = "An"
					}

				}
			}
		}
	}
	return tokens
}
func reconstructText(tokens []string) string {
	var builder strings.Builder
	for i, token := range tokens {
		if i > 0 {
			prevToken := tokens[i-1]
			// Determine if a space is needed
			needSpace := true

			if token == " " || prevToken == " " {
				needSpace = false
			} else if isPunctuation(token) || isPunctuation(prevToken) || token == "'" || prevToken == "'" {
				needSpace = false
			}

			if needSpace {
				builder.WriteString(" ")
			}
		}
		builder.WriteString(token)
	}
	return builder.String()
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

func removeTokenAndSpaces(tokens []string, index int) []string {
	// Remove any spaces before the token
	if index > 0 && tokens[index-1] == " " {
		tokens = append(tokens[:index-1], tokens[index:]...)
		index--
	}
	// Remove the token
	tokens = append(tokens[:index], tokens[index+1:]...)
	// Remove any spaces after the token
	if index < len(tokens) && tokens[index] == " " {
		tokens = append(tokens[:index], tokens[index+1:]...)
	}
	return tokens
}
