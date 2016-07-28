package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/camronlevanger/go-homedir"
	"github.com/camronlevanger/gobo/commands"
	"github.com/camronlevanger/gobo/models"
	"github.com/camronlevanger/gobo/utils"
)

func main() {

	var verbose bool
	var populate bool
	var file string
	var separator string
	var gopath string
	var home string
	var gobo string
	var goboInitial string
	var goboMaster string
	var initial bool

	separator = string(filepath.Separator)

	var err error

	home, err = homedir.Dir()
	if err != nil {
		FailOnError(err, "Unable to determine home directory")
	}

	gobo = home + separator + ".gobo" + separator
	goboInitial = gobo + "initial" + separator
	goboMaster = gobo + "gobo_master.toml"
	gopath = os.Getenv("GOPATH") + separator

	// print a gobo logo
	fmt.Println(models.GOBOSPEED + "\n")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    gobo create|delete|activate|save|get|list|install|tools [name] [args] ...\n")
		flag.PrintDefaults()
	}

	flag.BoolVar(
		&populate,
		"p",
		false,
		"Populate the new environment with everything that exists in the current environment.",
	)

	flag.BoolVar(&verbose, "v", false, "Print debug info to the console as it happens.")

	flag.StringVar(&file, "f", "gobo.toml", "The gobo toml file to install from.")

	flag.Parse()

	logger := utils.GetLogger(verbose)

	logger.Info(fmt.Sprintf("GOPATH interpreted as %s", gopath))
	logger.Info(fmt.Sprintf("populate: %v", populate))
	logger.Info(fmt.Sprintf("verbose: %v", verbose))
	logger.Info(fmt.Sprintf("Running cmd: %s", flag.Arg(0)))
	logger.Info(fmt.Sprintf("Virtual Environment Name: %s", flag.Arg(1)))

	backup := commands.GetBackupCommand(
		logger,
		utils.GetCopyService(),
		home,
		gopath,
		gobo,
		goboInitial,
		goboMaster,
	)

	initial, err = backup.Run()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating initial backup of GOPATH: %v\n", err))
	}

	command := flag.Arg(0)
	name := flag.Arg(1)

	logger.Info("Switching on: " + command)

	switch command {
	case "create":

		logger.Info("Creating new virtual environment: " + flag.Arg(1))

		configService := utils.GetConfigService(logger)
		copyService := utils.GetCopyService()
		packageService := utils.GetPackageService(logger, getHostInfo(), gopath, separator)
		moveService := utils.GetMoveService()

		create := commands.GetCreateCommand(
			logger,
			configService,
			copyService,
			moveService,
			packageService,
			populate,
			gopath,
			gobo,
			getHostInfo(),
			initial,
		)

		err := create.Run(name)
		if err != nil {
			logger.Fatal("Error running gobo create command: " + err.Error())
		}

		fmt.Println("Create command complete.")

	case "save":
		configService := utils.GetConfigService(logger)
		copyService := utils.GetCopyService()
		packageService := utils.GetPackageService(logger, getHostInfo(), gopath, separator)

		save := commands.GetSaveCommand(
			logger,
			configService,
			packageService,
			copyService,
			getHostInfo(),
			gopath,
			gobo,
		)

		err := save.Run(true)
		if err != nil {
			logger.Fatal("Error running gobo save command: " + err.Error())
		}

		fmt.Println("Save command complete.")

	case "activate":
		configService := utils.GetConfigService(logger)
		copyService := utils.GetCopyService()
		packageService := utils.GetPackageService(logger, getHostInfo(), gopath, separator)
		moveService := utils.GetMoveService()

		activate := commands.GetActivateCommand(
			logger,
			configService,
			packageService,
			copyService,
			moveService,
			getHostInfo(),
			gopath,
			gobo,
		)

		err := activate.Run(name)
		if err != nil {
			logger.Fatal("Error running gobo activate command: " + err.Error())
		}

		fmt.Println("Activate command complete.")

	case "restore":
		copyService := utils.GetCopyService()

		restore := commands.GetRestoreCommand(
			logger,
			copyService,
			gopath,
			gobo,
		)

		err := restore.Run()
		if err != nil {
			logger.Fatal("Error running gobo restore command: " + err.Error())
		}

		fmt.Println("Restore command complete.")

	case "list":

		configService := utils.GetConfigService(logger)
		copyService := utils.GetCopyService()
		packageService := utils.GetPackageService(logger, getHostInfo(), gopath, separator)
		moveService := utils.GetMoveService()

		list := commands.GetListCommand(
			logger,
			configService,
			packageService,
			copyService,
			moveService,
			getHostInfo(),
			gopath,
			gobo,
		)
		list.Run(gobo)

	case "delete":
		configService := utils.GetConfigService(logger)
		moveService := utils.GetMoveService()

		delete := commands.GetDeleteCommand(
			logger,
			moveService,
			configService,
			gobo,
			gopath,
		)
		err := delete.Run(name)
		if err != nil {
			logger.Fatal("Error running gobo delete command: " + err.Error())
		}

		fmt.Println("Delete command complete.")

	case "version":
		version := commands.GetVersionCommand()

		version.Run()

	case "install":
		configService := utils.GetConfigService(logger)
		copyService := utils.GetCopyService()
		packageService := utils.GetPackageService(logger, getHostInfo(), gopath, separator)

		install := commands.GetInstallCommand(
			logger,
			configService,
			packageService,
			copyService,
			getHostInfo(),
			gopath,
			gobo,
		)

		err := install.Run(file)
		if err != nil {
			logger.Fatal("Error running gobo install command: " + err.Error())
		}

		fmt.Println("Install command complete.")

	default:
		fmt.Println(command + " is not a known gobo command: ")
		flag.Usage()
	}

}

// FailOnError is the function to be called on fatal errors, this kills the app.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func getHostInfo() models.Host {
	host := models.Host{}
	version := runtime.Version()

	userInfo, err := user.Current()
	if err != nil {
		userInfo = nil
	}

	host.OS = runtime.GOOS
	host.Version = version
	host.User = *userInfo

	return host
}
