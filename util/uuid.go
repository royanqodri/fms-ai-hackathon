package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func ParseUUID(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func MustParseUUID(s string) uuid.UUID {
	return uuid.MustParse(s)
}

func replaceNumbersToDigits(text string) string {
	numberMap := map[string]string{
		// Indonesia
		"nol": "0", "kosong": "0",
		"satu": "1", "dua": "2", "tiga": "3", "empat": "4", "lima": "5",
		"enam": "6", "tujuh": "7", "delapan": "8", "sembilan": "9",

		// English
		"zero": "0", "oh": "0",
		"one": "1", "two": "2", "three": "3", "four": "4", "five": "5",
		"six": "6", "seven": "7", "eight": "8", "nine": "9",

		"double": "",
		"triple": "",
	}

	words := strings.Fields(text)
	var result []string

	for i := 0; i < len(words); i++ {

		cleanWord := strings.Trim(words[i], ".,!?-")
		lowerWord := strings.ToLower(cleanWord)

		if lowerWord == "double" && i+1 < len(words) {
			nextWord := strings.ToLower(strings.Trim(words[i+1], ".,!?-"))
			if val, ok := numberMap[nextWord]; ok && val != "" {
				result = append(result, val+val)
				i++
				continue
			}
		}

		if lowerWord == "triple" && i+1 < len(words) {
			nextWord := strings.ToLower(strings.Trim(words[i+1], ".,!?-"))
			if val, ok := numberMap[nextWord]; ok && val != "" {
				result = append(result, val+val+val)
				i++
				continue
			}
		}

		if val, ok := numberMap[lowerWord]; ok && val != "" {
			words[i] = strings.Replace(words[i], cleanWord, val, 1)
			result = append(result, words[i])
		} else {
			result = append(result, words[i])
		}
	}

	return strings.Join(result, " ")
}

func ExtractPlateNumber(text string) string {

	upperText := strings.ToUpper(text)

	convertedText := replaceNumbersToDigits(upperText)

	re := regexp.MustCompile(`([A-Z]{1,2})[\s\.\-]*([0-9]{1,4})[\s\.\-]*([A-Z]{1,3})`)
	matches := re.FindStringSubmatch(convertedText)

	if len(matches) > 0 {

		plateNumber := fmt.Sprintf("%s %s %s", matches[1], matches[2], matches[3])
		return plateNumber
	}

	return ""
}
