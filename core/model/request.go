package model

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	starlarkNamePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?mU)^.*meta:name:(.*)$`),
	}
	starlarkUrlPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?mU)^.*doc:url:(.*)$`),
		regexp.MustCompile(`(?mU)^\s*url\s*=(.*)$`),
	}
	starlarkMethodPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?mU)^.*doc:method:(.*)$`),
		regexp.MustCompile(`(?mU)^\s*method\s*=(.*)$`),
	}
	starlarkPrevReqPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?mU)^.*meta:prev_req:(.*)$`),
	}
	starlarkOutputPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?mU)^.*meta:output:(.*)$`),
	}
)

type Request struct {
	Url     string
	Method  string
	Headers Headers
	Body    Body
	Output  string
	Options map[string]interface{}
}

type RequestMold struct {
	Yaml        *YamlRequest
	Starlark    *StarlarkRequest
	ContentType string
	Root        string
	Filename    string
}

type YamlRequest struct {
	Name    string
	PrevReq string                 `yaml:"prev_req"`
	Url     string                 `yaml:"url"`
	Method  string                 `yaml:"method"`
	Headers Headers                `yaml:"headers"`
	Body    Body                   `yaml:"body"`
	Output  string                 `yaml:"output"`
	Options map[string]interface{} `yaml:"options"`
	Raw     string
}

type StarlarkRequest struct {
	Script string
}

func (r *Request) IsForm() bool {
	contentType, ok := r.Headers["Content-Type"]
	if !ok {
		return false
	}
	if len(contentType) == 0 {
		return false
	}
	return strings.ToLower(contentType[0]) == "application/x-www-form-urlencoded"
}

func (r *Request) BodyAsMap() (map[string]string, bool) {
	asMapInterface, ok := r.Body.(map[string]interface{})
	if !ok {
		return map[string]string{}, false
	}
	asMapString := make(map[string]string)
	for k, v := range asMapInterface {
		asMapString[k] = v.(string)
	}
	return asMapString, true
}

func (r *RequestMold) Name() string {
	if r.Yaml != nil {
		return r.Yaml.Name
	} else if r.Starlark != nil {
		return findWithPatterns(r.Starlark.Script, starlarkNamePatterns)
	}
	return ""
}

func (r *RequestMold) Url() string {
	if r.Yaml != nil {
		return r.Yaml.Url
	} else if r.Starlark != nil {
		return findWithPatterns(r.Starlark.Script, starlarkUrlPatterns)
	}
	return ""
}

func (r *RequestMold) Method() string {
	if r.Yaml != nil {
		return r.Yaml.Method
	} else if r.Starlark != nil {
		return findWithPatterns(r.Starlark.Script, starlarkMethodPatterns)
	}
	return ""
}

func (r *RequestMold) Raw() string {
	if r.Yaml != nil {
		return r.Yaml.Raw
	} else if r.Starlark != nil {
		return r.Starlark.Script
	}
	return ""
}

func (r *RequestMold) PreviousReq() string {
	if r.Yaml != nil {
		return r.Yaml.PrevReq
	} else if r.Starlark != nil {
		return findWithPatterns(r.Starlark.Script, starlarkPrevReqPatterns)
	}
	return ""
}

func (r *RequestMold) Output() string {
	if r.Yaml != nil {
		return r.Yaml.Output
	} else if r.Starlark != nil {
		return findWithPatterns(r.Starlark.Script, starlarkOutputPatterns)
	}
	return ""
}

func findWithPatterns(str string, patterns []*regexp.Regexp) string {
	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(str)
		if len(match) == 2 {
			return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(match[1]), "\"", ""), "'", "")
		}
	}
	return ""
}

func (r *RequestMold) DeleteFromFS() bool {
	err := os.Remove(filepath.Join(r.Root, r.Filename))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to remove file %s", r.Filename)
		return false
	}
	return true
}

func (r *RequestMold) Clone() RequestMold {
	copy := RequestMold{
		ContentType: r.ContentType,
		Root:        r.Root,
		Filename:    r.Filename,
	}

	if r.Yaml != nil {
		yamlRequest := YamlRequest{
			Name:    r.Yaml.Name,
			PrevReq: r.Yaml.PrevReq,
			Url:     r.Yaml.Url,
			Method:  r.Yaml.Method,
			Headers: r.Yaml.Headers,
			Body:    r.Yaml.Body,
			Output:  r.Yaml.Output,
			Raw:     r.Yaml.Raw,
		}
		copy.Yaml = &yamlRequest
	} else if r.Starlark != nil {
		starlarkRequest := StarlarkRequest{
			Script: r.Starlark.Script,
		}
		copy.Starlark = &starlarkRequest
	}

	return copy
}
