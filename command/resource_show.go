package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/minamijoyo/tfschema/tfschema"
	"github.com/posener/complete"
)

type ResourceShowCommand struct {
	Meta
	format string
}

func (c *ResourceShowCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("resource show", flag.ContinueOnError)
	cmdFlags.StringVar(&c.format, "format", "table", "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if len(cmdFlags.Args()) != 1 {
		c.Ui.Error("The resource show command expects RESOURCE_TYPE")
		c.Ui.Error(c.Help())
		return 1
	}

	resourceType := cmdFlags.Args()[0]
	providerName, err := detectProviderName(resourceType)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	client, err := tfschema.NewClient(providerName)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	defer client.Kill()

	block, err := client.GetResourceTypeSchema(resourceType)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	var out string
	switch c.format {
	case "table":
		out, err = block.FormatTable()
	case "json":
		out, err = block.FormatJSON()
	default:
		c.Ui.Error(fmt.Sprintf("Unknown output format: %s", c.format))
		c.Ui.Error(c.Help())
		return 1
	}

	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(out)

	return 0
}

func (c *ResourceShowCommand) AutocompleteArgs() complete.Predictor {
	return c.completePredictResourceType()
}

func (c *ResourceShowCommand) AutocompleteFlags() complete.Flags {
	return nil
}

func (c *ResourceShowCommand) Help() string {
	helpText := `
Usage: tfschema resource show [options] RESOURCE_TYPE

Options:

  -format=type    Set output format to table or json (default: table)
`
	return strings.TrimSpace(helpText)
}

func (c *ResourceShowCommand) Synopsis() string {
	return "Show a type definition of resource"
}
