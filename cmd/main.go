package main

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/encoding"
	"github.com/ihaiker/ngx/v2/hooks"
	"github.com/ihaiker/ngx/v2/query"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	rootCmd = &cobra.Command{
		Use: filepath.Base(os.Args[0]), Args: cobra.MaximumNArgs(1),
		Example: `    cat nginx.conf | ngx '.http.server.server_name'
    ngx -f nginx.conf '.http.server.server_name'
more info and document at https://ihaiker.github.io/ngx .`,
	}
	hook       = rootCmd.PersistentFlags().BoolP("hook", "H", false, "execute hooks")
	input      = rootCmd.PersistentFlags().StringP("file", "f", "stdin", "the parse file name")
	color      = rootCmd.PersistentFlags().BoolP("color", "C", true, "colorize output, ineffective if set -j or -J")
	output     = rootCmd.PersistentFlags().StringP("output", "o", "stdout", "print the result to file")
	normalJson = rootCmd.PersistentFlags().BoolP("json", "j", false, "print result simple JSON")
	simpleJson = rootCmd.PersistentFlags().BoolP("JSON", "J", false, "print result JSON")
	pwd, _     = os.Getwd()
	root       = rootCmd.PersistentFlags().String("root", pwd, "the include directive root path.")

	completionCmd = &cobra.Command{
		Use: "completion", Args: cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Short:     "Generates completion scripts",
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
)

func rune(cmd *cobra.Command, args []string) (err error) {
	var conf *config.Configuration
	if *input == "stdin" {
		if stat, err := os.Stdin.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			if conf, err = config.ParseIO(os.Stdin); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("inviold stdin")
		}
	} else {
		if conf, err = config.Parse(*input); err != nil {
			return
		}
	}

	if *hook {
		hooks.Defaults.Hook(&hooks.IncludeHooker{Merge: false, Search: hooks.WalkFiles(*root)}, "include")
		if err = hooks.Defaults.Execute(conf); err != nil {
			return
		}
	}

	if len(args) > 0 {
		if conf.Body, err = query.Selects(conf, args[0]); err != nil {
			return
		}
		conf.Source = "query: " + args[0]
	}

	var datas []byte
	if *simpleJson {
		datas, err = encoding.JsonIndent(conf, "    ", "    ")
	} else if *normalJson {
		datas, err = json.MarshalIndent(conf, "    ", "    ")
	} else {
		datas = []byte(conf.Pretty(*color))
	}
	_, err = os.Stdout.Write(datas)
	if len(datas) > 0 && datas[len(datas)-1] != '\n' {
		_, _ = os.Stdout.WriteString("\n")
	}
	return
}

func main() {
	rootCmd.AddCommand(completionCmd)
	rootCmd.RunE = rune
	rootCmd.SilenceUsage = true
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
