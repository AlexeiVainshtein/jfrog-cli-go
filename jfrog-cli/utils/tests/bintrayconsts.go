package tests

const (
	BintrayRepo1                = "jfrog-cli-tests-repo1"
	BintrayTestRepositoryConfig = "bintray_repository_config.json"

	BintrayUploadTestPackageName = "uploadTestPackage"
	BintrayUploadTestVersion     = "1.2.3"
)

var BintrayExpectedUploadFlatNonRecursive = []PackageSearchResultItem{{
	Repo:    BintrayRepo1,
	Path:    "a1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "507ac63c6b0f650fb6f36b5621e70ebca3b0965c"}, {
	Repo:    BintrayRepo1,
	Path:    "a2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "de2f31d77e2c2b1039a806f21b0c5f3243e45588"}, {
	Repo:    BintrayRepo1,
	Path:    "a3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "29d38faccfe74dee60d0142a716e8ea6fad67b49"}}

var BintrayExpectedUploadFlatNonRecursiveModified = []PackageSearchResultItem{{
	Repo:    BintrayRepo1,
	Path:    "a1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "954cf8f3f75c41f18540bb38460910b4f0074e6f"}, {
	Repo:    BintrayRepo1,
	Path:    "a2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "de2f31d77e2c2b1039a806f21b0c5f3243e45588"}, {
	Repo:    BintrayRepo1,
	Path:    "a3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "29d38faccfe74dee60d0142a716e8ea6fad67b49"}}

var BintrayExpectedUploadNonFlatNonRecursive = []PackageSearchResultItem{{
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/a1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "507ac63c6b0f650fb6f36b5621e70ebca3b0965c"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/a2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "de2f31d77e2c2b1039a806f21b0c5f3243e45588"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/a3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "29d38faccfe74dee60d0142a716e8ea6fad67b49"}}

var BintrayExpectedUploadRegexpPlaceholders = []PackageSearchResultItem{{
	Repo:    BintrayRepo1,
	Path:    "dir/a/a1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "507ac63c6b0f650fb6f36b5621e70ebca3b0965c"}, {
	Repo:    BintrayRepo1,
	Path:    "dir/a/a2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "de2f31d77e2c2b1039a806f21b0c5f3243e45588"}, {
	Repo:    BintrayRepo1,
	Path:    "dir/a/a3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "29d38faccfe74dee60d0142a716e8ea6fad67b49"}, {
	Repo:    BintrayRepo1,
	Path:    "dir/a/b/b1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "954cf8f3f75c41f18540bb38460910b4f0074e6f"}, {
	Repo:    BintrayRepo1,
	Path:    "dir/a/b/b2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "3b60b837e037568856bedc1dd4952d17b3f06972"}, {
	Repo:    BintrayRepo1,
	Path:    "dir/a/b/b3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "ec6420d2b5f708283619b25e68f9ddd351f555fe"}, {
	Repo:    BintrayRepo1,
	Path:    "dir/a/b/c/c1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "063041114949bf19f6fe7508aef639640e7edaac"}, {
	Repo:    BintrayRepo1,
	Path:    "dir/a/b/c/c2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "a4f912be11e7d1d346e34c300e6d4b90e136896e"}, {
	Repo:    BintrayRepo1,
	Path:    "dir/a/b/c/c3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "2d6ee506188db9b816a6bfb79c5df562fc1d8658"}}

var BintrayExpectedUploadRegexpPlaceholdersWithSpecialCharsFile = append(BintrayExpectedUploadRegexpPlaceholders, PackageSearchResultItem{
	Repo:    BintrayRepo1,
	Path:    "dir/a$+~&^a",
	Package: BintrayUploadTestPackageName,
	Name:    "a$+~&^a",
	Version: BintrayUploadTestVersion,
	Sha1:    "507ac63c6b0f650fb6f36b5621e70ebca3b0965c"})

var BintrayExpectedUploadFlatRecursive = []PackageSearchResultItem{{
	Repo:    BintrayRepo1,
	Path:    "a1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "507ac63c6b0f650fb6f36b5621e70ebca3b0965c"}, {
	Repo:    BintrayRepo1,
	Path:    "a2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "de2f31d77e2c2b1039a806f21b0c5f3243e45588"}, {
	Repo:    BintrayRepo1,
	Path:    "a3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "29d38faccfe74dee60d0142a716e8ea6fad67b49"}, {
	Repo:    BintrayRepo1,
	Path:    "b1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "954cf8f3f75c41f18540bb38460910b4f0074e6f"}, {
	Repo:    BintrayRepo1,
	Path:    "b2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "3b60b837e037568856bedc1dd4952d17b3f06972"}, {
	Repo:    BintrayRepo1,
	Path:    "b3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "ec6420d2b5f708283619b25e68f9ddd351f555fe"}, {
	Repo:    BintrayRepo1,
	Path:    "c1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "063041114949bf19f6fe7508aef639640e7edaac"}, {
	Repo:    BintrayRepo1,
	Path:    "c2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "a4f912be11e7d1d346e34c300e6d4b90e136896e"}, {
	Repo:    BintrayRepo1,
	Path:    "c3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "2d6ee506188db9b816a6bfb79c5df562fc1d8658"}}

var BintrayExpectedUploadNonFlatRecursive = []PackageSearchResultItem{{
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/a1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "507ac63c6b0f650fb6f36b5621e70ebca3b0965c"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/a2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "de2f31d77e2c2b1039a806f21b0c5f3243e45588"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/a3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "a3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "29d38faccfe74dee60d0142a716e8ea6fad67b49"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/b/b1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "954cf8f3f75c41f18540bb38460910b4f0074e6f"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/b/b2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "3b60b837e037568856bedc1dd4952d17b3f06972"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/b/b3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "b3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "ec6420d2b5f708283619b25e68f9ddd351f555fe"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/b/c/c1.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c1.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "063041114949bf19f6fe7508aef639640e7edaac"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/b/c/c2.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c2.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "a4f912be11e7d1d346e34c300e6d4b90e136896e"}, {
	Repo:    BintrayRepo1,
	Path:    GetTestResourcesPath()[1:] + "a/b/c/c3.in",
	Package: BintrayUploadTestPackageName,
	Name:    "c3.in",
	Version: BintrayUploadTestVersion,
	Sha1:    "2d6ee506188db9b816a6bfb79c5df562fc1d8658"}}
