package iniparser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const EntityDoesNotExist = ParserError("This entity does not exist in the ini data")

const (
    commentOperator = ";"
    entityAssignmentOperator = "="
    sectionRgx = `\[.*?\]`
)

type ParserError string

func (e ParserError) Error() string {
    return string(e)
}

type Sections map[string]Entities

type Entities map[string]string

type Parser struct {
    iniDataMap Sections
}

func New() *Parser {
    p := Parser{make(Sections)}
    return &p
}

func (p *Parser) LoadFromString(iniData string) {
    var currentSectionName string
    sectionRgx := regexp.MustCompile(sectionRgx)
    scanner := bufio.NewScanner(strings.NewReader(iniData))

    parseEntity := func (entity string) (string, string) {
        nameValueList := strings.Split(entity, entityAssignmentOperator)
        return strings.Trim(nameValueList[0], " "), strings.Trim(nameValueList[1], " ")
    }

    for scanner.Scan() {
        line := scanner.Text()
        if len(line) > 0 {
            line = strings.Trim(line, " ")
            if !strings.HasPrefix(line, commentOperator) {
                if sectionRgx.MatchString(line) {
                    currentSectionName = sectionRgx.FindString(line)
                    currentSectionName = strings.Trim(currentSectionName, " [] ")
                    if _, isExist := p.iniDataMap[currentSectionName]; !isExist {
                        p.iniDataMap[currentSectionName] = make(Entities)
                    }
                } else {
                    name, value := parseEntity(line)
                    p.iniDataMap[currentSectionName][name] = value
                }
            }
        }
    }
}

func (p *Parser) LoadFromFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    } else {
        p.LoadFromString(string(data))
        return nil
    }
}

func (p *Parser) GetSectionNames() []string {
    sectionNames := make([]string, 0)
    
    for name := range p.iniDataMap {
        sectionNames = append(sectionNames, name)
    }

    return sectionNames
}

func (p *Parser) GetSections() Sections {
    resultMap := make(Sections)

    for sectionName, sectionData := range p.iniDataMap {
        resultMap[sectionName] = make(Entities)
        for key, name := range sectionData {
            resultMap[sectionName][key] = name
        }
    }
    return resultMap
}

func (p *Parser) Get(sectionName, key string) (string, error) {
    if value, isExist := p.iniDataMap[sectionName][key]; isExist {
        return value, nil
    } else {
        return "", EntityDoesNotExist
    }
}

func (p *Parser) Set(sectionName, name, value string) {
    if _, isExist := p.iniDataMap[sectionName]; !isExist {
        p.iniDataMap = make(Sections)
    }
    if _, isExist := p.iniDataMap[sectionName][name]; !isExist {
        p.iniDataMap[sectionName] = make(Entities)
    }
    p.iniDataMap[sectionName][name] = value
}

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
