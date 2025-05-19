package cmds

import (
	"cryphtron/internal/utils"
	"cryphtron/internal/utils/log"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
)

const WgetGhCmdID = "ghgt"

var WgetGhCmdName = CmdName(WgetGhCmdID)

func init() {
	AppendCmd(WgetGhCmdName, WgetGhCmd)
}

type ghGetFlagSet struct {
	*flag.FlagSet
	Owner        string
	Project      string
	Domain       string
	Tag          string
	ArticaftName string
	Output       string
}

func ghGetFlags() *ghGetFlagSet {
	result := &ghGetFlagSet{FlagSet: flag.NewFlagSet(WgetGhCmdName, flag.ContinueOnError)}
	result.StringVar(&result.Domain, "domain", "github.com", "the domain of github articafts")
	result.StringVar(&result.Project, "project", "", "the project of the github articaft")
	result.StringVar(&result.Owner, "owner", "", "the owner of the github articaft")
	result.StringVar(&result.Tag, "tag", "latest", "the tag of the github articaft")
	result.StringVar(&result.ArticaftName, "name", "", "the name of the github articaft")
	result.StringVar(&result.Output, "out", "", "the path of output")
	return result
}

func (f *ghGetFlagSet) normalizeFlags(args []string) ([]string, error) {
	err := f.Parse(args)
	if err != nil {
		return nil, err
	}
	if len(f.Owner) == 0 || len(f.Project) == 0 || len(f.ArticaftName) == 0 {
		return nil, errors.New("owner, project and articaft name must be specified")
	}
	if len(strings.TrimSpace(f.Domain)) == 0 {
		f.Domain = "github.com"
	}
	if len(strings.TrimSpace(f.Tag)) == 0 {
		f.Tag = "latest"
	}
	if len(f.Output) == 0 {
		wd, _ := os.Getwd()
		f.Output = wd
	}
	return f.Args(), nil
}

func WgetGhCmd(args []string) error {
	flags := ghGetFlags()
	args, err := flags.normalizeFlags(args)
	if err != nil {
		return err
	}
	if len(args) != 0 {
		return errors.New("too many arguments")
	}
	apiUrl := "api." + flags.Domain + "/repos/" + flags.Owner + "/" + flags.Project + "/releases"
	if flags.Tag == "latest" {
		apiUrl = apiUrl + "/latest"
	} else {
		apiUrl = apiUrl + "/tags/" + flags.Tag
	}
	utils.LogCLI(log.InfoLevel).Println(WgetGhCmdName, ":", "Getting latest version of", flags.ArticaftName, "from", apiUrl)
	result, err := utils.HttpRequestBytes("https://"+apiUrl, "GET")
	if err != nil {
		return err
	}
	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name     string `json:"name"`
			Download string `json:"browser_download_url"`
		} `json:"assets"`
	}
	err = json.Unmarshal(result, &release)
	if err != nil {
		return err
	}
	if len(release.TagName) == 0 {
		return errors.New("no release found")
	}
	utils.LogCLI(log.InfoLevel).Println(WgetGhCmdName, ":", "Downloading", flags.ArticaftName, "from", apiUrl, "with tag", release.TagName)
	url := ""
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, flags.ArticaftName) {
			url = asset.Download
			break
		}
	}
	if url == "" {
		return errors.New("no asset found")
	}
	return utils.HttpRequestFile(url, flags.Output)
}
