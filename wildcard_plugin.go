package main

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/cli/cf/formatters"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/jeaniejung/Wildcard_Plugin/table"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Wildcard struct {
}

func (cmd *Wildcard) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "wildcard_plugin",
		Version: plugin.VersionType{
			Major: 2,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "wildcard-apps",
				Alias:    "wc-a",
				HelpText: "List all apps in the target space matching the wildcard pattern",
				UsageDetails: plugin.Usage{
					Usage: "cf wildcard-apps APP_NAME_WITH_WILDCARD",
				},
			},
			{
				Name:     "wildcard-delete",
				Alias:    "wc-d",
				HelpText: "Delete apps in the target space matching the wildcard pattern",
				UsageDetails: plugin.Usage{
					Usage: "cf wildcard-delete APP_NAME_WITH_WILDCARD [-f -r]",
				},
			},
			{
				Name:     "wildcard-unbind",
				Alias:    "wc-us",
				HelpText: "Unbind service on apps in the target space matching the wildcard pattern",
				UsageDetails: plugin.Usage{
					Usage: "cf wildcard-unbind APP_NAME_WITH_WILDCARD [-f]",
				},
			},
		},
	}
}

func main() {
	plugin.Start(newWildcard())
}

func newWildcard() *Wildcard {
	return &Wildcard{}
}

func reverseOrder(s []string) []string {
	reversedString := make([]string, len(s))
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		reversedString[i], reversedString[j] = s[j], s[i]
	}
	return reversedString
}

func (cmd *Wildcard) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "wildcard-apps" && len(args) == 2 {
		cmd.WildcardCommandApps(cliConnection, args[1])
	} else if args[0] == "wildcard-delete" && len(args) >= 2 && len(args) <= 4 {
		pattern := args[1]
		fmt.Println(pattern)
		var _force bool = false
		var _routes = false
        var force* bool = &_force
        var routes* bool = &_routes
		if len(args) > 2 {
			wildcardFlagSet := flag.NewFlagSet("echo", flag.ExitOnError)
			force := wildcardFlagSet.Bool("f", false, "forces deletion of all apps matching APP_NAME_WITH_WILDCARD")
			routes := wildcardFlagSet.Bool("r", false, "delete routes asssociated with APP_NAME_WITH_WILDCARD")
			err := wildcardFlagSet.Parse(args[1:])
			fmt.Println(args)
			fmt.Println(wildcardFlagSet)

			reversedArgs := reverseOrder(args)

			fmt.Println(reversedArgs)

			wildcardFlagSet2 := flag.NewFlagSet("echo", flag.ExitOnError)
			force2 := wildcardFlagSet2.Bool("f", false, "forces deletion of all apps matching APP_NAME_WITH_WILDCARD")
			routes2 := wildcardFlagSet2.Bool("r", false, "delete routes asssociated with APP_NAME_WITH_WILDCARD")
			err2 := wildcardFlagSet2.Parse(reversedArgs)
			fmt.Println(reversedArgs)
			fmt.Println(wildcardFlagSet2)

			pattern = wildcardFlagSet.Arg(0)
			*force = *force || *force2
			*routes = *routes || *routes2

			// Parse starting from [1] because the [0]th element is the
			// name of the command and
			//wildcardFlagSet's parsing begins with 1 (*asd -f -r) and args left with unchanged
			//and wildcardFlagSet.Args() returns (*asd)
			checkError(err)
			checkError(err2)
		}
		cmd.WildcardCommandDelete(cliConnection, pattern, force, routes)
	} else if args[0] == "wildcard-unbind" && len(args) >= 3 && len(args) <= 4 {
		service := args[1]
		fmt.Println(service)
		pattern := args[2]
		fmt.Println(pattern)
		var _force bool = false
        var force* bool = &_force
		if len(args) > 3 {
			wildcardFlagSet := flag.NewFlagSet("echo", flag.ExitOnError)
			force := wildcardFlagSet.Bool("f", false, "forces unbind of all apps matching APP_NAME_WITH_WILDCARD")
			err := wildcardFlagSet.Parse(args[1:])
			fmt.Println(args)
			fmt.Println(wildcardFlagSet)

			reversedArgs := reverseOrder(args)

			fmt.Println(reversedArgs)

			wildcardFlagSet2 := flag.NewFlagSet("echo", flag.ExitOnError)
			force2 := wildcardFlagSet2.Bool("f", false, "forces unbind of all apps matching APP_NAME_WITH_WILDCARD")
			err2 := wildcardFlagSet2.Parse(reversedArgs)
			fmt.Println(reversedArgs)
			fmt.Println(wildcardFlagSet2)

			pattern = wildcardFlagSet.Arg(0)
			*force = *force || *force2

			// Parse starting from [1] because the [0]th element is the
			// name of the command and
			//wildcardFlagSet's parsing begins with 1 (*asd -f) and args left with unchanged
			//and wildcardFlagSet.Args() returns (*asd)
			checkError(err)
			checkError(err2)
		}
		cmd.WildcardCommandUnbind(cliConnection, service, pattern, force)
	} else {
		usage(args)
	}
}

func usage(args []string) {
	if args[0] == "wildcard-apps" {
		fmt.Println("Usage: cf wildcard-apps\n\tcf wildcard-apps APP_NAME_WITH_WILDCARD")
	} else if args[0] == "wildcard-delete" {
		fmt.Println("Usage: cf wildcard-delete\n\tcf wildcard-delete APP_NAME_WITH_WILDCARD [-f -r]")
	} else if args[0] == "wildcard-unbind" {
		fmt.Println("Usage: cf wildcard-unbind\n\tcf wildcard-unbind APP_NAME_WITH_WILDCARD [-f]")
	}
}

func checkError(err error) {
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
}

func introduction(cliConnection plugin.CliConnection, args string) {
	currOrg, _ := cliConnection.GetCurrentOrg()
	currSpace, _ := cliConnection.GetCurrentSpace()
	currUsername, _ := cliConnection.Username()
	fmt.Println("Getting apps matching", table.EntityNameColor(args), "in org", table.EntityNameColor(currOrg.Name), "/ space", table.EntityNameColor(currSpace.Name), "as", table.EntityNameColor(currUsername))
	fmt.Println(table.SuccessColor("OK"))
	fmt.Println("")
}

func getMatchedApps(cliConnection plugin.CliConnection, args string) []plugin_models.GetAppsModel {
	pattern := args
	output, err := cliConnection.GetApps()
	checkError(err)
	matchedApps := []plugin_models.GetAppsModel{}
	for i := 0; i < (len(output)); i++ {
		ok, _ := filepath.Match(pattern, output[i].Name)
		if ok {
			matchedApps = append(matchedApps, output[i])
		}
	}
	return matchedApps
}

func (cmd *Wildcard) WildcardCommandApps(cliConnection plugin.CliConnection, args string) {
	introduction(cliConnection, args)
	output := getMatchedApps(cliConnection, args)
	mytable := table.NewTable([]string{("name"), ("requested state"), ("instances"), ("memory"), ("disk"), ("urls")})
	for _, app := range output {
		var urls []string
		for _, route := range app.Routes {
			if route.Host == "" {
				urls = append(urls, route.Domain.Name)
			}
			urls = append(urls, fmt.Sprintf("%s.%s", route.Host, route.Domain.Name))
		}
		runningInstances := strconv.Itoa(app.RunningInstances)
		if app.RunningInstances < 0 {
			runningInstances = "?"
		}
		mytable.Add(
			app.Name,
			app.State,
			runningInstances+"/"+strconv.Itoa(app.TotalInstances),
			formatters.ByteSize(app.Memory*formatters.MEGABYTE),
			formatters.ByteSize(app.DiskQuota*formatters.MEGABYTE),
			strings.Join(urls, ", "),
		)
	}
	mytable.Print()
	if len(output) == 0 {
		fmt.Println(table.WarningColor("No apps found matching"), table.WarningColor(args))
	}
}

func (cmd *Wildcard) WildcardCommandDelete(cliConnection plugin.CliConnection, args string, force *bool, routes *bool) {
	output := getMatchedApps(cliConnection, args)
	exit := false
	if !*force && len(output) > 0 {
		cmd.WildcardCommandApps(cliConnection, args)
		fmt.Println("")
		fmt.Printf("Would you like to delete the apps (%s)nteractively, (%s)ll, or (%s)ancel this command?%s", table.PromptColor("i"), table.PromptColor("a"), table.PromptColor("c"), table.PromptColor(">"))
		var mode string
		fmt.Scanf("%s", &mode)
		if strings.EqualFold(mode, "a") || strings.EqualFold(mode, "all") {
			*force = true
		} else if strings.EqualFold(mode, "i") || strings.EqualFold(mode, "interactively") {
		} else {
			fmt.Println(table.WarningColor("Delete cancelled"))
			exit = true
		}
	} else {
		introduction(cliConnection, args)
	}
	if !exit {
		for _, app := range output {
			coloredAppName := table.EntityNameColor(app.Name)
			if *force && *routes {
				cliConnection.CliCommandWithoutTerminalOutput("delete", app.Name, "-f", "-r")
				fmt.Println("Deleting app", coloredAppName, "and its mapped routes")
			} else if *force {
				cliConnection.CliCommandWithoutTerminalOutput("delete", app.Name, "-f")
				fmt.Println("Deleting app", coloredAppName)
			} else {
				var confirmation string
				fmt.Printf("Really delete the app %s?%s ", table.PromptColor(app.Name), table.PromptColor(">"))
				fmt.Scanf("%s", &confirmation)
				if strings.EqualFold(confirmation, "y") || strings.EqualFold(confirmation, "yes") {
					if *routes {
						cliConnection.CliCommandWithoutTerminalOutput("delete", app.Name, "-f", "-r")
						fmt.Println("Deleting app", coloredAppName, "and its mapped routes")
					} else {
						cliConnection.CliCommandWithoutTerminalOutput("delete", app.Name, "-f")
						fmt.Println("Deleting app", coloredAppName)
					}
				}
			}
		}
	}
	if len(output) == 0 {
		fmt.Println(table.WarningColor("No apps found matching"), table.WarningColor(args))
	} else {
		fmt.Println(table.SuccessColor("OK"))
	}
}

func (cmd *Wildcard) WildcardCommandUnbind(cliConnection plugin.CliConnection, service string, args string, force *bool) {
	output := getMatchedApps(cliConnection, args)
	exit := false
	if !*force && len(output) > 0 {
		cmd.WildcardCommandApps(cliConnection, args)
		fmt.Println("")
		fmt.Printf("Would you like to run this command (%s)nteractively, (%s)ll, or (%s)ancel ?%s", table.PromptColor("i"), table.PromptColor("a"), table.PromptColor("c"), table.PromptColor(">"))
		var mode string
		fmt.Scanf("%s", &mode)
		if strings.EqualFold(mode, "a") || strings.EqualFold(mode, "all") {
			*force = true
		} else if strings.EqualFold(mode, "i") || strings.EqualFold(mode, "interactively") {
		} else {
			fmt.Println(table.WarningColor("Cancelled"))
			exit = true
		}
	} else {
		introduction(cliConnection, args)
	}
	if !exit {
		for _, app := range output {
			coloredService := table.EntityNameColor(service)
			coloredAppName := table.EntityNameColor(app.Name)
			if *force {
				cliConnection.CliCommandWithoutTerminalOutput("unbind-service", app.Name, service, "-f")
				fmt.Println("Unbinding service", coloredService, "from app", coloredAppName)
			} else {
				var confirmation string
				fmt.Printf("Really unbind service %s from the app %s?%s ", table.PromptColor(service), table.PromptColor(app.Name), table.PromptColor(">"))
				fmt.Scanf("%s", &confirmation)
				if strings.EqualFold(confirmation, "y") || strings.EqualFold(confirmation, "yes") {
        			cliConnection.CliCommandWithoutTerminalOutput("unbind-service", app.Name, service, "-f")
					fmt.Println("Unbinding service", coloredService, "from app", coloredAppName)
				}
			}
		}
	}
	if len(output) == 0 {
		fmt.Println(table.WarningColor("No apps found matching"), table.WarningColor(args))
	} else {
		fmt.Println(table.SuccessColor("OK"))
	}
}
