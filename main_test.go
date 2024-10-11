package main

import (
	"strings"
	"testing"
)

func TestProcessData(t *testing.T) {
	testCases := []struct {
		name           string
		inputText      string
		expectedOutput string
	}{
		{
			name:           "Test Case 1: Hex and Bin Conversion",
			inputText:      "I have 1E (hex) files and 10 (bin) directories.",
			expectedOutput: "I have 30 files and 2 directories.",
		},
		{
			name:           "Test Case 2: Uppercase Transformation",
			inputText:      "this is so exciting (up, 2)!",
			expectedOutput: "this is SO EXCITING!",
		},
		{
			name:           "Test Case 3: Lowercase Transformation",
			inputText:      "I should stop SHOUTING (low).",
			expectedOutput: "I should stop shouting.",
		},
		{
			name:           "Test Case 4: Capitalization Transformation",
			inputText:      "welcome to the brooklyn bridge (cap).",
			expectedOutput: "welcome to the brooklyn Bridge.",
		},
		{
			name:           "Test Case 5: Punctuation Spacing",
			inputText:      "Wait , what !?",
			expectedOutput: "Wait, what!?",
		},
		{
			name:           "Test Case 6: Handling Quotes",
			inputText:      "She said : ' hello world '.",
			expectedOutput: "She said: 'hello world'.",
		},
		{
			name:           "Test Case 7: Article Correction",
			inputText:      "It was a honest mistake.",
			expectedOutput: "It was an honest mistake.",
		},
		{
			name:           "Test Case 8: Combined Transformations",
			inputText:      "There it was. A amazing rock! It has been 10 (bin) years. So exciting (up, 2)!",
			expectedOutput: "There it was. An amazing rock! It has been 2 years. SO EXCITING!",
		},
		{
			name:           "Test Case 9: Complex Punctuation and Quotes",
			inputText:      "As Elton John said: ' I am the most well-known homosexual in the world (cap) '.",
			expectedOutput: "As Elton John said: 'I am the most well-known homosexual in the world'.",
		},
		{
			name:           "Test Case 10: Multiple Punctuations",
			inputText:      "I was thinking ... You were right .",
			expectedOutput: "I was thinking... You were right.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens := tokenize(tc.inputText)
			tokens = processHexBin(tokens)
			tokens = processCaseTransform(tokens)
			tokens = processQuotes(tokens)
			tokens = processPunctuation(tokens)
			tokens = processArticles(tokens)
			//tokens = capitalizeSentences(tokens) // If implemented
			outputText := reconstructText(tokens)

			// Trim spaces for accurate comparison
			outputText = strings.TrimSpace(outputText)
			expectedOutput := strings.TrimSpace(tc.expectedOutput)

			if outputText != expectedOutput {
				t.Errorf("Failed %s\nExpected:\n%s\nGot:\n%s", tc.name, expectedOutput, outputText)
			}
		})
	}
}
