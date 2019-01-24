package dependencies

import (
	"fmt"
	golangutil "github.com/jfrog/jfrog-cli-go/jfrog-cli/artifactory/utils/golang"
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/utils/config"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Represents go dependency when running with deps-tidy set to true.
type PackageWithDeps struct {
	Dependency             *Package
	transitiveDependencies []PackageWithDeps
	regExp                 *RegExp
	recursiveTidyOverwrite bool
	shouldRevertToEmptyMod bool
	cachePath              string
	GoModEditMessage       string
	originalModContent     []byte
}

// Creates a new dependency
func (pwd *PackageWithDeps) New(cachePath string, dependency Package, recursiveTidyOverwrite bool) GoPackage {
	pwd.Dependency = &dependency
	pwd.cachePath = cachePath
	pwd.transitiveDependencies = nil
	pwd.recursiveTidyOverwrite = recursiveTidyOverwrite
	return pwd
}

// Populate the mod file and publish the dependency and it's transitive dependencies to Artifactory
func (pwd *PackageWithDeps) PopulateModAndPublish(targetRepo string, cache *golangutil.DependenciesCache, details *config.ArtifactoryDetails) error {
	var path string
	log.Debug("Starting to work on", pwd.Dependency.GetId())
	dependenciesMap := cache.GetMap()
	published, _ := dependenciesMap[pwd.Dependency.GetId()]
	if published {
		log.Debug("Overwriting the mod file in the cache from the one from Artifactory", pwd.Dependency.GetId())
		moduleAndVersion := strings.Split(pwd.Dependency.GetId(), ":")
		path = downloadModFileFromArtifactoryToLocalCache(pwd.cachePath, targetRepo, moduleAndVersion[0], moduleAndVersion[1], details, httpclient.NewDefaultHttpClient())
		err := pwd.updateModContent(path, cache)
		logError(err)
	}
	// Checks if mod is empty, need to run go mod tidy command to populate the mod file.
	log.Debug(fmt.Sprintf("Dependency %s mod file is empty: %t", pwd.Dependency.GetId(), !pwd.PatternMatched(pwd.regExp.GetNotEmptyModRegex())))

	// Creates the dependency in the temp folder and runs go commands: go mod tidy and go mod graph.
	// Returns the path to the project in the temp and the a map with the project dependencies
	path, output, err := pwd.createDependencyAndPrepareMod(cache)
	logError(err)
	pwd.publishDependencyAndPopulateTransitive(path, targetRepo, output, cache, details)
	return nil
}

// Updating the new mod content
func (pwd *PackageWithDeps) updateModContent(path string, cache *golangutil.DependenciesCache) error {
	if path != "" {
		modContent, err := ioutil.ReadFile(path)
		if err != nil {
			cache.IncrementFailures()
			return errorutils.CheckError(err)
		}
		pwd.Dependency.SetModContent(modContent)
	}
	return nil
}

// Init the dependency information if needed.
func (pwd *PackageWithDeps) Init() error {
	var err error
	pwd.regExp, err = GetRegex()
	if err != nil {
		return err
	}
	return nil
}

// Returns true if regex found a match otherwise false.
func (pwd *PackageWithDeps) PatternMatched(regExp *regexp.Regexp) bool {
	lines := strings.Split(string(pwd.Dependency.modContent), "\n")
	for _, line := range lines {
		if regExp.FindString(line) != "" {
			return true
		}
	}
	return false
}

// Creates the dependency in the temp folder and runs go mod tidy and go mod graph
// Returns the path to the project in the temp and the a map with the project dependencies
func (pwd *PackageWithDeps) createDependencyAndPrepareMod(cache *golangutil.DependenciesCache) (path string, output map[string]bool, err error) {
	path, err = pwd.getModPathAndUnzipDependency(path)
	if err != nil {
		return
	}
	pwd.shouldRevertToEmptyMod = false
	// Check the mod in the cache if empty or not
	if pwd.PatternMatched(pwd.regExp.GetNotEmptyModRegex()) {
		err = pwd.useCachedMod(path)
		if err != nil {
			return
		}
	} else {
		published, _ := cache.GetMap()[pwd.Dependency.GetId()]
		if !published {
			output, err = pwd.prepareUnpublishedDependency(path)
			return
		} else {
			pwd.prepareResolvedDependency(path)
		}
	}
	output, err = runGoModGraph()
	return
}

func (pwd *PackageWithDeps) prepareResolvedDependency(path string) {
	// Put the mod file to temp
	err := writeModContentToModFile(path, pwd.Dependency.GetModContent())
	logError(err)
	// If not empty --> use the mod file and don't run go mod tidy
	// If empty --> Run go mod tidy. Publish the package with empty mod file.
	if !pwd.PatternMatched(pwd.regExp.GetNotEmptyModRegex()) {
		log.Debug("The mod still empty after downloading from Artifactory:", pwd.Dependency.GetId())
		originalModContent := pwd.Dependency.GetModContent()
		pwd.prepareAndRunTidy(path, originalModContent)
	} else {
		log.Debug("Project mod file is not empty after downloading from Artifactory", pwd.Dependency.id)
	}
}

func (pwd *PackageWithDeps) prepareAndRunTidy(path string, originalModContent []byte) {
	err := populateModWithTidy(path)
	logError(err)
	// Need to remember here to revert to the empty mod file.
	pwd.shouldRevertToEmptyMod = true
	pwd.originalModContent = originalModContent
}

func (pwd *PackageWithDeps) prepareUnpublishedDependency(pathToModFile string) (output map[string]bool, err error) {
	err = pwd.prepareAndRunInit(pathToModFile)
	if err != nil {
		log.Error(err)
		exists, err := fileutils.IsFileExists(pathToModFile, false)
		logError(err)
		if !exists {
			// Create a mod file
			err = writeModContentToModFile(pathToModFile, pwd.Dependency.GetModContent())
			logError(err)
		}
	}
	// Got here means init worked or mod was created. Need to check the content if mod is empty or not
	modContent, err := ioutil.ReadFile(pathToModFile)
	logError(err)
	originalModContent := pwd.Dependency.GetModContent()
	pwd.Dependency.SetModContent(modContent)
	// If not empty --> use the mod file and don't run go mod tidy
	// If empty --> Run go mod tidy. Publish the package with empty mod file.
	if !pwd.PatternMatched(pwd.regExp.GetNotEmptyModRegex()) {
		log.Debug("The mod still empty after running 'go mod init' for:", pwd.Dependency.GetId())
		pwd.prepareAndRunTidy(pathToModFile, originalModContent)
		output, err = runGoModGraph()
		return
	} else {
		log.Debug("Project mod file after init is not empty", pwd.Dependency.id)
		output, err = runGoModGraph()
		if err != nil {
			log.Debug(fmt.Sprintf("Command go mod graph finished with the following error: %s for dependency %s", err.Error(), pwd.Dependency.GetId()))
			// Graph failed after init. Lets return to empty mod and then run tidy on it and graph again.
			// First create an empty mod.
			pwd.Dependency.SetModContent(originalModContent)
			pwd.prepareAndRunTidy(pathToModFile, originalModContent)
			output, err = runGoModGraph()
		}
	}
	return
}

func (pwd *PackageWithDeps) useCachedMod(path string) error {
	// Mod not empty in the cache. Use it.
	log.Debug("Using the mod in the cache since not empty:", pwd.Dependency.GetId())
	err := writeModContentToModFile(path, pwd.Dependency.GetModContent())
	logError(err)
	err = os.Chdir(filepath.Dir(path))
	if errorutils.CheckError(err) != nil {
		return err
	}
	logError(removeGoSum(path))
	return nil
}

func (pwd *PackageWithDeps) getModPathAndUnzipDependency(path string) (string, error) {
	err := os.Unsetenv(golangutil.GOPROXY)
	if err != nil {
		return "", err
	}
	// Unzips the zip file into temp
	tempDir, err := createDependencyInTemp(pwd.Dependency.GetZipPath())
	if err != nil {
		return "", err
	}
	path = pwd.getModPathInTemp(tempDir)
	return path, err
}

func (pwd *PackageWithDeps) prepareAndRunInit(pathToModFile string) error {
	log.Debug("Preparing to init", pathToModFile)
	err := os.Chdir(filepath.Dir(pathToModFile))
	if errorutils.CheckError(err) != nil {
		return err
	}
	exists, err := fileutils.IsFileExists(pathToModFile, false)
	logError(err)
	if exists {
		err = os.Remove(pathToModFile)
		logError(err)
	}
	// Mod empty.
	// If empty, run go mod init
	moduleId := pwd.Dependency.GetId()
	moduleInfo := strings.Split(moduleId, ":")
	return golangutil.RunGoModInit(replaceExclamationMarkWithUpperCase(moduleInfo[0]), pwd.GoModEditMessage)
}

func writeModContentToModFile(path string, modContent []byte) error {
	return ioutil.WriteFile(path, modContent, 0700)
}

func (pwd *PackageWithDeps) getModPathInTemp(tempDir string) string {
	moduleId := pwd.Dependency.GetId()
	moduleInfo := strings.Split(moduleId, ":")
	moduleInfo[0] = replaceExclamationMarkWithUpperCase(moduleInfo[0])
	moduleId = strings.Join(moduleInfo, ":")
	modulePath := strings.Replace(moduleId, ":", "@", 1)
	path := filepath.Join(tempDir, modulePath, "go.mod")
	return path
}

func (pwd *PackageWithDeps) publishDependencyAndPopulateTransitive(pathToMod, targetRepo string, graphDependencies map[string]bool, cache *golangutil.DependenciesCache, details *config.ArtifactoryDetails) {
	// If the mod is not empty, populate transitive dependencies
	if len(graphDependencies) > 0 {
		sumFileContent, sumFileStat, err := golangutil.GetSumContentAndRemove(filepath.Dir(pathToMod))
		logError(err)
		pwd.setTransitiveDependencies(targetRepo, graphDependencies, cache, details)
		if len(sumFileContent) > 0 && sumFileStat != nil {
			golangutil.RestoreSumFile(filepath.Dir(pathToMod), sumFileContent, sumFileStat)
		}
	}

	published, _ := cache.GetMap()[pwd.Dependency.GetId()]
	if !published && pwd.PatternMatched(pwd.regExp.GetNotEmptyModRegex()) {
		err := pwd.writeModContentToGoCache()
		logError(err)
	}

	// Populate and publish the transitive dependencies.
	if pwd.transitiveDependencies != nil {
		pwd.populateTransitive(targetRepo, cache, details)
	}

	if !published && pwd.shouldRevertToEmptyMod {
		log.Debug("Reverting to the original mod of", pwd.Dependency.GetId())
		editedBy := pwd.regExp.GetEditedByJFrogCli()
		if editedBy.FindString(string(pwd.originalModContent)) == "" {
			pwd.originalModContent = append([]byte(pwd.GoModEditMessage+"\n\n"), pwd.originalModContent...)
		}
		writeModContentToModFile(pathToMod, pwd.originalModContent)
		pwd.Dependency.SetModContent(pwd.originalModContent)
		err := pwd.writeModContentToGoCache()
		logError(err)
	}
	// Publish to Artifactory the dependency if needed.
	if !published {
		err := pwd.prepareAndPublish(targetRepo, cache, details)
		logError(err)
	}

	// Remove from temp folder the dependency.
	err := os.RemoveAll(filepath.Dir(pathToMod))
	if errorutils.CheckError(err) != nil {
		log.Error("Received an error removing dir:", err, filepath.Dir(pathToMod))
	}
}

// Prepare for publishing and publish the dependency to Artifactory
// Mark this dependency as published
func (pwd *PackageWithDeps) prepareAndPublish(targetRepo string, cache *golangutil.DependenciesCache, details *config.ArtifactoryDetails) error {
	err := pwd.Dependency.prepareAndPublish(targetRepo, cache, details)
	cache.GetMap()[pwd.Dependency.GetId()] = true
	return err
}

func populateModWithTidy(path string) error {
	err := os.Chdir(filepath.Dir(path))
	if errorutils.CheckError(err) != nil {
		return err
	}
	log.Debug("Preparing to populate mod", filepath.Dir(path))
	err = removeGoSum(path)
	logError(err)
	// Running go mod tidy command
	err = golangutil.RunGoModTidy()
	if err != nil {
		return err
	}

	return nil
}

func removeGoSum(path string) error {
	// Remove go.sum file to avoid checksum conflicts with the old go.sum
	goSum := filepath.Join(filepath.Dir(path), "go.sum")
	exists, err := fileutils.IsFileExists(goSum, false)
	if err != nil {
		return err
	}
	if exists {
		err = os.Remove(goSum)
		if errorutils.CheckError(err) != nil {
			return err
		}
	}
	return nil
}

func runGoModGraph() (output map[string]bool, err error) {
	// Running go mod graph command
	return golangutil.GetDependenciesGraph()
}

func (pwd *PackageWithDeps) setTransitiveDependencies(targetRepo string, graphDependencies map[string]bool, cache *golangutil.DependenciesCache, details *config.ArtifactoryDetails) {
	var dependencies []PackageWithDeps
	for transitiveDependency := range graphDependencies {
		module := strings.Split(transitiveDependency, "@")
		if len(module) == 2 {
			dependenciesMap := cache.GetMap()
			name := getDependencyName(module[0])
			_, exists := dependenciesMap[name+":"+module[1]]
			if !exists {
				// Check if the dependency is in the local cache.
				dep, err := createDependency(pwd.cachePath, name, module[1])
				logError(err)
				if err != nil {
					continue
				}
				// Check if this dependency exists in Artifactory.
				client := httpclient.NewDefaultHttpClient()
				downloadedFromArtifactory, err := shouldDownloadFromArtifactory(module[0], module[1], targetRepo, details, client)
				logError(err)
				if err != nil {
					continue
				}
				if dep == nil {
					// Dependency is missing in the local cache. Need to download it...
					dep, err = downloadAndCreateDependency(pwd.cachePath, name, module[1], transitiveDependency, targetRepo, downloadedFromArtifactory, details)
					logError(err)
					if err != nil {
						continue
					}
				}

				if dep != nil {
					log.Debug(fmt.Sprintf("Dependency %s has transitive dependency %s", pwd.Dependency.GetId(), dep.GetId()))
					depsWithTrans := &PackageWithDeps{Dependency: dep,
						regExp:           pwd.regExp,
						cachePath:        pwd.cachePath,
						GoModEditMessage: pwd.GoModEditMessage}
					dependencies = append(dependencies, *depsWithTrans)
					dependenciesMap[name+":"+module[1]] = downloadedFromArtifactory
				}
			} else {
				log.Debug("Dependency", transitiveDependency, "has been previously added.")
			}
		}
	}
	pwd.transitiveDependencies = dependencies
}

func (pwd *PackageWithDeps) writeModContentToGoCache() error {
	moduleAndVersion := strings.Split(pwd.Dependency.GetId(), ":")
	pathToModule := strings.Split(moduleAndVersion[0], "/")
	path := filepath.Join(pwd.cachePath, strings.Join(pathToModule, fileutils.GetFileSeparator()), "@v", moduleAndVersion[1]+".mod")
	err := ioutil.WriteFile(path, pwd.Dependency.GetModContent(), 0700)
	return errorutils.CheckError(err)
}

// Runs over the transitive dependencies, populate the mod files of those transitive dependencies
func (pwd *PackageWithDeps) populateTransitive(targetRepo string, cache *golangutil.DependenciesCache, details *config.ArtifactoryDetails) {
	cache.IncrementTotal(len(pwd.transitiveDependencies))
	for _, transitiveDep := range pwd.transitiveDependencies {
		published, _ := cache.GetMap()[transitiveDep.Dependency.GetId()]
		if !published {
			log.Debug("Starting to work on transitive dependency:", transitiveDep.Dependency.GetId())
			transitiveDep.PopulateModAndPublish(targetRepo, cache, details)
		} else {
			log.Debug("The dependency", transitiveDep.Dependency.GetId(), "was already handled")
		}
	}
}
