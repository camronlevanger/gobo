package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/camronlevanger/gobo/models"
	"github.com/camronlevanger/gobo/utils"
)

// IListCommand is the interface to implement for listing a directory.
type IListCommand interface {
	Run(dir string)
}

// ListCommand is the struct for this implementation of IListCommand.
type ListCommand struct {
	logger         utils.ILogger
	configService  utils.IConfigService
	packageService utils.IPackageService
	copyService    utils.ICopyService
	moveService    utils.IMoveService
	host           models.Host
	gopath         string
	gobopath       string
}

// GetListCommand returns a pointer to an implementation of IListCommand.
func GetListCommand(
	logger utils.ILogger,
	configService utils.IConfigService,
	packageService utils.IPackageService,
	copyService utils.ICopyService,
	moveService utils.IMoveService,
	host models.Host,
	gopath string,
	gobopath string,
) *ListCommand {
	list := ListCommand{
		logger,
		configService,
		packageService,
		copyService,
		moveService,
		host,
		gopath,
		gobopath,
	}

	return &list
}

// Run lists all directories in the given filepath, and gives the option to activate them.
func (list *ListCommand) Run(dir string) {

	fmt.Println("Available environments:")
	var envs []string

	files, _ := ioutil.ReadDir(dir)
	count := 0
	for _, f := range files {
		if f.IsDir() && f.Name() != "initial" {
			count++
			fmt.Println(strconv.Itoa(count) + ". " + f.Name())
			envs = append(envs, f.Name())
		}
	}

	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the environment number to activate (Enter to cancel): ")
	answer, _ := reader.ReadString('\n')
	if answer == "\n" {
		list.logger.Info("Not saving environment updates.")
		fmt.Println("Goodbye.")
		return
	}

	selection := strings.Split(answer, "\n")
	list.logger.Info("Chose to activate option: " + selection[0])

	activate := GetActivateCommand(
		list.logger,
		list.configService,
		list.packageService,
		list.copyService,
		list.moveService,
		list.host,
		list.gopath,
		list.gobopath,
	)

	num, _ := strconv.Atoi(selection[0])
	name := envs[num-1]

	list.logger.Info("Option " + selection[0] + " is " + name + ", activating...")

	err := activate.Run(name)
	if err != nil {
		list.logger.Fatal("Error running gobo activate command: " + err.Error())
	}
}
