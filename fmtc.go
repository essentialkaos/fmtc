package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io"
	"os"
	"runtime"

	_ "embed"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"

	"github.com/essentialkaos/fmtc/support"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic utility info
const (
	APP  = "fmtc"
	VER  = "0.0.1"
	DESC = "Utility for rendering fmtc formatted data"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_ERROR    = "E:error"
	OPT_LINE     = "L:line"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_ERROR:    {Type: options.BOOL},
	OPT_LINE:     {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

//go:embed go.mod
var gomod []byte

// gitrev is short hash of the latest git commit
var gitRev string

// ////////////////////////////////////////////////////////////////////////////////// //

// colorTagApp is tag used for app name
var colorTagApp string

// colorTagVer is tag used for app version
var colorTagVer string

// ////////////////////////////////////////////////////////////////////////////////// //

// main is main utility function
func main() {
	runtime.GOMAXPROCS(1)

	preConfigureUI()

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Print(APP, VER, gitRev, gomod)
		os.Exit(0)
	case options.GetB(OPT_HELP):
		genUsage().Print()
		os.Exit(0)
	}

	colorData(args)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if os.Getenv("NO_COLOR") != "" {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{&}{#FF1D7C}", "{#FF1D7C}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{&}{#197}", "{#197}"
	default:
		colorTagApp, colorTagVer = "{*}{&}{m}", "{m}"
	}
}

// colorData processes input data
func colorData(args options.Arguments) {
	var data string
	var isStdin bool

	if len(args) != 0 {
		data = args.Flatten()
	} else {
		if !fsutil.IsCharacterDevice("/dev/stdin") {
			stdinData, err := io.ReadAll(os.Stdin)

			if err != nil {
				os.Exit(1)
			}

			data = string(stdinData)
			isStdin = true
		} else {
			os.Exit(1)
		}
	}

	if options.GetB(OPT_LINE) || isStdin {
		if options.GetB(OPT_ERROR) {
			fmtc.Fprint(os.Stderr, data)
		} else {
			fmtc.Fprint(os.Stdout, data)
		}
	} else {
		if options.GetB(OPT_ERROR) {
			fmtc.Fprintln(os.Stderr, data)
		} else {
			fmtc.Fprintln(os.Stdout, data)
		}
	}
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	if len(a) == 0 {
		fmtc.Fprintln(os.Stderr, "{r}"+f+"{!}")
	} else {
		fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "fmtc"))
	case "fish":
		fmt.Printf(fish.Generate(info, "fmtc"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "fmtc"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "data…")

	info.AppNameColorTag = colorTagApp

	info.AddOption(OPT_ERROR, "Print data to stderr")
	info.AddOption(OPT_LINE, "Don't print newline at the end")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		`"{*}Done!{!} File {#87}$file{!} successfully uploaded to {g_}$host{!}"`,
		"Print fmtc formatted message",
	)

	info.AddExample(
		`-E "{r*}There is no user bob{!}"`,
		"Print fmtc formatted message to stderr",
	)

	info.AddExample(
		`-nc "{*}Done!{!} File {#87}$file{!} successfully uploaded to {g_}$host{!}"`,
		"Print message without colors using -nc/--no-color option",
	)

	info.AddRawExample(
		`echo "{*}Done!{!} File {#87}$file{!} successfully uploaded to {g_}$host{!}" | fmtc`,
		"Use stdin as a source of data",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,

		DescSeparator: "—",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //
