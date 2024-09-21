package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/flevin58/embed/resources/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Produces the embed.go file in the current folder
func ProduceEmbedGo(root string, files []string) error {

	// Define the template variables
	type Entry struct {
		File string
		Var  string
		Type string
	}
	type Embed struct {
		Package string
		Entries []Entry
	}

	// Initialize them as empty
	embed := Embed{
		Package: path.Base(root),
		Entries: make([]Entry, 0),
	}

	// Populate with each file
	const (
		TypeByte   = "[]byte"
		TypeString = "string"
	)

	for _, file := range files {
		fvar := strings.Replace(file, ".", "_", -1)
		fvar = cases.Title(language.Und).String(fvar)

		// Set the Type to the default value
		var ftype = TypeByte

		// Flag -s modifies the default if items are included
		if len(StringExts) > 0 && IsItemIn(path.Ext(file), StringExts) {
			ftype = TypeString
		}

		// Flag -b modifies the default if items are excluded
		if len(ByteExts) > 0 && !IsItemIn(path.Ext(file), ByteExts) {
			ftype = TypeString
		}

		embed.Entries = append(embed.Entries, Entry{
			File: file,
			Var:  fvar,
			Type: ftype,
		})
	}

	// Parse the template
	var tmplString = templates.Embed_tmpl
	tmpl, err := template.New("embed").Parse(tmplString)
	if err != nil {
		return fmt.Errorf("error parsing template: %s", err.Error())
	}

	if DryRun {
		fmt.Printf("Package %s: %d asset(s) would have been embedded\n", embed.Package, len(embed.Entries))
	} else {
		// Generate the "embed.go" file
		fh, err := os.Create(path.Join(root, "embed.go"))
		if err != nil {
			return fmt.Errorf("error creating embed.go: %s", err.Error())
		}
		defer fh.Close()

		err = tmpl.Execute(fh, embed)
		if err != nil {
			return fmt.Errorf("error creating embed.go: %s", err.Error())
		}
		if Verbose {
			fmt.Printf("Package %s: %d asset(s) embedded\n", embed.Package, len(embed.Entries))
		}
	}
	return nil
}
