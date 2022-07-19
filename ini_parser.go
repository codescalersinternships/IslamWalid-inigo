// Package iniparser implements a utility to parse ini files.
package iniparser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Parser error messages
const (
    EntityDoesNotExist = "This entity does not exist in the ini data"
    FileDoesNotExist = "no such file or directory"
)

const (
    // constant used to repersent comment character in ini files.
    commentCharacter = ";"
    entityAssignmentOperator = "="
)

// Regular expression used to match lines that contain secion starter.
var sectionRgx = regexp.MustCompile(`\[.*?\]`)

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
func (p *Parser) LoadFromString(iniData string) {
    var currentSectionName string
    scanner := bufio.NewScanner(strings.NewReader(iniData))

    parseEntity := func (entity string) (string, string) {
        nameValueList := strings.Split(entity, entityAssignmentOperator)
        return strings.Trim(nameValueList[0], " "), strings.Trim(nameValueList[1], " ")
    }

    for scanner.Scan() {
        line := scanner.Text()
        if len(line) > 0 {
            line = strings.Trim(line, " ")
            if !strings.HasPrefix(line, commentCharacter) {
                if sectionRgx.MatchString(line) {
                    currentSectionName = sectionRgx.FindString(line)
                    currentSectionName = strings.Trim(currentSectionName, " [] ")
                    if _, isExist := p.iniDataMap[currentSectionName]; !isExist {
                        p.iniDataMap[currentSectionName] = make(Entity)
                    }
                } else {
                    name, value := parseEntity(line)
                    p.iniDataMap[currentSectionName][name] = value
                }
            }
        }
    }
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
        p.iniDataMap = make(Section)
    }
    if _, isExist := p.iniDataMap[sectionName][key]; !isExist {
        p.iniDataMap[sectionName] = make(Entity)
    }
    p.iniDataMap[sectionName][key] = value
}

// 
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

func (p *Parser) SaveToFile(path string) error {
    file, err := os.Create(path)
    defer file.Close()

    file.WriteString(p.String())
    return err
}
