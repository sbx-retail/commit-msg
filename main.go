package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	fileMessage := os.Args[1]

	log.Println(fileMessage)
	commitFileContent, err := ioutil.ReadFile(fileMessage)
	if err != nil {
		log.Fatalln(err)
	}

	commitMessage := strings.TrimSpace(string(commitFileContent))
	version := ""

	// TODO: Read version from embarcadero projects (ProductVersion=3.3.0.66;)
	patternEmbarcaderoProjects := "*.dproj"
	matchesEmbarcadero, err := filepath.Glob(patternEmbarcaderoProjects)
	if err != nil {
		log.Fatalln(err)
	}

	if len(matchesEmbarcadero) > 0 {
		dprojFile := matchesEmbarcadero[0]
		embarcaderoFileContent, err := ioutil.ReadFile(dprojFile)
		if err != nil {
			log.Fatalln(err)
		}
		embarcaderoFile := strings.TrimSpace(string(embarcaderoFileContent))
		reEmbarcaderoVersion, _ :=
			regexp.Compile("ProductVersion=[0-9]{1,}.[0-9]{1,}.[0-9]{1,}.[0-9]{1,};")
		matchVersionEmbarcadero := reEmbarcaderoVersion.FindStringSubmatch(embarcaderoFile)
		log.Print(matchVersionEmbarcadero)
		if len(matchVersionEmbarcadero) == 0 {
			log.Fatal("Not found product version in file " + dprojFile)
		}

		version = strings.ReplaceAll(matchVersionEmbarcadero[0], "ProductVersion=", "")
		version = strings.ReplaceAll(version, ";", "")
	}

	// TODO: Read version from dotNET projects
	patternDotNetProjects := "*.csproj"
	matchDotNet, err := filepath.Glob(patternDotNetProjects)
	if err != nil {
		log.Fatalln(err)
	}
	// log.Println(matchDotNet)
	if len(matchDotNet) > 0 {
		versionProjectFile := "Properties/AssemblyInfo.cs"
		assemblyFileContent, err := ioutil.ReadFile(versionProjectFile)
		if err != nil {
			log.Fatalln(err)
		}
		assemblyFile := strings.TrimSpace(string(assemblyFileContent))
		reDotNetVersion, _ :=
			regexp.Compile("AssemblyVersion\\(\"[0-9]{1,}.[0-9]{1,}.[0-9]{1,}.[0-9]{1,}\"\\)")
		matchVersionDotNet := reDotNetVersion.FindStringSubmatch(assemblyFile)
		log.Print(matchVersionDotNet)
		if len(matchVersionDotNet) == 0 {
			log.Fatal("Not found product version in file " + versionProjectFile)
		}

		version = strings.ReplaceAll(matchVersionDotNet[0], "AssemblyVersion(\"", "")
		version = strings.ReplaceAll(version, "\")", "")
	}
	// TODO: Read version from nodeJS projects
	nodejsProjects := "package.json"
	matchNodeJS, err := filepath.Glob(nodejsProjects)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(matchNodeJS)
	if len(matchNodeJS) > 0 {
		packageFilePath := matchNodeJS[0]
		packageFileContent, err := ioutil.ReadFile(packageFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		packageFile := strings.TrimSpace(string(packageFileContent))
		rePackageVersion, _ :=
			regexp.Compile("\"version\": \"[0-9]{1,}.[0-9]{1,}.[0-9]{1,}(.[0-9]{1,}|)\"")
		matchVersionPackage := rePackageVersion.FindStringSubmatch(packageFile)
		log.Print(matchVersionPackage)
		if len(matchVersionPackage) == 0 {
			log.Fatal("Not found product version in file " + packageFilePath)
		}

		version = strings.ReplaceAll(matchVersionPackage[0], "\"version\": \"", "")
		version = strings.ReplaceAll(version, "\"", "")
	}

	version = "[v" + version + "]"

	regEx, _ := regexp.Compile(version)
	match := regEx.Match(commitFileContent)

	if !match {
		ioutil.WriteFile(fileMessage, []byte(commitMessage+" "+version), 0)
	}
	// os.Exit(42)

}
