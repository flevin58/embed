package args

import (
	"flag"
	"fmt"
	"os"
)

var Usage = func() {
	fmt.Println(`
Usage: embed [flags] <assets-folder>
Flags:
	-h  --help         This screen
	-s  --silent       No output. Useful in unattended scripts
	-x  --exclude      The extensions to be exluded from embedding
	-xd --exclude_dir  The folders to be excluded from embedding
	`)
}

// Argument flags
var (
	Root   string
	Xe     string
	Xd     string
	Silent bool
	Help   bool
)

func init() {
	flag.Usage = Usage

	flag.StringVar(&Xe, "exclude_ext", "*", "")
	flag.StringVar(&Xe, "x", "", "")

	flag.StringVar(&Xd, "exclude_dir", "", "")
	flag.StringVar(&Xd, "xd", "", "")

	flag.BoolVar(&Silent, "silent", false, "")
	flag.BoolVar(&Silent, "s", false, "")

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

	Root = flag.Arg(0)

}
