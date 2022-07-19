package iniparser

import (
	"path"
	"reflect"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
    t.Run("file exists", func(t *testing.T) {
        want := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
        "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

        got := New()
        got.LoadFromFile(path.Join("ini_files", "ref.ini"))
        assertIniDataMap(t, got.iniDataMap, want.iniDataMap)
    })

    t.Run("file does not exist", func(t *testing.T) {
        got := New()
        err := got.LoadFromFile("unknown file")
        assertError(t, err, "open unknown file: no such file or directory")
    })
}

func TestLoadFromString(t *testing.T) {
    want := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

    got := New()
    got.LoadFromString(`; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"
`)

    assertIniDataMap(t, got.iniDataMap, want.iniDataMap)
}

func TestGetSectionNames(t *testing.T) {
    p := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

    want := []string{"owner", "database"}
    got := p.GetSectionNames()

    assertSectionList(t, got, want)
}

func TestGestSections(t *testing.T) {
    p := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

    want := p.iniDataMap
    got := p.GetSections()

    assertIniDataMap(t, want, got)
}

func TestGet(t *testing.T) {
    t.Run("entity exists", func(t *testing.T) {
        p := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
        "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

        want := "John Doe"
        got, _ := p.Get("owner", "name")

        assertString(t, got, want)
    })

    t.Run("entity does not exist", func(t *testing.T) {
        p := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
        "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

        _, err := p.Get("unknown section", "unknown key")

        assertError(t, err, "This entity does not exist in the ini data")
    })
}

func TestSet(t *testing.T) {
    p := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

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
    refPasrser := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

    genString := refPasrser.String()
    genParser := New()
    genParser.LoadFromString(genString)

    assertIniDataMap(t, refPasrser.iniDataMap, genParser.iniDataMap)
}

func TestSaveToFile(t *testing.T) {
    want := Parser{iniDataMap: Sections{"owner": Entities{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": Entities{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}}

    genFilePath := path.Join("ini_files", "gen.ini")
    want.SaveToFile(genFilePath)

    got := New()
    got.LoadFromFile(genFilePath)

    assertIniDataMap(t, got.iniDataMap, want.iniDataMap)
}

func assertString(t testing.TB, got, want string)  {
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

func assertIniDataMap(t testing.TB, got, want Sections) {
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
