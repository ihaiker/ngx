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
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use: filepath.Base(os.Args[0]), Args: cobra.RangeArgs(0, 20),
		Example: `    cat nginx.conf | ngx '.http.server.server_na
    ngx -f nginx.conf '.http.server.server_name'
	more information and document at https://ihaiker.github.io/ngx .`,
	}
	input_p    = rootCmd.PersistentFlags().StringP("input", "i", "stdin", "the parse filename")
	output_p   = rootCmd.PersistentFlags().StringP("output", "o", "stdout", "print the result to file")
	hook       = rootCmd.PersistentFlags().BoolP("hook", "H", false, "execute hooks")
	color      = rootCmd.Flags().BoolP("color", "C", true, "colorize output, ineffective if set -j or -J")
	normalJson = rootCmd.Flags().BoolP("json", "j", false, "print result simple JSON")
	simpleJson = rootCmd.Flags().BoolP("JSON", "J", false, "print result JSON")
	root       = rootCmd.PersistentFlags().String("root", "${PWD}", "the root path for include hook .")
	merge      = rootCmd.PersistentFlags().Bool("merge", false, "enable the merge mode of include")
	ignore     = rootCmd.PersistentFlags().Bool("ignore-errors", false, "Ignore error when not found")
	args       = rootCmd.PersistentFlags().StringToString("args", nil, "Specify parameters, It takes effect when the --hook, -H parameter is set")
	value      = rootCmd.Flags().BoolP("value", "v", false, "only show args value, When it is effective that only one selected directive found")

	setArgsCmd = &cobra.Command{
		Use: "set", Short: "Set config directive args.",
		Args:    cobra.MinimumNArgs(1),
		Example: "  ngx set <pql> arg0 arg1",
	}

	addBodyCmd = &cobra.Command{
		Use: "add", Short: "Append directive to select directive", Args: cobra.MinimumNArgs(2),
		Example: "ngx append <pql> bodyDirective0 bodyDirective1",
	}
	first = rootCmd.PersistentFlags().Bool("first", false, "add directive to the first")
	index = rootCmd.PersistentFlags().Int("index", -1, "add directive location index, -1 is last")

	delCmd = &cobra.Command{
		Use: "del", Short: "Delete config directive.", Args: cobra.MinimumNArgs(1),
		Example: "ngx del <delete_pql1> <delete_pql2>",
	}

	completionCmd = &cobra.Command{
		Use: "completion", Args: cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Short:     "Generates completion scripts",
		Example:   "ngx completion bash",
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

func init() {
	rootCmd.RunE = rune
	setArgsCmd.RunE = setRune
	addBodyCmd.RunE = addRune
	delCmd.RunE = deleteRune
	rootCmd.AddCommand(setArgsCmd, addBodyCmd, delCmd, completionCmd)
	rootCmd.SilenceUsage = true
}

func stdin() (conf *config.Configuration, err error) {
	var stat os.FileInfo
	if stat, err = os.Stdin.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		if conf, err = config.ParseIO(os.Stdin); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("inviold stdin")
	}
	return
}

func input() (conf *config.Configuration, err error) {
	if *input_p == "stdin" {
		conf, err = stdin()
	} else {
		conf, err = config.Parse(*input_p)
	}
	return
}

func output(datas []byte) (err error) {
	stream := os.Stdout
	if *output_p != "stdout" {
		if stream, err = os.OpenFile(*output_p, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			return
		}
	}
	if _, err = stream.Write(datas); err != nil {
		return
	}
	if *output_p == "stdout" {
		if len(datas) > 0 && datas[len(datas)-1] != '\n' {
			_, _ = stream.Write([]byte("\n"))
		}
	}
	return
}

func marshal(conf *config.Configuration) (datas []byte, err error) {
	*color = (*output_p == "stdout") && *color

	//当只有一个结果的时候，用户指定了--value直接返回内容
	if *value && len(conf.Body) == 1 && len(conf.Body[0].Body) == 0 {
		if len(conf.Body[0].Args) == 0 {
			datas = []byte{}
		} else if len(conf.Body[0].Args) == 1 {
			datas = []byte(conf.Body[0].Args[0])
		} else {
			datas, _ = json.Marshal(conf.Body[0].Args)
		}
		return
	}

	if *simpleJson {
		datas, err = encoding.JsonIndent(conf, "    ", "    ")
	} else if *normalJson {
		datas, err = json.MarshalIndent(conf, "    ", "    ")
	} else {
		datas = []byte(conf.Pretty(*color))
	}
	return
}

func executeHook(conf *config.Configuration) error {
	if *hook {
		path := os.ExpandEnv(*root)
		hooks.Defaults.Hook(&hooks.IncludeHooker{Merge: *merge, Search: hooks.WalkFiles(path)}, "include")
		for name, value := range *args {
			hooks.Defaults.Variables.Parameter(name, value)
		}
		if err := hooks.Defaults.Execute(conf); err != nil {
			return err
		}
	}
	return nil
}

func rune(cmd *cobra.Command, args []string) (err error) {
	var conf *config.Configuration
	if conf, err = input(); err != nil {
		return
	}

	if err = executeHook(conf); err != nil {
		return
	}

	if len(args) > 0 && args[0] != "." {
		if conf.Body, err = query.Selects(conf, strings.Join(args, " | ")); err != nil {
			if !query.IsNotFound(err) || !*ignore {
				return
			}
		}
		conf.Source = conf.Source + " and query: " + strings.Join(args, " | ")
	}

	var datas []byte
	if datas, err = marshal(conf); err != nil {
		return
	}
	return output(datas)
}

func setRune(cmd *cobra.Command, args []string) (err error) {
	filter := args[0]
	args = args[1:]

	var conf *config.Configuration
	if conf, err = input(); err != nil {
		return
	}

	if err = executeHook(conf); err != nil {
		return
	}

	//设置内容
	var items config.Directives
	if items, err = query.Selects(conf, filter); err != nil {
		//ignore not found
		if !query.IsNotFound(err) || !*ignore {
			return
		}
	} else {
		for _, item := range items {
			item.Args = args
		}
	}

	var datas []byte
	if datas, err = marshal(conf); err != nil {
		return
	}
	return output(datas)
}

func addRune(cmd *cobra.Command, args []string) (err error) {
	filter := args[0]
	args = args[1:]

	var conf *config.Configuration
	if conf, err = input(); err != nil {
		return
	}

	if err = executeHook(conf); err != nil {
		return
	}

	//设置内容
	var items config.Directives
	if items, err = query.Selects(conf, filter); err != nil {
		//ignore not found
		if err != query.ErrNotFound || *ignore {
			return
		}
	} else {
		newItems := config.Directives{}
		for _, arg := range args {
			if appConf, err := config.ParseBytes([]byte(arg)); err == nil {
				newItems = append(newItems, appConf.Body...)
			} else {
				return err
			}
		}

		for _, item := range items {
			if *first {
				item.Body = append(newItems, item.Body...)
			} else if *index != -1 {
				if !(*index < len(item.Body)) {
					return fmt.Errorf("out of index, %d > %d", *index, len(item.Body))
				}
				item.Body = append(item.Body[:*index], append(newItems, item.Body[*index:]...)...)
			} else {
				item.Body = append(item.Body, newItems...)
			}
		}
	}

	var datas []byte
	if datas, err = marshal(conf); err != nil {
		return
	}
	return output(datas)
}

func deleteRune(cmd *cobra.Command, args []string) (err error) {
	var conf *config.Configuration
	if conf, err = input(); err != nil {
		return
	}

	if err = executeHook(conf); err != nil {
		return
	}

	var items config.Directives
	for _, arg := range args {
		if items, err = query.Selects(conf, arg); err != nil {
			if query.IsNotFound(err) && !*ignore {
				continue
			}
			return fmt.Errorf("select %s: %s", arg, err.Error())
		}
		for _, item := range items {
			item.Name = ""
		}
	}

	var datas []byte
	if datas, err = marshal(conf); err != nil {
		return
	}
	return output(datas)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
