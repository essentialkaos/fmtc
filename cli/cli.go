package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic utility info
const (
	APP  = "fmtc"
	VER  = "2.0.0"
	DESC = "Utility for rendering fmtc formatted data"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_ERROR = "E:error"
	OPT_LINE  = "L:line"
	OPT_HELP  = "h:help"
	OPT_VER   = "v:version"

	OPT_UPDATE       = "U:update"
	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_ERROR: {Type: options.BOOL},
	OPT_LINE:  {Type: options.BOOL},
	OPT_HELP:  {Type: options.BOOL},
	OPT_VER:   {Type: options.MIXED},

	OPT_UPDATE:       {Type: options.MIXED},
	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// colorTagApp is tag used for app name
var colorTagApp string

// colorTagVer is tag used for app version
var colorTagVer string

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	runtime.GOMAXPROCS(1)

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error(errs.Error(" - "))
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
		support.Collect(APP, VER).WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).Print()
		os.Exit(0)
	case withSelfUpdate && options.GetB(OPT_UPDATE):
		os.Exit(updateBinary())
	case options.GetB(OPT_HELP):
		genUsage().Print()
		os.Exit(0)
	}

	colorData(args)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// configureUI configures user interface
func configureUI() {
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

	// Eval all escape sequences
	data, _ = strconv.Unquote(`"` + strings.ReplaceAll(data, `"`, `\"`) + `"`)

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

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, "fmtc"))
	case "fish":
		fmt.Print(fish.Generate(info, "fmtc"))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, "fmtc"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "data…")

	info.AppNameColorTag = colorTagApp

	info.AddOption(OPT_ERROR, "Print data to stderr")
	info.AddOption(OPT_LINE, "Don't print newline at the end")

	if withSelfUpdate {
		info.AddOption(OPT_UPDATE, "Update application to the latest version")
	}

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

		BugTracker:    "https://github.com/essentialkaos/fmtc/issues",
		UpdateChecker: getUpdateChecker(),
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //
