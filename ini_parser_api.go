package iniparser

type INIParser struct {
    dataMap Sections
}

func LoadFromFile(path string) (*INIParser, error) {
    parser := INIParser{dataMap: make(Sections)}
    dataString, err := parser.dataMap.readFile(path)

    if err != nil {
        return nil, err
    }
    parser.dataMap.parseIniString(dataString)
    return &parser, nil
}

func LoadFromString(data string) *INIParser {
    parser := INIParser{dataMap: make(Sections)}
    parser.dataMap.parseIniString(data)
    return &parser
}

func (this *INIParser) GetSectionNames() []string {
    return this.dataMap.getSectionNames()
}

func (this *INIParser) GetSections() Sections {
    userCopy := make(Sections)
    for sectionName, section := range this.dataMap {
        userCopy[sectionName] = make(keys)
        for name, value := range section {
            userCopy[sectionName][name] = value
        }
    }
    return userCopy
}

func (this *INIParser) Get(sectionName, key string) (string, error) {
    return this.dataMap.get(sectionName, key)
}

func (this *INIParser) Set(sectionName, key, value string) error {
    return this.dataMap.set(sectionName, key, value)
}

func (this *INIParser) ToString() string {
    return this.dataMap.toString()
}

func (this *INIParser) SaveToFile(path string) error {
    return this.SaveToFile(path)
}
