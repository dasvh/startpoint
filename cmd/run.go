/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	requestchain "startpoint/core/chaining"
	"startpoint/core/client/runner"
	"startpoint/core/loader"
	"startpoint/core/model"
	"startpoint/core/print"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RunConfig struct {
	Plain        bool
	PrintHeaders bool
	PrintBody    bool
}

type RunArgs struct {
	Request string
	Profile string
}

var runConfig RunConfig

var runCmd = &cobra.Command{
	Use:   "run [REQUEST NAME] [PROFILE NAME]",
	Short: "Run a http request",
	Long:  `Run a http request`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(0, 2)(cmd, args); err != nil {
			return err
		}

		if len(args) == 0 {
			return nil
		}

		parsedArgs := ParseArgs(args)

		if len(parsedArgs.Request) == 0 {
			return errors.New("Request name is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		runArgs := ParseArgs(args)

		requests, err := loader.ReadRequests(viper.GetString("workspace"))
		if err != nil {
			fmt.Print(fmt.Errorf("error %v", err))
			return
		}
		var request *model.RequestMold
		for _, m := range requests {
			if m.Name() == runArgs.Request {
				request = m
				break
			}
		}
		if request == nil {
			fmt.Printf("Could not find a request with name '%s' under workspace '%s'", runArgs.Request, viper.GetString("workspace"))
			return
		}

		profiles, err := loader.ReadProfiles(viper.GetString("workspace"))
		if err != nil {
			fmt.Print(fmt.Errorf("error %v", err))
			return
		}
		profileName := runArgs.Profile
		if len(profileName) == 0 {
			profileName = "default"
		}
		var profile *model.Profile
		for _, p := range profiles {
			if p.Name == profileName {
				profile = p
				break
			}
		}

		runRequests := requestchain.ResolveRequestChain(request, requests)
		responses, err := runner.RunRequestChain(runRequests, profile, func(took time.Duration, statusCode int) {
			log.Info().Msgf("Request responded with status %d and took %s", statusCode, took)
		})
		if err != nil {
			fmt.Print(fmt.Errorf("error %v", err))
			return
		}

		for _, response := range responses {
			var responseStr string
			if runConfig.Plain {
				responseStr, err = print.SprintPlainResponse(response, runConfig.PrintHeaders, runConfig.PrintBody)
			} else {
				responseStr, err = print.SprintPrettyResponse(response, runConfig.PrintHeaders, runConfig.PrintBody)
			}
			if err != nil {
				fmt.Print(fmt.Errorf("error %v", err))
				return
			}
			fmt.Println(responseStr)
		}

	},
}

func ParseArgs(args []string) RunArgs {
	if len(args) == 0 {
		return RunArgs{}
	} else if len(args) == 1 {
		return RunArgs{args[0], ""}
	}
	return RunArgs{args[0], args[1]}
}

func init() {
	rootCmd.AddCommand(runCmd)

	const printHeadersP = "h"
	const printBodyP = "b"

	runConfig.PrintBody = true

	runCmd.PersistentFlags().BoolVarP(&runConfig.Plain, "plain", "p", false, "Print plain response without styling")
	runCmd.PersistentFlags().Bool("no-body", false, "Print no body")
	runCmd.PersistentFlags().StringSlice("print", []string{}, fmt.Sprintf("Print WHAT\n- '%s'\tPrint response headers\n- '%s'\tPrint response body", printHeadersP, printBodyP))
	// TODO some problem with rootCmd's zerolog setting -- needed for print flags to work properly
	/*
		runCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			if cmd == runCmd {
				printFlags, _ := cmd.Flags().GetStringSlice("print")
				for _, flag := range printFlags {
					if flag == printHeadersP {
						runConfig.PrintHeaders = true
					} else if flag == printBodyP {
						runConfig.PrintBody = true
					}
				}
				noBody, _ := cmd.PersistentFlags().GetBool("no-body")
				runConfig.PrintBody = !noBody
			}
		}
	*/
}
