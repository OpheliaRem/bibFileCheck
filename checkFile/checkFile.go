package checkFile

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const yearOfStart = 2022

func CheckFile(f *os.File) {
	checkQuantityOfSources(f)
	checkFieldsOfSources(f)
}

func checkFieldsOfSources(f *os.File) {
	sources := getAllSources(f)

	for _, source := range sources {
		key := strings.Split(source, "{")[0]
		key = strings.ToLower(key)

		fields, ok := requiredFieldsOfSources[key]
		if !ok {
			return
		}

		validateFieldsOfSource(source, fields)
	}

	fmt.Println()
}

func validateFieldsOfSource(source string, requiredFields []string) {
	fieldsWithValues := strings.Split(source, ",")

	name := strings.Split(fieldsWithValues[0], "{")[1]

	var fields []string
	for i := 1; i < len(fieldsWithValues); i++ {
		field := strings.TrimSpace(strings.Split(fieldsWithValues[i], "=")[0])
		fields = append(fields, field)
	}

	for _, requiredField := range requiredFields {
		found := false
		for _, field := range fields {
			if field == requiredField {
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Источник", name+": не было найдено необходимое поле", requiredField)
		}
	}
}

func determineAcademicYear() int {
	now := time.Now()

	var toSubtract int
	if now.Month() > 7 {
		toSubtract = 0
	} else {
		toSubtract = 1
	}

	return now.Year() - yearOfStart - 1 - toSubtract
}

func checkQuantityOfSources(f *os.File) {
	academicYear := determineAcademicYear()

	validateQuantity(
		"Всего источников",
		numOfSourcesRequirements[academicYear]["total"],
		getNumOfAllSources(f),
	)

	validateQuantity(
		"Литература на иностранных языках",
		numOfSourcesRequirements[academicYear]["foreign"],
		getNumOfForeignSources(f),
	)

	validateQuantity(
		"Научно-периодическая литература после 2010 года",
		numOfSourcesRequirements[academicYear]["periodic2010"],
		getNumOfPeriodic2010(f),
	)

	validateQuantity(
		"Литература XXI века",
		numOfSourcesRequirements[academicYear]["modern"],
		getNumOfModern(f),
	)

	fmt.Println()
}

func validateQuantity(reason string, expected int, got int) {
	var message string
	if expected > got {
		message = "Ошибка!"
	} else {
		message = "Верно."
	}

	fmt.Printf(message+" "+reason+": ожидалось %v, найдено %v\n", expected, got)
}

func getNumOfAllSources(file *os.File) int {
	sources := getAllSources(file)

	return len(sources)
}

func getNumOfForeignSources(file *os.File) int {
	sources := getAllSources(file)
	foreignSources := filter(sources, func(s string) bool {
		return isHyphenation(s, "english")
	})

	return len(foreignSources)
}

func getNumOfModern(file *os.File) int {
	sources := getAllSources(file)
	modernSources := filter(sources, func(s string) bool {
		return isSourceNewerThan(s, 2000)
	})

	return len(modernSources)
}

func getNumOfPeriodic2010(file *os.File) int {
	periodicSources := getAllPeriodic(file)
	periodicSources2010 := filter(periodicSources, func(s string) bool {
		return isSourceNewerThan(s, 2010)
	})

	return len(periodicSources2010)
}

func getAllSources(f *os.File) []string {
	_, err := f.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	content := string(data)
	parts := strings.Split(content, "@")
	var sources []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			sources = append(sources, part)
		}
	}
	return sources
}

func filter(s []string, predicate func(string) bool) []string {
	var result []string
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func getAllPeriodic(f *os.File) []string {
	sources := getAllSources(f)

	namesForPeriodicSources := []string{
		"article",
		"Article",
		"ARTICLE",
		"inproceedings",
		"InProceedings",
		"Inproceeding",
		"INPROCEEDINGS",
	}

	periodicSources := filter(sources, func(s string) bool {
		for _, name := range namesForPeriodicSources {
			if strings.Contains(s, name) {
				return true
			}
		}
		return false
	})

	return periodicSources
}

func isSourceNewerThan(source string, targetYear int) bool {
	index := strings.Index(source, "year")

	if index == -1 {
		return false
	}

	for !unicode.IsDigit(rune(source[index])) {
		index++
	}

	var begin int
	var end int

	begin = index
	for unicode.IsDigit(rune(source[index])) {
		index++
	}
	end = index

	year, err := strconv.Atoi(source[begin:end])
	if err != nil {
		panic(err)
	}

	return year >= targetYear
}

func isHyphenation(source string, targetHyphenation string) bool {
	index := strings.Index(source, "hyphenation")

	if index == -1 {
		return false
	}

	for string(source[index]) != "=" {
		index++
	}

	for !unicode.IsLetter(rune(source[index])) {
		index++
	}

	var begin int
	var end int

	begin = index
	for unicode.IsLetter(rune(source[index])) {
		index++
	}
	end = index

	return targetHyphenation == source[begin:end]
}
