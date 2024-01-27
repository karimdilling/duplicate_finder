package options

import (
	"flag"
	"fmt"
)

var ExcludeFlag = flag.String("exclude", "", "Space seperated list of files to exclude in the search, i.e. -exclude 'node_modules venv .git'")
var HelpFlag = flag.Bool("h", false, "Prints help output")

func PrintProgramDescription() {
	const programDescription = "Usage: duplicate_finder [OPTION]... [PATH TO FOLDER]...\n" +
		"Use '.' as [PATH TO FOLDER] for the current directory\n" +
		"\n" +
		"This program checks for files with duplicate contents in the specified directory" +
		"\n" +
		"The following [OPTION]s can be set:"

	fmt.Println(programDescription)
}
