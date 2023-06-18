package util

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
)

const bashTemplateDef = `_<<.Context.App.Name>>_bash_autocomplete() {
	local cur opts;
	COMPREPLY=();
	cur=$(_get_cword)
	opts='<<range $idx, $flag := .Context.App.VisibleFlags>><<range .Names>><<opthyphens .>> <<end>><<end>>'
	COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) );
	return 0;
};
complete -F _<<.Context.App.Name>>_bash_autocomplete <<.Program>>
`

const zshTemplateDef = `#compdef _<<.Program>> <<.Program>>

function _<<.Program>> {
	local line

	_arguments -C \
		<<range $idx, $flag := .Context.App.VisibleFlags>><<range .Names>>"<<opthyphens .>>[<<if $flag.Usage>><<$flag.Usage>><<end>>]" \
		<<end>><<end>>"1: :(<<range .Context.App.Commands>><<.Name>> <<end>>)" \
		"'*:URL:_urls'"
}
`

// InitCompletionFlag generates completion code
var InitCompletionFlag = &cli.StringFlag{
	Name:  "create-completion",
	Usage: "generate completion code. Value must be 'bash' or 'zsh'",
}

func createTemplateContext(ctx *cli.Context) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["Program"] = os.Args[0]
	ret["CompleteFlag"] = cli.BashCompletionFlag.Names()[0]
	ret["Context"] = ctx

	return ret
}

func createTemplateFunctions() template.FuncMap {
	ret := make(template.FuncMap)
	ret["join"] = strings.Join
	ret["opthyphens"] = func(optName string) string {
		if len(optName) == 1 {
			return "-" + optName
		} else {
			return "--" + optName
		}
	}

	return ret
}

// PrintCompletion prints a generated bash or zsh completion script.
func PrintCompletion(ctx *cli.Context) {
	shell := ctx.String(InitCompletionFlag.Name)

	bashTemplate := template.Must(template.New("bash").
		Delims("<<", ">>").
		Funcs(createTemplateFunctions()).
		Parse(bashTemplateDef))

	zshTemplate := template.Must(template.New("zsh").
		Delims("<<", ">>").
		Funcs(createTemplateFunctions()).
		Parse(zshTemplateDef))

	switch shell {
	case "bash":
		if err := bashTemplate.Execute(os.Stdout, createTemplateContext(ctx)); err != nil {
			panic(err)
		}
	case "zsh":
		if err := zshTemplate.Execute(os.Stdout, createTemplateContext(ctx)); err != nil {
			panic(err)
		}
	default:
		cli.ShowAppHelp(ctx)
		fmt.Printf("no autocomplete support for %s", shell)
	}
}
