package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/camronlevanger/gobo/models"
)

// IPackageService is the interface to implement for working with go packages.
type IPackageService interface {
	Install(url string) error
	Checkout(url string, bookmark string) error
	Get(url string, bookmark string) error
	DetermineBookmark(path string) string
	IsATag(path string) (bool, string)
	PathVisited(path string, f os.FileInfo, err error) error
	GetInstalledPackages() []models.Package
	DiffAndUpdatePackages(currentPackages []models.Package) (bool, []models.Package)
}

// PackageService is the struct for this implementation of IPackageService.
type PackageService struct {
	logger    ILogger
	host      models.Host
	gopath    string
	separator string
}

// GetPackageService returns a pointer to an implementation of IPackageService.
func GetPackageService(
	logger ILogger,
	host models.Host,
	gopath string,
	separator string,
) *PackageService {
	var packageService = PackageService{
		logger,
		host,
		gopath,
		separator,
	}

	return &packageService
}

var installedPackages []models.Package

// Install is a wrapper for the `go install` command.
func (packageService *PackageService) Install(path string) error {

	packageService.logger.Info("Running go install " + path + "...")

	app := "go"
	cmdArgs := []string{"install", path}

	cmd := exec.Command(app, cmdArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		packageService.logger.Info("Error running go install: " + err.Error() + ": " + stderr.String())
		return errors.New("Error running go install: " + err.Error() + ": " + stderr.String())
	}

	return nil
}

// Checkout is a wrapper for the `git checkout` command.
func (packageService *PackageService) Checkout(path string, revision string) error {

	packageService.logger.Info("Running git checkout to: " + revision + " on: " + path + "...")

	app := "git"
	cmdArgs := []string{"checkout", revision}

	os.Chdir(packageService.gopath + path)

	cmd := exec.Command(app, cmdArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		packageService.logger.Info("Error running git checkout: " + err.Error() + ": " + stderr.String())
		return errors.New("Error running git checkout: " + err.Error() + ": " + stderr.String())
	}

	return nil
}

// Get is a wrapper for 'go install' which also subsequently checks out the project at specified revision.
func (packageService *PackageService) Get(path string, revision string) error {

	packageService.logger.Info("Running go get " + path + "...")

	app := "go"
	cmdArgs := []string{"get", path}

	cmd := exec.Command(app, cmdArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		packageService.logger.Info("Error running go get: " + err.Error() + ": " + stderr.String())
		return errors.New("Error running go get: " + err.Error() + ": " + stderr.String())
	}

	err = packageService.Checkout(path, revision)

	err = packageService.Install(path)

	return nil
}

// DetermineBookmark takes the path of a git repo and returns the tag or hash that it is currently checked out at.
func (packageService *PackageService) DetermineBookmark(path string) string {
	tag, version := packageService.IsATag(path)

	if tag {
		packageService.logger.Info("Returning tag as bookmark: " + version)

		return version
	}

	cmd := "git"
	cmdArgs := []string{"rev-parse", "HEAD"}

	os.Chdir(path)
	hash, err := exec.Command(cmd, cmdArgs...).Output()
	if err != nil {
		packageService.logger.Info("Error running git rev-parse HEAD: " + err.Error())
	}

	packageService.logger.Info("Bookmark, HEAD is at: " + string(hash))

	return string(hash)
}

// IsATag takes the path of a git repository and returns whether or not the repo is checked out at a tag, and the tag.
func (packageService *PackageService) IsATag(path string) (bool, string) {

	app := "git"
	cmdArgs := []string{"describe", "--exact-match"}

	os.Chdir(path)
	cmd := exec.Command(app, cmdArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		packageService.logger.Info("Error running git describe --exact-match: " + err.Error() + ": " + stderr.String())
		packageService.logger.Info(fmt.Sprintf("Package %s is not checked out at a tag, using HEAD for bookmark.", path))
		return false, out.String()
	}

	return true, out.String()
}

// PathVisited determines whether or not the visited path is a Golang package, and if so creates a models.Package object.
func (packageService *PackageService) PathVisited(path string, f os.FileInfo, err error) error {
	//packageService.logger.Info(fmt.Sprintf("Visited: %s\n", path))
	newPackage := models.Package{}
	if f.IsDir() {
		_, existsErr := os.Stat(path + packageService.separator + ".git")
		if existsErr == nil {
			newPackage.Path = packageService.GetURLFromPath(path)
			newPackage.Origin = packageService.GetURLFromPath(path)
			newPackage.Revision = packageService.DetermineBookmark(path)
			newPackage.RevisionTime = time.Now().String()
			installedPackages = append(installedPackages, newPackage)
		}
	}
	return nil
}

// GetInstalledPackages returns an array of models.Package that exist in the current directory.
func (packageService *PackageService) GetInstalledPackages() []models.Package {
	packageService.logger.Info("Walking source directory to find installed packages.")
	err := filepath.Walk(packageService.gopath+"src", packageService.PathVisited)
	packageService.logger.Info(fmt.Sprintf("filepath.Walk() returned %v\n", err))

	return installedPackages
}

// GetURLFromPath takes the system path of the package and tries to figure out the http address of the repo.
func (packageService *PackageService) GetURLFromPath(path string) string {

	var url string

	splitPath := strings.Split(path, packageService.gopath+"src/")
	if len(splitPath) > 1 {
		url = splitPath[1]
	} else {
		url = "Unable to determine package URL from path: " + path
	}

	return url
}

// DiffAndUpdatePackages takes the current toml packages and compares them to the filesystem and returns a new array
// representing the current state of the environment filesystem.
func (packageService *PackageService) DiffAndUpdatePackages(currentPackages []models.Package) (bool, []models.Package) {

	var updatedPackages []models.Package
	var changesDetected bool

	systemPackages := packageService.GetInstalledPackages()
	// check if current package bookmarks have changed
	for i := 0; i < len(currentPackages); i++ {
		var sysbook string
		sysbook = packageService.DetermineBookmark(packageService.gopath + currentPackages[i].Path)
		if sysbook != currentPackages[i].Revision {
			packageService.logger.Info(sysbook + " does not match " + currentPackages[i].Revision + " at path " + packageService.gopath + currentPackages[i].Path)
			changesDetected = true
			packageService.logger.Info(currentPackages[i].Path + " has been updated on the filesystem.")
			currentPackages[i].Revision = sysbook
			currentPackages[i].RevisionTime = time.Now().String()
		}
		updatedPackages = append(updatedPackages, currentPackages[i])
	}

	packageService.logger.Info(fmt.Sprintf("Changed packages: %+v", updatedPackages))

	var addedPackages []models.Package
	for sys := 0; sys < len(systemPackages); sys++ {
		match := false
		for cur := 0; cur < len(currentPackages); cur++ {
			if systemPackages[sys].Path == currentPackages[cur].Path {
				match = true
				break
			}
		}

		if !match {
			updatedPackages = append(updatedPackages, systemPackages[sys])
			addedPackages = append(addedPackages, systemPackages[sys])
			packageService.logger.Info("New package " + systemPackages[sys].Path + " has been identified.")
		}

	}

	packageService.logger.Info(fmt.Sprintf("Added packages: %+v", addedPackages))

	return changesDetected, updatedPackages

}
