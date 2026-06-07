package cmds

import (
	"encoding/json"
	"flag"
	"net/url"
	"strings"
)

const WgetGhCmdID = "ghgt"

var WgetGhCmdName = CmdName(WgetGhCmdID)

func init() {
	AppendCmd(WgetGhCmdName, WgetGhCmd)
}

type ghGetFlagSet struct {
	*flag.FlagSet
	Owner          string
	Project        string
	Domain         string
	Tag            string
	ArticaftName   string
	Output         string
	ApiDomain      string
	IsConcatDomain bool
}

func ghGetFlags() *ghGetFlagSet {
	result := &ghGetFlagSet{FlagSet: flag.NewFlagSet(WgetGhCmdName, flag.ContinueOnError)}
	result.StringVar(&result.Domain, "domain", "github.com", "the domain of github articafts")
	result.StringVar(&result.Project, "project", "", "the project of the github articaft")
	result.StringVar(&result.Owner, "owner", "", "the owner of the github articaft")
	result.StringVar(&result.Tag, "tag", "latest", "the tag of the github articaft")
	result.StringVar(&result.ArticaftName, "name", "", "the name of the github articaft")
	result.StringVar(&result.Output, "out", "", "the path of output")
	result.StringVar(&result.ApiDomain, "api-domain", "api.github.com", "the api domain of github")
	result.BoolVar(&result.IsConcatDomain, "is-concat", false, "is domain need to concat original")
	return result
}

func (f *ghGetFlagSet) normalizeFlags(args []string) ([]string, error) {
	err := f.Parse(args)
	if err != nil {
		return nil, err
	}
	if len(f.Owner) == 0 || len(f.Project) == 0 || len(f.ArticaftName) == 0 {
		return nil, cmdError("owner, project and articaft name must be specified")
	}
	if len(strings.TrimSpace(f.Domain)) == 0 {
		f.Domain = "github.com"
	}
	if len(strings.TrimSpace(f.Tag)) == 0 {
		f.Tag = ""
	}
	if len(f.Output) == 0 {
		f.Output = GetWorkDir()
	}
	if len(strings.TrimSpace(f.ApiDomain)) == 0 {
		f.ApiDomain = "api.github.com"
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
		return cmdError("too many arguments")
	}
	apiUrl := flags.ApiDomain + "/repos/" + flags.Owner + "/" + flags.Project + "/releases"
	if len(flags.Tag) == 0 {
		apiUrl = apiUrl + "/latest"
	} else {
		apiUrl = apiUrl + "/tags/" + flags.Tag
	}
	LogInfo(WgetGhCmdName, ":", "Getting latest version of", flags.ArticaftName, "from", apiUrl)
	result, err := HttpRequestBytesWithAgentSuffix("https://"+apiUrl, AgentName(WgetGhCmdID), "GET")
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
		return cmdError("no release found")
	}
	downloadUrlStr := ""
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, flags.ArticaftName) {
			downloadUrlStr = asset.Download
			break
		}
	}
	LogInfo(WgetGhCmdName, ":", "Downloading", flags.ArticaftName, "from", downloadUrlStr, "with tag", release.TagName)
	if downloadUrlStr == "" {
		return cmdError("no asset found")
	}
	if flags.IsConcatDomain {
		downloadUrlStr = "https://" + strings.TrimSpace(flags.Domain) + "/" + downloadUrlStr
	}
	LogDebug(WgetGhCmdName, ":", "Raw URL from", downloadUrlStr)
	downloadUrl, err := url.Parse(downloadUrlStr)
	if err != nil {
		return err
	}
	if len(flags.Domain) > 0 && !flags.IsConcatDomain {
		downloadUrl.Host = flags.Domain
	}
	LogDebug(WgetGhCmdName, ":", "Final URL is ", downloadUrl.String())
	return HttpRequestFileWithAgentSuffix(downloadUrl.String(), AgentName(WgetCmdID), flags.Output)
}
