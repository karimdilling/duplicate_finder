package options

import "flag"

var ExcludeFlag = flag.String("exclude", "", "Space seperated list of files to exclude in the search, i.e. -exclude 'node_modules venv .git'")
