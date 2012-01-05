package main

import "os"
import "tritium/packager"
import "tritium/linker"
import "tritium/spec"

func show_usage() {
	println("General purpose Tritium command line interface. Commands are: package, link, test")
}

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]
		if command == "package" {
			pkg := packager.BuildDefaultPackage()
			println(string(pkg.Marshal()))
		} else if command == "link" {
			println("Linking files found in the directory:", os.Args[2])
			linker.RunLinker(os.Args[2])
		} else if command == "test" {
			println("Running tests")
			spec.RunTests()
		} else {
			println("No such command", command)
			show_usage()
		}
	} else {
		show_usage()
	}
}