package glide

import (
	"fmt"
	// "log"
	"os"
	"strings"

	glidecfg "github.com/Masterminds/glide/cfg"
	json2csv "github.com/sniperkit/xutil/plugin/format/convert/json2csv"
	tablib "github.com/sniperkit/xutil/plugin/format/convert/tabular"
	jsoniter "github.com/sniperkit/xutil/plugin/format/json"
)

var (
	json                                    = jsoniter.ConfigCompatibleWithStandardLibrary
	writers  map[string]*json2csv.CSVWriter = make(map[string]*json2csv.CSVWriter, 0)
	sheets   map[string][]interface{}       = make(map[string][]interface{}, 0)
	datasets map[string]*tablib.Dataset     = make(map[string]*tablib.Dataset, 0) // := NewDataset([]string{"firstName", "lastName"})
)

type FileRequest struct {
	Name    string `json:"name" yaml:"name" toml:"name" csv:"name"`
	Content string `json:"content" yaml:"content" toml:"content" csv:"content"`
}

type Request struct {
	Files []FileRequest `json:"files" yaml:"files" toml:"files" csv:"files"`
}

type Specifier struct {
	Operator string `json:"operator,omitempty" yaml:"operator,omitempty" toml:"operator,omitempty" csv:"operator,omitempty"`
	Version  string `json:"version" yaml:"version" toml:"version" csv:"version"`
}

type Dependency struct {
	Name       string                 `json:"name" yaml:"name" toml:"name" csv:"name"`
	Specifiers []Specifier            `json:"specifiers" yaml:"specifiers" toml:"specifiers" csv:"specifiers"`
	Group      string                 `json:"group,omitempty" yaml:"group,omitempty" toml:"group,omitempty" csv:"group,omitempty"`
	Extras     map[string]interface{} `json:"extras" yaml:"extras" toml:"extras" csv:"extras"`
}

type FileResponse struct {
	Name         string       `json:"name" yaml:"name" toml:"name" csv:"name"`
	Dependencies []Dependency `json:"dependencies" dependencies:"dependencies" toml:"dependencies" csv:"dependencies"`
	Error        string       `json:"error,omitempty" yaml:"error,omitempty" toml:"error,omitempty" csv:"error,omitempty"`
}

type Response struct {
	Files []FileResponse `json:"files"`
}

func ParseBytes(input []byte]) (map[string]interface{}, error) {
	var req Request
	if err := json.NewDecoder(input).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "can't parse input JSON: %s\n", err)
		os.Exit(1)
	}
	resp := make(map[string]interface{}, 0)
	return parse(req), nil
}

func ParseString(input string) (map[string]interface{}, error) {
	var req Request
	if err := json.NewDecoder(strings.NewReader(input)).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "can't parse input JSON: %s\n", err)
		os.Exit(1)
	}
	resp := make(map[string]interface{}, 0)
	return parse(req), nil
}

func ParseReader(input io.Reader) (map[string]interface{}, error) {
	var req Request
	if err := json.NewDecoder(input).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "can't parse input JSON: %s\n", err)
		return nil, error
	}
	return parse(req), nil
}

func parse(req *Request) map[string]interface{} {
	var resp Response
	for _, fileReq := range req.Files {
		fileResp := FileResponse{
			Name: fileReq.Name,
		}

		cfg, err := glidecfg.ConfigFromYaml([]byte(fileReq.Content))
		if err == nil {
			for _, src := range []struct {
				group string
				deps  glidecfg.Dependencies
			}{
				{"", cfg.Imports},
				{"development", cfg.DevImports},
			} {
				for _, dep := range src.deps {
					depResp := Dependency{
						Name:   dep.Name,
						Extras: make(map[string]interface{}),
						Group:  src.group,
					}

					if dep.Reference != "" {
						depResp.Specifiers = []Specifier{
							{"", dep.Reference},
						}
						depResp.Extras["reference"] = dep.Reference
					}

					if dep.Repository != "" {
						depResp.Extras["repo"] = dep.Repository
					}

					if dep.VcsType != "" {
						depResp.Extras["vcs"] = dep.VcsType
					}

					if len(dep.Subpackages) > 0 {
						depResp.Extras["subpackages"] = dep.Subpackages
					}

					if len(dep.Os) > 0 {
						depResp.Extras["os"] = dep.Os
					}

					if len(dep.Arch) > 0 {
						depResp.Extras["arch"] = dep.Arch
					}

					fileResp.Dependencies = append(fileResp.Dependencies, depResp)
				}
			}
		} else {
			fileResp.Error = fmt.Sprintf("can't parse input YAML: %s", err)
		}

		resp.Files = append(resp.Files, fileResp)
	}

	return json.NewEncoder(os.Stdout).Encode(resp), nil
}
