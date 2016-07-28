package models

// GOPATHDIRECTORIES is an array of directories to operate on in the GOPATH.
var GOPATHDIRECTORIES = [...]string{"src", "pkg", "bin"}

// GOPATHFILES is an array of files to operate on in the GOPATH.
var GOPATHFILES = [...]string{"gobo.json"}

// GOBOVERSION is the version of the app.
const GOBOVERSION = "0.0.2"

// FILEMODE is the constant value for the filemode to create directories and files with.
const FILEMODE = 0755

// GOBOSPEED is the ascii art gobo logo.
const GOBOSPEED = "" +
	"              ______          \n" +
	"_______ _________  /_______      \n" +
	"__  __ `/  __ \\_  __ \\  __ \\     \n" +
	"_  /_/ // /_/ /  /_/ / /_/ /     \n" +
	"_\\__, / \\____//_.___/\\____/      \n" +
	"/____/                           \n"
