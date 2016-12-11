package tests

import (
	"testing"
	"strconv"
	"flag"
	"strings"
	"os"
	"fmt"
	"runtime"
	"github.com/jfrogdev/jfrog-cli-go/artifactory/commands"
	"github.com/jfrogdev/jfrog-cli-go/utils/ioutils"
	"github.com/jfrogdev/jfrog-cli-go/utils/cliutils/log"
)

var PrintSearchResult *bool
var RtUrl *string
var RtUser *string
var RtPassword *string
var RtApiKey *string
var BtUser *string
var BtKey *string
var TestArtifactory *bool
var TestBintray *bool

func init() {
	PrintSearchResult = flag.Bool("printSearchResult", false, "Set to true for printing search results")
	RtUrl = flag.String("rt.url", "http://localhost:8081/artifactory/", "Artifactory url")
	RtUser = flag.String("rt.user", "admin", "Artifactory username")
	RtPassword = flag.String("rt.password", "password", "Artifactory password")
	RtApiKey = flag.String("rt.apikey", "", "Artifactory user API key")
	TestArtifactory = flag.Bool("test.artifactory", true, "Test Artifactory")
	TestBintray = flag.Bool("test.bintray", false, "Test Bintray")
	BtUser = flag.String("bt.user", "", "Bintray username")
	BtKey = flag.String("bt.key", "", "Bintray password")
}

func CleanFileSystem() {
	isExist, err := ioutils.IsDirExists(Out)
	if err != nil {
		log.Error(err)
	}
	if isExist {
		os.RemoveAll(Out)
	}
}

func IsExistLocally(expected, actual []string, t *testing.T) {
	if len(actual) == 0 && len(expected) != 0 {
		t.Error("Couldn't find all expected files, expected: " + strconv.Itoa(len(expected)) + ", found: " + strconv.Itoa(len(actual)))
	}
	for _, v := range expected {
		for i, r := range actual {
			if v == r {
				break
			}
			if i == len(actual) - 1 {
				t.Error("Missing file : " + v)
			}
		}
	}
}

func CompareExpectedVsActuals(expected []string, actual []commands.SearchResult, t *testing.T) {
	if len(actual) != len(expected) {
		t.Error("Couldn't find all expected files, expected: " + strconv.Itoa(len(expected)) + ", found: " + strconv.Itoa(len(actual)))
	}
	for _, v := range expected {
		for i, r := range actual {
			if v == r.Path {
				break
			}
			if i == len(actual) - 1 {
				t.Error("Missing file: " + v)
			}
		}
	}
}

func GetTestResourcesPath() string {
	dir, _ := os.Getwd()
	fileSepatatr := ioutils.GetFileSeperator()
	index := strings.LastIndex(dir, fileSepatatr)
	dir = dir[:index]
	return dir + ioutils.GetFileSeperator() + "testsdata" + ioutils.GetFileSeperator()
}

func getFileByOs(fileName string) string {
	var currentOs string;
	fileSepatatr := ioutils.GetFileSeperator()
	if runtime.GOOS == "windows" {
		currentOs = "win"
	} else {
		currentOs = "unix"
	}
	return GetTestResourcesPath() + "specs" + fileSepatatr + currentOs + fileSepatatr + fileName
}

func GetFilePath(fileName string) string {
	filePath := GetTestResourcesPath() + "specs/common" + ioutils.GetFileSeperator() + fileName
	isExists, _ := ioutils.IsFileExists(filePath)
	if isExists {
		return filePath
	}
	return getFileByOs(fileName)
}

func FixWinPath(filePath string) string {
	fixedPath := strings.Replace(filePath, "\\", "\\\\", -1)
	return fixedPath
}

type PackageSearchResultItem struct {
	Name      string
	Path      string
	Package   string
	Version   string
	Repo      string
	Owner     string
	Created   string
	Size      int64
	Sha1      string
	Published bool
}

type JfrogCli struct {
	main   func()
	prefix string
	suffix string
}

func NewJfrogCli(mainFunc func(), prefix, suffix string) (*JfrogCli) {
	return &JfrogCli{mainFunc, prefix, suffix}
}

func (cli *JfrogCli) Exec(args ...string) {
	spaceSplit := " "
	os.Args = strings.Split(cli.prefix, spaceSplit)
	for _, v := range args {
		if v == "" {
			continue
		}
		os.Args = append(os.Args, strings.Split(v, spaceSplit)...)
	}
	if cli.suffix != "" {
		os.Args = append(os.Args, strings.Split(cli.suffix, spaceSplit)...)
	}

	fmt.Println("[Command]", strings.Join(os.Args, " "))
	cli.main()
}

func (cli *JfrogCli) WithSuffix(suffix string) *JfrogCli {
	return &JfrogCli{cli.main, cli.prefix, suffix}
}