package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func embedCmd(cmd *cobra.Command, args []string) {

	// We expect a single argument = folder to process
	if len(args) != 1 {
		cmd.PrintErrln("You must specify the source folder")
		os.Exit(1)
	}

	// Check that -x and -i are not used together
	if len(IncludedExts) > 0 && len(ExcludedExts) > 0 {
		cmd.PrintErrln("Flags -i and -x can't be used together")
		os.Exit(1)
	}

	// Check that -s and -b are not used together
	if len(StringExts) > 0 && len(ByteExts) > 0 {
		cmd.PrintErrln("Flags -b and -s can't be used together")
		os.Exit(1)
	}

	startFolder := args[0]

	if DryRun {
		Verbose = true
	}

	// Fix dots in extensions to make them compatible with path.Ext()
	IncludedExts = FixDot(IncludedExts)
	ExcludedExts = FixDot(ExcludedExts)
	StringExts = FixDot(StringExts)
	ByteExts = FixDot(ByteExts)

	// Print the contents of the flags after processing
	if Verbose {
		cmd.Println()
		cmd.Printf("Root folder.: %v\n", args[0])
		cmd.Printf("DryRun......: %v\n", DryRun)
		cmd.Printf("IncludedExts: %v\n", IncludedExts)
		cmd.Printf("ExcludedExts: %v\n", ExcludedExts)
		cmd.Printf("ExcludedDirs: %v\n", ExcludedDirs)
		cmd.Printf("ByteExts....: %v\n", ByteExts)
		cmd.Printf("StringExts..: %v\n", StringExts)
		cmd.Println()
	}

	if err := TraverseDir(startFolder); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}
	if Verbose {
		cmd.Println("Done.")
	}
}
