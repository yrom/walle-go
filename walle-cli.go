package main

import (
	"flag"
	"os"
	"fmt"
	"walle"
	"strings"
	"bytes"
	"path/filepath"
)

type extraInfo map[string]string

// Default value of extraInfo
func (e *extraInfo) String() string {
	return ""
}
func (e *extraInfo) Set(val string) error {
	pairs := strings.Split(val, ",")
	*e = make(extraInfo)
	for _, p := range pairs {
		a := strings.Split(p, "=")
		(*e)[a[0]] = a[1]
	}
	return nil
}

type channels []string

// Default value of channels
func (c *channels) String() string {
	return ""
}

func (c *channels) Set(val string) error {
	*c = strings.Split(val, ",")
	return nil
}

var (
	command     = filepath.Base(os.Args[0])
	show        = flag.NewFlagSet("show", flag.ExitOnError)
	gen         = flag.NewFlagSet("gen", flag.ExitOnError)
	showRaw     bool
	showHelp    bool
	genOut      string
	genChannels channels
	genExtras   extraInfo
	genForce    bool
	genDebug    bool
	genHelp     bool
)

func init() {

	show.BoolVar(&showRaw, "r", false, "print `raw` text associated to id 0x71777777")
	show.BoolVar(&showHelp, "h", false, "print `help` message of show command")
	gen.StringVar(&genOut, "o", "", "`output` dir, generated channel apk(s) will store in here. default is input's dir")
	gen.Var(&genChannels, "c", "generate apk with the `channel(s)`, split multiple channels with ','")
	gen.Var(&genExtras, "e", "generate apk with the `extras` info (key value pairs, e.g thing=test,boom=1)")
	gen.BoolVar(&genHelp, "h", false, "print `help` message of gen command")
	gen.BoolVar(&genForce, "f", false, "`force` to overwrite exist channeled apk in output")
	gen.BoolVar(&genDebug, "d", false, "print `debug` log")
}

// ./walle show xxxx.apk
func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	var cmd = os.Args[1]
	switch cmd {
	case "show":
		show.Parse(os.Args[2:])
		if showHelp {
			printUsageOfShow()
			break
		}
		args := show.Args()
		if showHelp || len(args) == 0 {
			exit("Error: no apk files!")
		}

		if showRaw {
			walle.PrintRaw(args)
		} else {
			walle.PrintChannel(args)
		}
		break
	case "gen":
		gen.Parse(os.Args[2:])
		if genHelp {
			printUsageOfGen()
			break
		}
		args := gen.Args()
		if len(args) == 0 {
			exit("Error: no input file!")
		}
		if len(args) > 1 {
			fmt.Println("Warning: too many input files, only first one will be used!")
		}
		walle.GenerateChannelApk(genOut, genChannels, genExtras, args[0], genForce, genDebug)

		break
	case "help":
		printHelp()
		fmt.Println()
		printUsageOfShow()
		fmt.Println()
		printUsageOfGen()
		break;
	default:
		printHelp()
		break
	}
}
func exit(v string) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
func printUsageOfGen() {
	fmt.Printf("%s  gen [-o out] -c <channels> [-e extras] <file>\n", command)
	gen.VisitAll(printFlag)
	fmt.Println("  e.g gen -c test /foo/bar/A.apk")
	fmt.Println("      gen -o /foo/bar/channel/ -c test /foo/bar/A.apk")
	fmt.Println("      gen -o /foo/bar/channel/ -c test1,test2 /foo/bar/A.apk")
}

func printUsageOfShow() {
	fmt.Printf("%s  show [-r] <files...>\n", command)
	show.VisitAll(printFlag)
	fmt.Println("  e.g show /foo/bar/A.apk /foo/bar/bar/B.apk")
	fmt.Println("      show -r /foo/bar/A.apk")
}

func printFlag(f *flag.Flag) {
	var buf bytes.Buffer
	buf.WriteString("      -")
	buf.WriteString(f.Name)

	name, usage := flag.UnquoteUsage(f)
	if len(name) > 0 {
		buf.WriteString("  ")
		buf.WriteString(name)
	}
	buf.WriteString("\n  \t")
	buf.WriteString(usage)
	fmt.Println(buf.String())
}
func printHelp() {
	fmt.Println("walle-cli is a command line tool for processing channel apk")
	fmt.Printf("Usage:\n  %s command [args]\n", command)
	fmt.Println("Commands")
	fmt.Println("  show \tget channel info from apk and show all by default")
	fmt.Println("  gen \tgenerate apk with channel info")
	fmt.Println("  help \tprint help message")
	fmt.Println()
	fmt.Printf("%s <command> -h for more useful info\n", command)
}
