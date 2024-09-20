package args

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var Usage = func() {
	fmt.Println(`
Usage: embed [flags] <assets-folder>

Flags:
-h  -help         This screen
-d  -dry_run      Sets verbose mode and does not touch the file system (use for testing your flags)
-v  -verbose      Outputs info to standard output
-x  -exclude      Specify extensions to exclude (all others are included)
-i  -include      Specify extensions to include (all others are excluded)
-b  -byte         Files with given extension(s) are []byte (all others are string)
-s  -string       Files with given extension(s) are string (all others are []byte)
-xd -exclude_dir  The folder(s) to be excluded from embedding

Note:  🔸 Extensions must be provided as comma separated without spaces:
          Ex. embed -x exe,bin,lib -s txt,md,json resources
       🔸 Flags -b and -s are mutually exclusive
       🔸 Flags -x and -i are mutually exclusive
-`)
}

// Argument flags
var (
	Root         string
	ie           string
	xe           string
	xd           string
	se           string
	be           string
	Verbose      bool
	Help         bool
	DryRun       bool
	IncludedExts []string
	ExcludedExts []string
	ExcludedDirs []string
	StringExts   []string
	ByteExts     []string
)

func listFromFlag(extFlag string, addDot bool) []string {
	extensions := strings.Split(extFlag, ",")
	list := make([]string, 0)
	for _, iext := range extensions {
		iext = strings.TrimSpace(iext)
		if addDot && iext[0] != '.' {
			iext = "." + iext
		}
		list = append(list, iext)
	}
	return list
}

func collectExts(flg string, addDot bool) []string {
	switch flg {
	case "":
		return make([]string, 0)
	case "all", "none":
		list := make([]string, 1)
		list[0] = flg
		return list
	default:
		return listFromFlag(flg, addDot)
	}
}

func collectByteExts() []string {
	return collectExts(be, true)
}

func collectStringExts() []string {
	return collectExts(se, true)
}

func collectExcludedExts() []string {
	return collectExts(xe, true)
}

func collectIncludedExts() []string {
	return collectExts(ie, true)
}

func collectExcludedDirs() []string {
	return collectExts(xd, false)
}

func init() {
	flag.Usage = Usage

	flag.StringVar(&ie, "include_ext", "*", "")
	flag.StringVar(&ie, "i", "", "")

	flag.StringVar(&xe, "exclude_ext", "*", "")
	flag.StringVar(&xe, "x", "", "")

	flag.StringVar(&se, "string", "", "")
	flag.StringVar(&se, "s", "", "")

	flag.StringVar(&be, "byte", "", "")
	flag.StringVar(&be, "b", "", "")

	flag.StringVar(&xd, "exclude_dir", "", "")
	flag.StringVar(&xd, "xd", "", "")

	flag.BoolVar(&Verbose, "verbose", false, "")
	flag.BoolVar(&Verbose, "v", false, "")

	flag.BoolVar(&DryRun, "dry_run", false, "")
	flag.BoolVar(&DryRun, "d", false, "")

	flag.BoolVar(&Help, "help", false, "")
	flag.BoolVar(&Help, "h", false, "")

	flag.Parse()

	if Help {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	IncludedExts = collectIncludedExts()
	ExcludedExts = collectExcludedExts()
	if len(IncludedExts) > 0 && len(ExcludedExts) > 0 {
		fmt.Println("Error: -x and -i flags cannot be used together")
		flag.Usage()
		os.Exit(1)
	}

	ExcludedDirs = collectExcludedDirs()

	StringExts = collectStringExts()
	ByteExts = collectByteExts()
	if len(StringExts) > 0 && len(ByteExts) > 0 {
		fmt.Println("Error: -b and -s flags cannot be used together")
		flag.Usage()
		os.Exit(1)
	}

	if Verbose {
		fmt.Println()
		fmt.Printf("Extensions to discard....: %v\n", ExcludedExts)
		fmt.Printf("Extensions to embed......: %v\n", IncludedExts)
		fmt.Printf("Folders to discard.......: %v\n", ExcludedDirs)
		fmt.Printf("Extensions of type string: %v\n", StringExts)
		fmt.Printf("Extensions of type []byte: %v\n", ByteExts)
		fmt.Println()
		if len(IncludedExts)+len(ExcludedExts) == 0 {
			fmt.Println("Note: all files will be embedded")
		}
		fmt.Println()
	}
	Root = flag.Arg(0)
}
