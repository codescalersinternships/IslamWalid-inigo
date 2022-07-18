package iniparser

import (
	"reflect"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
    iniMap := make(Sections)

    want := `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"
`
    got, _ := iniMap.readFile("ini_files/ref.ini")

    if got != want {
        t.Errorf("\ngot:\n %q\n\nwant:\n%q", got, want)
    }
}

func TestLoadFromString(t *testing.T) {
    got := make(Sections)
    want := Sections{"owner": keys{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": keys{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}

    got.parseIniString(`; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"
`)

    if !reflect.DeepEqual(got, want) {
        t.Errorf("\ngot:\n%v\n\nwant:\n%v", got, want)
    }
}

func TestGetSectionNames(t *testing.T) {
    iniMap := Sections{"owner": keys{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": keys{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}

    want := []string{"owner", "database"}
    got := iniMap.getSectionNames()

    assertList := func (got, want []string) {
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
            t.Errorf("\ngot:\n%v\n\nwant:\n%v", got, want)
        }
    }

    assertList(got, want)
}

func TestGet(t *testing.T) {
    iniName := Sections{"owner": keys{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": keys{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}

    t.Run("key exists", func(t *testing.T) {
        want := "John Doe"
        got, _ := iniName.get("owner", "name")

        if want != got {
            t.Errorf("\ngot:\n%q\n\nwant:\n%q", got, want)
        }
    })

    t.Run("key does not exist", func(t *testing.T) {
        want := "This key does not exist in the ini data"
        _, err := iniName.get("unkown section", "unknown key")

        if err == nil {
            t.Fatal("Exptected \"KeyDoesNotExist\" error.")
        }

        if want != err.Error() {
            t.Errorf("\ngot:\n%q\n\nwant:\n%q", err.Error(), want)
        }
    })
}

func TestSet(t *testing.T) {
    iniMap := Sections{"owner": keys{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": keys{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}

    t.Run("key exists", func(t *testing.T) {
        want := "CodeScalser"
        iniMap.set("owner", "organization", "CodeScalser")
        got, _ := iniMap.get("owner", "organization")

        if got != want {
            t.Errorf("\ngot:\n%q\n\nwant:\n%q", got, want)
        }
    })

    t.Run("key does not exist", func(t *testing.T) {
        want := "This key does not exist in the ini data"
        err := iniMap.set("unkown section", "unknown key", "some value")

        if err == nil {
            t.Fatal("Exptected \"KeyDoesNotExist\" error.")
        }

        if want != err.Error() {
            t.Errorf("\ngot:\n%q\n\nwant:\n%q", err.Error(), want)
        }
    })
}

func TestString(t *testing.T) {
    refIniMap := Sections{"owner": keys{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": keys{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}

    refString := `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"
`

    genString := refIniMap.toString()
    genIniMap := make(Sections)
    genIniMap.parseIniString(genString)

    if !reflect.DeepEqual(genIniMap, refIniMap) {
        t.Errorf("\ngenerated string:\n%q\n\nreference string:\n%q\n\nthe two strings are not equivalent",
        genString, refString)
    }
}

func TestSaveToFile(t *testing.T) {
    want := Sections{"owner": keys{"name": "John Doe", "organization": "Acme Widgets Inc."},
    "database": keys{"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""}}

    want.saveToFile("ini_files/gen.ini")

    got := make(Sections)
    gotString, _:= got.readFile("ini_files/gen.ini")
    got.parseIniString(gotString)

    if !reflect.DeepEqual(got, want) {
        t.Errorf("\ngenerated map from generated file:\n%v\n\nwanted map:\n%v\nthe two files are not equivalent",
        got, want)
    }
}
