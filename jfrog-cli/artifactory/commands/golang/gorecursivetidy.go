package golang

import (
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/artifactory/utils/golang"
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/artifactory/utils/golang/project"
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/artifactory/utils/golang/project/dependencies"
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/utils/cliutils"
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"os"
)

func Execute(targetRepo string, details *config.ArtifactoryDetails) error {
	exists, err := fileutils.IsFileExists("go.mod", false)
	if err != nil {
		return err
	}
	var wd string
	var modContent []byte
	shouldRevertMod := false
	if !exists {
		err = golang.RunGoModInit( "", "// Generated by GoCenter")
		if err != nil {
			return err
		}

		regExp , err := dependencies.GetRegex()
		if err != nil {
			return err
		}
		notEmptyModRegex := regExp.GetNotEmptyModRegex()
		modContent, err = fileutils.ReadFile("go.mod")
		if err != nil {
			return err
		}

		projectPackage := dependencies.Package{}
		projectPackage.SetModContent(modContent)
		packageWithDep := dependencies.PackageWithDeps{Dependency:&projectPackage}
		if !packageWithDep.PatternMatched(notEmptyModRegex) {
			wd, err = os.Getwd()
			if err != nil {
				return err
			}
			log.Debug("Root mod is empty, preparing to run 'go mod tidy'")
			err = golang.RunGoModTidy()
			if err != nil {
				return err
			}
			shouldRevertMod = true
		} else {
			log.Debug("Root project mod not empty.")
		}
	}

	goProject, err := project.Load("-")
	if err != nil {
		cliutils.ExitOnErr(err)
	}

	err = goProject.DownloadFromVcsAndPublish(targetRepo, "", true, false,true, details)
	if err != nil {
		cliutils.ExitOnErr(err)
	}

	if shouldRevertMod {
		log.Debug("Reverting to original mod")
		err = os.Chdir(wd)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile("go.mod", modContent, 0700)
	}
	return err
}