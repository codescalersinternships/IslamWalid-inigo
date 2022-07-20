package iniparser

import (
	"fmt"
	"path"
	"reflect"
	"sort"
	"testing"
)

var (
    refIniFile = path.Join("ini_files", "ref.ini")
    genIniFile = path.Join("ini_files", "gen.ini")
)

// Test functions

func TestLoadFromFile(t *testing.T) {
    t.Run("file exists", func(t *testing.T) {
        want := newRefParser()

        got := New()
        got.LoadFromFile(refIniFile)
        assertIniDataMap(t, got.iniDataMap, want.iniDataMap)
    })

    t.Run("file does not exist", func(t *testing.T) {
        got := New()
        err := got.LoadFromFile("unknown file")
        assertError(t, err, "unknown file: No such file or directory")
    })
}

func TestLoadFromString(t *testing.T) {
    t.Run("No format errors", func(t *testing.T) {
        want := newRefParser()
        got := New()
        got.LoadFromString(refIniString())

        assertIniDataMap(t, got.iniDataMap, want.iniDataMap)
    })

    t.Run("Format errors", func(t *testing.T) {
        p := New()

        for _, wrongCase := range wrongIniStrings() {
            err := p.LoadFromString(wrongCase)
            if err == nil {
                t.Fatalf("Excepected error: %q\n when parsing wrong case: %q", WrongFormat, wrongCase)
            } else {
                assertString(t, err.Error(), WrongFormat)
            }
        }
    })

}

func TestGetSectionNames(t *testing.T) {
    p := newRefParser()

    want := []string{"owner", "database"}
    got := p.GetSectionNames()

    assertSectionList(t, got, want)
}

func TestGetSections(t *testing.T) {
    p := newRefParser()

    want := p.iniDataMap
    got := p.GetSections()

    assertIniDataMap(t, want, got)
}

func TestGet(t *testing.T) {
    t.Run("entity exists", func(t *testing.T) {
        p := newRefParser()
        want := "John Doe"
        got, _ := p.Get("owner", "name")
        assertString(t, got, want)
    })

    t.Run("entity does not exist", func(t *testing.T) {
        p := newRefParser()
        _, err := p.Get("unknown section", "unknown key")
        assertError(t, err, "This entity does not exist in the ini data")
    })
}

func TestSet(t *testing.T) {
    p := newRefParser()

    t.Run("entity exists", func(t *testing.T) {
        want := "CodeScalser"
        p.Set("owner", "organization", "CodeScalser")
        got, _ := p.Get("owner", "organization")

        assertString(t, got, want)
    })

    t.Run("entity does not exist", func(t *testing.T) {
        want := "new value"
        p.Set("new section", "new name", "new value")
        got, _ := p.Get("new section", "new name")

        assertString(t, got, want)
    })
}

func TestString(t *testing.T) {
    refPasrser := newRefParser()

    genString := refPasrser.String()
    genParser := New()
    genParser.LoadFromString(genString)

    assertIniDataMap(t, refPasrser.iniDataMap, genParser.iniDataMap)
}

func TestSaveToFile(t *testing.T) {
    want := newRefParser()

    want.SaveToFile(genIniFile)

    got := New()
    got.LoadFromFile(genIniFile)

    assertIniDataMap(t, got.iniDataMap, want.iniDataMap)
}

// Example functions

func ExampleNew() {
    myParser := New()
    // the output is empty since we did not parse any data yet
    fmt.Println(myParser)
    // output:
}

func ExampleParser_LoadFromFile() {
    // Create parser1 and fill it with data parsed from ini_files/ref.ini
    p1 := New()
    p1.LoadFromFile(refIniFile)

    // Create parser2 and fill it with data parsed from the same file
    p2 := New()
    p2.LoadFromFile(refIniFile)

    fmt.Println(reflect.DeepEqual(p1, p2))
    //output: true
}

func ExampleParser_LoadFromString() {
    // Create reference parser and fill it with data parsed from ini_files/ref.ini
    refParser := New()
    refParser.LoadFromFile(refIniFile)

    // Create generated parser from the string resulted from the reference parser
    genParser := New()
    genParser.LoadFromString(refParser.String())

    fmt.Println(reflect.DeepEqual(refParser, genParser))
    //output: true
}

func ExampleParser_GetSectionNames() {
    // Create new parser and fill it with data parsed from ini_files/ref.ini
    p := New()
    p.LoadFromFile(refIniFile)

    sectionNames := p.GetSectionNames()

    // Sort the resulting slice to always match the output example
    sort.Strings(sectionNames)

    fmt.Println(sectionNames)
    // output: [database owner]
}

func ExampleParser_GetSections() {
    // Create parser1 and fill it with data parsed from ini_files/ref.ini
    p1 := New()
    p1.LoadFromFile(refIniFile)

    // Create parser2 and fill it with data parsed from the same file
    p2 := New()
    p2.LoadFromFile(refIniFile)

    fmt.Println(reflect.DeepEqual(p1, p2))
    // output: true
}

func ExampleParser_Get() {
    // Create new parser and fill it with data parsed from ini_files/ref.ini
    p := New()
    p.LoadFromFile(refIniFile)

    valueField, _ := p.Get("owner", "name")
    fmt.Println(valueField)
    // output: John Doe
}

func ExampleParser_Set() {
    // Create new parser and fill it with data parsed from ini_files/ref.ini
    p := New()
    p.LoadFromFile(refIniFile)

    // Sets the entity with key "name" in section "owner" to value "person"
    p.Set("owner", "name", "person")
    valueField, _ := p.Get("owner", "name")
    fmt.Println(valueField)
    // output: person
}

func ExampleParser_String() {
    // Create new parser and fill it with data parsed from ini_files/ref.ini
    p := New()
    p.LoadFromFile(refIniFile)

    // Parse implements String so it is printed with ini form
    fmt.Println(p)
}

func ExampleParser_SaveToFile() {
    // Create new parser
    p := New()

    // Add entity with key "new key", value "new value" in section "new section"
    p.Set("new section", "new key", "new value")

    // Save the date to the given file in ini form
    p.SaveToFile(genIniFile)
}

// Helper functions

func newRefParser() *Parser {
    newObj := Parser{iniDataMap: Section{"owner": Entity{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": Entity{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

    return &newObj
}

func refIniString() string {
    return `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"
`
}

func wrongIniStrings() []string {
    var wrongCases []string
    wrongCases = append(wrongCases, "last modified 1 April 2001 by John Doe")
    wrongCases = append(wrongCases, "[owner")
    wrongCases = append(wrongCases, "name  John Doe")
    wrongCases = append(wrongCases, "[1234.9;890]")
    wrongCases = append(wrongCases, "serv==er = ;192.0.2.62     ")
    wrongCases = append(wrongCases, "port =")
    wrongCases = append(wrongCases, "[]")

    return wrongCases
}

func assertString(t testing.TB, got, want string) {
    t.Helper()

    if got != want {
        t.Errorf("got:\n%q\nwant:\n%q", got, want)
    }
}

func assertError(t testing.TB, err error, want string) {
    t.Helper()
    if err == nil {
        t.Fatalf("Exptected Error: %q", want)
    }
    assertString(t, err.Error(), want)
}

func assertIniDataMap(t testing.TB, got, want Section) {
    t.Helper()

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got:\n%v\nwant:\n%v", got, want)
    }
}

func assertSectionList(t testing.TB, got, want []string) {
    t.Helper()
    wantMap := make(map[string]int)
    gotMap := make(map[string]int)

    for _, elem := range want {
        wantMap[elem]++
    }

    for _, elem := range got {
        gotMap[elem]++
    }

    if !reflect.DeepEqual(wantMap, gotMap) {
        t.Errorf("got:\n%v\n\nwant:\n%v", got, want)
    }
}
