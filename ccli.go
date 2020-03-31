package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

var Trace *log.Logger // custom logger for multiwriting

type Person struct {
	FirstName   string
	LastName    string
	City        string
	Address     string
	Email       string
	PhoneNumber int64
}

type personSorter struct {
	p  []Person
	by func(p1, p2 *Person) bool
}

type By func(p1, p2 *Person) bool

// Sort methods that need to be defined for sorting Person structs
func (by By) Sort(persons []Person) {
	sortMethod := &personSorter{
		p:  persons,
		by: by,
	}
	sort.Sort(sortMethod)
}

func (p *personSorter) Len() int {
	return len(p.p)
}

func (p *personSorter) Less(i, j int) bool {
	return p.by(&p.p[i], &p.p[j])
}

func (p *personSorter) Swap(i, j int) {
	p.p[i], p.p[j] = p.p[j], p.p[i]
}

// band-aid function until specific Usage() functions are written for flags
func displayusage() string {
	return string("usage: ccli [options] <command> <subcommand(s)> [parameters]\nFor help try:\n\nccli help\nccli <command> help\nccli <command> <subcommand> help")
}

// band-aid function until specific Usage() functions are written for flags
func displayErrorUsage(s string) string {
	return string("usage: of command [" + s + "] is not supported\nFor help try:\n\nccli help\nccli <command> help\nccli <command> <subcommand> help")
}

func convertJSONToContact(json []Person, p *[]Person) {
	*p = json
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

/**
 * Description: 	file loader
 * Purpose: 		loads a list of contact when the cli is called
**/
func loadFile(fname string) []Person {
	jsonFile, err := os.Open(fname)

	//	check if opening file throws an error
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer jsonFile.Close()
	var persons []Person
	rawBytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(rawBytes), &persons)

	return persons
}

/**
 * Description: 	file saver
 * Purpose: 		saves list of contacts that is loaded in when the cli is called
**/
func save(contacts []Person, savePath string) {

	data, _ := json.Marshal(contacts)

	// writes file to path provided along with []byte{}
	_ = ioutil.WriteFile(savePath, data, 0644)
}

func fileExists(path string) bool {
	src, err := os.Stat(path)
	if os.IsNotExist(err) && src == nil {
		return false
	}
	return true
}

/**
 * Description: 	path checker
 * Purpose: 		determines if the path provide exits
**/
func dirExists(dname string, homedir string) bool {
	path := homedir + "/" + dname
	// log.Println("Path:", path)
	src, err := os.Stat(path)

	if os.IsNotExist(err) && src == nil {
		// create new dir
		newDir := os.MkdirAll(path, 0755)
		// log.Println("newDir:", newDir)
		if newDir != nil {
			panic(err)
		}
		return true
	}
	return true
}

/**
 * Description: 	sorting method
 * Purpose: 		sorts list of contacts dependent upon user input
**/
func sortContactsBy(method string, list []Person) {

	// normalize method
	method = strings.ToLower(method)

	whitelist := []string{"a", "ascend", "ascending", "d", "descend", "descending"}
	_, found := Find(whitelist, method)

	if !found {
		flag.Usage()
		os.Exit(1)
	}
	// anon function for Sort method to sort by ascedning order
	acsend := func(p1, p2 *Person) bool {
		return p1.FirstName < p2.FirstName
	}
	// anon function for Sort method to sort by ascedning order
	descend := func(p1, p2 *Person) bool {
		return p1.FirstName > p2.FirstName
	}

	if method[0:1] == "a" {
		printPerson(list)
		By(acsend).Sort(list)
		log.Println()
		printPerson(list)
	} else {
		printPerson(list)
		By(descend).Sort(list)
		log.Println()
		printPerson(list)
	}
}

// used for debugging
func printPerson(p []Person) {
	for _, person := range p {
		log.Println(person)
	}
}

/**
 * Description: 	initializes custom logger
 * Purpose: 		creates a multistream logger for redundency
**/
func InitLogger() {
	t := time.Now()
	d := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	fpath := "./.logs/log-" + d.String()
	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file")
	}

	multilogger := io.MultiWriter(file, os.Stdout)
	// config logger
	log.SetOutput(multilogger)
	log.SetFlags(0)
}

func main() {
	InitLogger()
	fileName := "contacts.json"
	dataDirName := ".ccli/"
	homePath := os.Getenv("HOME")
	var contactList []Person // holds person's info

	// checks if home path exists
	if homePath == "" {
		log.Fatal("Error, no envvar for $HOME")
		os.Exit(1)
	}

	log.Printf("homePath: %s\n", homePath)
	log.Println("dirExists:", dirExists(dataDirName, string(homePath)))

	if dirExists(dataDirName, string(homePath)) {
		// check for data file to load
		pathToFile := homePath + "/" + dataDirName + fileName
		// log.Println("fileExists:", fileExists(pathToFile))
		if fileExists(pathToFile) {
			data := loadFile(pathToFile)
			if data != nil {
				convertJSONToContact(data, &contactList)
			}
			log.Println("contactList:", contactList)

		}
	}

	// sub commands for cli
	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	// flags for create command
	createFirstNamePtr := createCommand.String("fname", "0", "-fname - field value for first name (REQUIRED])")
	createLastNamePtr := createCommand.String("lname", "0", "-lname - sets value for last name (REQUIRED)")
	createCityPtr := createCommand.String("city", "0", "-c - sets value for city name")
	createAddrPtr := createCommand.String("addr", "0", "-addr - sets value for address")
	createEmailPtr := createCommand.String("email", "0", "-email - sets value for email (REQUIRED)")
	createPhonePtr := createCommand.Int64("pnum", -1, "-pnum - sets value for phoneNumber (OPTIONAL)")

	// flags for list command
	listOrderPtr := listCommand.String("order", "0", "-order - lets you choose how to display your contacts")

	// exists if sub command wasn't sepcificed
	if len(os.Args) < 2 {
		log.Println(displayusage())
		os.Exit(1)
	}

	flag.Parse()

	//switch statement to for what command and sub commands to run
	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	case "create":
		createCommand.Parse(os.Args[2:])
	case "help":
		helpCommand.Parse(os.Args[2:])
	default:
		log.Println(displayErrorUsage(string(os.Args[1])))
		os.Exit(1)
	}

	if createCommand.Parsed() {
		// verify parsed content can execute
		if *createFirstNamePtr == "0" {
			flag.Usage()
			os.Exit(1)
		}
		if *createLastNamePtr == "0" {
			flag.Usage()
			os.Exit(1)
		}
		if *createEmailPtr == "0" {
			flag.Usage()
			os.Exit(1)
		}
		if *createPhonePtr == -1 {
			flag.Usage()
			os.Exit(1)
		}

		p := Person{*createFirstNamePtr, *createLastNamePtr, *createCityPtr, *createAddrPtr, *createEmailPtr, *createPhonePtr}

		contactList = append(contactList, p)

	}
	// List Commmands
	if listCommand.Parsed() {
		log.Println(*listCommand)

		if *listOrderPtr == "0" {
			flag.Usage()
			os.Exit(1)
		}

		sortContactsBy(*listOrderPtr, contactList)

	}
	if helpCommand.Parsed() {
		log.Println(*helpCommand)
	}

	// saves contacts
	savePath := homePath + "/" + dataDirName + fileName
	save(contactList, savePath)

	return

}
