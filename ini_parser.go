// Package iniparser implements a utility to parse ini files.
package iniparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parser error messages provided when error occurs during parsing and manipulation the data.
const (
    EntityDoesNotExist = "This entity does not exist in the ini data"
    FileDoesNotExist = "No such file or directory"
    WrongFormat = "Wrong INI Format"
)

const (
    // constant used to repersent comment character in ini files.
    commentCharacter = ";"
    entityAssignmentOperator = "="
    openSetionBracket = "["
    closeSetionBracket = "]"
)

// ParserError implements error interface defines errors encountered while using the parser.
type ParserError string

// Error is implementation to Error() method in error interface.
func (e ParserError) Error() string {
    return string(e)
}

// Section type repersent the sections in ini files.
// it is an alias to a map of string and Entity (map[string]Entity).
type Section map[string]Entity

// Entity type repersent the key-value entities in ini files.
// it is an alias to a map of string and string (map[string]string).
type Entity map[string]string

// Parser is used to repersent the parser object used by the user.
type Parser struct {
    iniDataMap Section
}

// New creates a new parser object.
// It returns a pointer to the created parser object.
func New() *Parser {
    p := Parser{make(Section)}
    return &p
}

// LoadFromFile takes an ini file path as its argument then converts it to map of section names and section entities.
func (p *Parser) LoadFromFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return ParserError(fmt.Sprintf("%s: %s", path, FileDoesNotExist))
    } else {
        p.LoadFromString(string(data))
        return nil
    }
}

// LoadFromString is a Parser method reads the data in ini string.
// it converts the string into map of section names and section entities.
func (p *Parser) LoadFromString(iniData string) error {
    var currentSectionName string
    scanner := bufio.NewScanner(strings.NewReader(iniData))

    // Check if reserved characters is used wrongly
    checkWrongCharacthers := func (token string) bool {
        return strings.Contains(token, commentCharacter) ||
        strings.Contains(token, entityAssignmentOperator) ||
        strings.Contains(token, openSetionBracket) ||
        strings.Contains(token, closeSetionBracket)
    }

    // Extract name and value from entity line and return an error if the format is unvalid
    parseEntity := func (entityLine string) (string, string, bool) {
        entityLine = strings.Trim(entityLine, " ")
        keyValueList := strings.Split(entityLine, entityAssignmentOperator)

        // Check that keyValueList contains two values only
        if len(keyValueList) == 2 {
            key := keyValueList[0]
            value := keyValueList[1]
            
            key = strings.Trim(key, " ")
            value = strings.Trim(value, " ")

            if !checkWrongCharacthers(key) && !checkWrongCharacthers(value) {
                // check if key and value not empty strings
                if len(key) > 0 && len(value) > 0 {
                    return key, value, true
                }
            }
        }
        return "", "", false
    }

    // Extract the name of the section and return an error if the format is unvalid
    parseSectionName := func (sectionLine string) (string, bool) {
        sectionLine = strings.Trim(sectionLine, " ")

        // Check if section name is surrounded by "[]"
        if strings.HasPrefix(sectionLine, openSetionBracket) &&
        strings.HasSuffix(sectionLine, closeSetionBracket) {
            sectionLine = strings.Trim(sectionLine, "[ ]")

            // Check if wrong characters are not used
            if !checkWrongCharacthers(sectionLine) {
                return sectionLine, true
            }
        }

        return "", false
    }

    // Parse the data file by iterating over it line by line and extract the data from it.
    for scanner.Scan() {
        line := scanner.Text()
        // Ignore empty files
        if len(line) > 0 {
            // Ignore comment lines
            if !strings.HasPrefix(line, commentCharacter) {
                if sectionName, isCorrect := parseSectionName(line); isCorrect {
                    // Hold the current correct section name
                    currentSectionName = sectionName
                    // Create a new section if it does not exist
                    if _, isExist := p.iniDataMap[currentSectionName]; !isExist {
                        p.iniDataMap[currentSectionName] = make(Entity)
                    }
                } else if name, value, isCorrect := parseEntity(line); isCorrect {
                    p.iniDataMap[currentSectionName][name] = value
                } else {
                    return ParserError(WrongFormat)
                }
            }
        }
    }
    return nil
}

// GetSectionNames returns a slice of the section names.
func (p *Parser) GetSectionNames() []string {
    sectionNames := make([]string, 0)

    for name := range p.iniDataMap {
        sectionNames = append(sectionNames, name)
    }

    return sectionNames
}

// GetSections returns the data parsed as a map of section names and section entities.
func (p *Parser) GetSections() Section {
    resultMap := make(Section)

    for sectionName, sectionData := range p.iniDataMap {
        resultMap[sectionName] = make(Entity)
        for key, name := range sectionData {
            resultMap[sectionName][key] = name
        }
    }
    return resultMap
}

// Get returns the value associated with the given section name and key.
func (p *Parser) Get(sectionName, key string) (string, error) {
    if value, isExist := p.iniDataMap[sectionName][key]; isExist {
        return value, nil
    } else {
        return "", ParserError(EntityDoesNotExist)
    }
}

// Set assign the given value to the given section name and key.
func (p *Parser) Set(sectionName, key, value string) {
    if _, isExist := p.iniDataMap[sectionName]; !isExist {
        p.iniDataMap[sectionName] = make(Entity)
    }
    p.iniDataMap[sectionName][key] = value
}

// String converts the ini data map into string.
func (p *Parser) String() string {
    var result string

    for sectionName, section := range p.iniDataMap {
        result += fmt.Sprintf("[%s]\n", sectionName)
        for name, value := range section {
            result += fmt.Sprintf("%s = %s\n", name, value)
        }
    }

    return result
}

// SaveToFile writes the data map converted into string in a file with the given path.
func (p *Parser) SaveToFile(path string) error {
    file, err := os.Create(path)
    defer file.Close()

    file.WriteString(p.String())
    return err
}
