package args

import (
	"flag"
	"fmt"
	"os"
)

var Usage = func() {
	fmt.Println(`
Usage: embed [flags] <assets-folder>
Note: extensions must be provided as comma separated without spaces:
      Ex. -s txt,md,json -b ogg,wav,png,jpg,ttf
Flags:
	-h  -help         This screen
	-v  -verbose      Outputs info to standard output
	-s  -string       Files with given extension(s) are 'string' instead of '[]byte'
	-x  -exclude      The extension(s) to be exluded from embedding
	-xd -exclude_dir  The folder(s) to be excluded from embedding
	`)
}

// Argument flags
var (
	Root    string
	Xe      string
	Xd      string
	Verbose bool
	Help    bool
	StrVal  string
)

func init() {
	flag.Usage = Usage

	flag.StringVar(&Xe, "exclude_ext", "*", "")
	flag.StringVar(&Xe, "x", "", "")

	flag.StringVar(&StrVal, "string", "", "")
	flag.StringVar(&StrVal, "s", "", "")

	flag.StringVar(&Xd, "exclude_dir", "", "")
	flag.StringVar(&Xd, "xd", "", "")

	flag.BoolVar(&Verbose, "verbose", false, "")
	flag.BoolVar(&Verbose, "v", false, "")

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
