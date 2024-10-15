package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
)

func (c *CLI) Run(ctx *kong.Context) error {

	if c.DryRun {
		c.Verbose = true
	}

	// Fix dots in extensions to make them compatible with path.Ext()
	c.IncludedExts = FixDot(c.IncludedExts)
	c.ExcludedExts = FixDot(c.ExcludedExts)
	c.StringExts = FixDot(c.StringExts)
	c.ByteExts = FixDot(c.ByteExts)

	// Print the contents of the flags after processing
	if c.Verbose {
		fmt.Println()
		fmt.Printf("Root folder.: %v\n", c.Folder)
		fmt.Printf("DryRun......: %v\n", c.DryRun)
		fmt.Printf("IncludedExts: %v\n", c.IncludedExts)
		fmt.Printf("ExcludedExts: %v\n", c.ExcludedExts)
		fmt.Printf("ExcludedDirs: %v\n", c.ExcludedDirs)
		fmt.Printf("ByteExts....: %v\n", c.ByteExts)
		fmt.Printf("StringExts..: %v\n", c.StringExts)
		fmt.Println()
	}

	if err := c.TraverseDir(c.Folder); err != nil {
		return err
	}
	if c.Verbose {
		fmt.Println("Done.")
	}
	return nil
}
