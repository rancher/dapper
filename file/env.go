package file

import (
	"strings"
)

type Context map[string]string

func (c Context) Source() string {
	source := "/source/"
	if v, ok := c["DAPPER_SOURCE"]; ok && v != "" {
		source = v
	}

	if !strings.HasSuffix(source, "/") {
		source += "/"
	}

	return source
}

func (c Context) Cp() string {
	if v, ok := c["DAPPER_CP"]; ok && v != "" {
		return v
	}
	return "."
}

func (c Context) Socket() bool {
	if v, ok := c["DAPPER_DOCKER_SOCKET"]; ok && v != "" {
		return "true" == v
	}
	return false
}

func (c Context) Mode(mode string) string {
	switch mode {
	case "cp", "bind":
		return mode
	}
	return "cp"
}

func (c Context) Env() []string {
	val := []string{}
	if v, ok := c["DAPPER_ENV"]; ok && v != "" {
		val = strings.Split(v, " ")
	}

	ret := []string{}

	for _, i := range val {
		i = strings.TrimSpace(i)
		if i != "" {
			ret = append(ret, i)
		}
	}

	return ret
}

func (c Context) Shell() string {
	if shell, ok := c["SHELL"]; ok && shell != "" {
		return shell
	}
	return "/bin/bash"
}

func (c Context) Output() []string {
	if v, ok := c["DAPPER_OUTPUT"]; ok {
		ret := []string{}
		for _, i := range strings.Split(v, " ") {
			i = strings.TrimSpace(i)
			if i != "" {
				ret = append(ret, i)
			}
		}
		return ret
	}
	return []string{}
}

func (c Context) RunArgs() []string {
	if v, ok := c["DAPPER_RUN_ARGS"]; ok {
		ret := []string{}
		for _, i := range strings.Split(v, " ") {
			i = strings.TrimSpace(i)
			if i != "" {
				ret = append(ret, i)
			}
		}
		return ret
	}
	return []string{}
}
