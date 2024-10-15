//go:generate embed -i tmpl -s all resources

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

package main

import "github.com/flevin58/embed/cmd"

func main() {
	//cmd.Execute()
	cmd.ParseAndRun()
}
