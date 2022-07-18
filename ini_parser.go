package iniparser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const KeyDoesNotExist = ParserError("This key does not exist in the ini data")

type ParserError string

func (e ParserError) Error() string {
    return string(e)
}

type Sections map[string]keys

type keys map[string]string

func (this Sections) readFile(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err == nil {
        return string(data), nil
    } else {
        return "", err
    }
}

func (this Sections) parseIniString(iniData string) {
    var currentSectionName string
    dataLines := strings.Split(iniData, "\n")
    sectionRgx := regexp.MustCompile(`\[.*?\]`)

    parseKey := func (key string) (string, string) {
        keyList := strings.Split(key, "=")
        return strings.Trim(keyList[0], " "), strings.Trim(keyList[1], " ")
    }

    for _, line := range dataLines {
        if len(line) > 0 {
            line = strings.Trim(line, " ")
            if line[0] != ';' {
                if sectionRgx.MatchString(line) {
                    currentSectionName = sectionRgx.FindString(line)
                    currentSectionName = strings.Trim(currentSectionName, " [] ")
                    if _, isExist := this[currentSectionName]; !isExist {
                        this[currentSectionName] = make(keys)
                    }
                } else {
                    name, value := parseKey(line)
                    this[currentSectionName][name] = value
                }
            }
        }
    }
}

func (this Sections) getSectionNames() []string {
    sectionNames := make([]string, 0)
    
    for name := range this {
        sectionNames = append(sectionNames, name)
    }

    return sectionNames
}

func (this Sections) get(sectionName, name string) (string, error) {
    if value, isExist := this[sectionName][name]; isExist {
        return value, nil
    } else {
        return "", KeyDoesNotExist
    }
}

func (this Sections) set(sectionName, name, value string) error{
    if _, isExist := this[sectionName][name]; isExist {
        this[sectionName][name] = value
        return nil
    } else {
        return KeyDoesNotExist
    }
}

func (this Sections) toString() string {
    var result string

    for sectionName, section := range this {
        result += fmt.Sprintf("[%s]\n", sectionName)
        for name, value := range section {
            result += fmt.Sprintf("%s = %s\n", name, value)
        }
    }

    return result
}

func (this Sections) saveToFile(path string) error {
    file, err := os.Create(path)
    defer file.Close()

    file.WriteString(this.toString())
    return err
}
