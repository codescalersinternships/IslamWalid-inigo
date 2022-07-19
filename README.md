# INI Parser

Go package use to parse [ini files](https://en.wikipedia.org/wiki/INI_file)

## Functionality
- `New`: Creates new parse object
- `LoadFromFile(filePath)`: Fills the parser with the data parsed from the given file path.
- `LoadFromString(iniString)`: Fills the parser with the data parsed from the given string.
- `GetSectionNames`: Returns a slice of section names.
- `GetSections`: Returns the data parsed as a map of section names and section entities.
- `Get(sectionName, key)`: Returns the value with key `key` and section `sectionName`
- `Set(sectionName, key, value)`: assign the given `value` to the given `sectionName` and `key`.
- `String`: Returns the parsed data in a string form.
- `SaveToFile(filePath)`: Save the date to the give `filePath`.

## Examples
Create the parser object:
```go
parser := iniparse.New()
```
Parse from a file:
```go
parser.LoadFromFile("file.ini")
```
Or from a string:
```go
iniText := `[section]
domain = wikipedia.org
[section.subsection]
foo = bar`
parser.LoadFromString(iniText)
```
Get the sections names of your parser object:
```go
sliceOfStringsOfSectionNames := parser.GetSectionNames()
```
Get your ini file as a map:
```go
mapOfIni := parser.GetSections()
```
Set a property:
```go
parser.Set("sectionName", "key", "value")
```
Get a property:
```go
value := parser.Get("sectionName", "key")
```
Return the object as string:
```go
parserAsString := parser.String()
```
Save your object in a file:
```go
parser.SaveToFile("filePath")
```
