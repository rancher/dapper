package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/term"
	"path"
)

type Dapperfile struct {
	File   string
	Mode   string
	docker string
	env    Context
	Socket bool
	NoOut  bool
}

func Lookup(file string) (*Dapperfile, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	d := &Dapperfile{
		File: file,
	}

	return d, d.init()
}

func (d *Dapperfile) init() error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}
	d.docker = docker
	return nil
}

func (d *Dapperfile) Run(commandArgs []string) error {
	tag, err := d.build()
	if err != nil {
		return err
	}

	logrus.Debugf("Running build in %s", tag)
	name, args := d.runArgs(tag, "", commandArgs)
	defer func() {
		logrus.Debugf("Deleting temp container %s", name)
		d.execWithOutput("rm", "-fv", name)
	}()

	if err := d.run(args...); err != nil {
		return err
	}

	source := d.env.Source()
	output := d.env.Output()
	if !d.IsBind() && !d.NoOut {
		for _, i := range output {
			p := i
			if !strings.HasPrefix(p, "/") {
				p = path.Join(source, i)
			}
			targetDir := path.Dir(i)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return err
			}
			logrus.Infof("docker cp %s %s", p, targetDir)
			if err := d.exec("cp", name+":"+p, targetDir); err != nil {
				logrus.Debugf("Error copying back '%s': %s", i, err)
			}
		}
	}

	return nil
}

func (d *Dapperfile) Shell(commandArgs []string) error {
	tag, err := d.build()
	if err != nil {
		return err
	}

	logrus.Debugf("Running shell in %s", tag)
	_, args := d.runArgs(tag, d.env.Shell(), nil)
	args = append([]string{"--rm"}, args...)

	return d.runExec(args...)
}

func (d *Dapperfile) runArgs(tag, shell string, commandArgs []string) (string, []string) {
	name := fmt.Sprintf("%s-%s", strings.Split(tag, ":")[0], randString())

	args := []string{"-i", "--name", name}

	if term.IsTerminal(0) {
		args = append(args, "-t")
	}

	if d.env.Socket() || d.Socket {
		args = append(args, "-v", "/var/run/docker.sock:/var/run/docker.sock")
	}

	if d.IsBind() {
		wd, err := os.Getwd()
		if err == nil {
			args = append(args, "-v", fmt.Sprintf("%s:%s", fmt.Sprintf("%s/%s", wd, d.env.Cp()), d.env.Source()))
		}
	}

	args = append(args, "-e", fmt.Sprintf("DAPPER_UID=%d", os.Getuid()))
	args = append(args, "-e", fmt.Sprintf("DAPPER_GID=%d", os.Getgid()))

	for _, env := range d.env.Env() {
		args = append(args, "-e", env)
	}

	if shell != "" {
		args = append(args, "--entrypoint", shell)
		args = append(args, "-e", "TERM")
	}

	args = append(args, d.env.RunArgs()...)
	args = append(args, tag)

	if shell != "" && len(commandArgs) == 0 {
		args = append(args, "-")
	} else {
		args = append(args, commandArgs...)
	}

	return name, args
}

func buildArgs() []string {
	v := os.Getenv("DAPPER_BUILD_ARGS")
	logrus.Debugf("DAPPER_BUILD_ARGS=%s", v)
	if len(v) > 0 {
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

func (d *Dapperfile) build() (string, error) {
	tag := d.tag()
	logrus.Debugf("Building %s using %s", tag, d.File)
	args := []string{"build", "-t", tag, "-f", d.File}
	args = append(args, buildArgs()...)
	args = append(args, ".")

	logrus.Debugf("Build args %s", args)

	if err := d.exec(args...); err != nil {
		return "", err
	}

	if err := d.readEnv(tag); err != nil {
		return "", err
	}

	if !d.IsBind() {
		text := fmt.Sprintf("FROM %s\nCOPY %s %s", tag, d.env.Cp(), d.env.Source())
		if err := d.buildWithContent(tag, text); err != nil {
			return "", err
		}
	}

	return tag, nil
}

func (d *Dapperfile) buildWithContent(tag, content string) error {
	tempfile, err := ioutil.TempFile(".", d.File)
	if err != nil {
		return err
	}

	logrus.Debugf("Created tempfile %s", tempfile.Name())
	defer func() {
		logrus.Debugf("Deleting tempfile %s", tempfile.Name())
		if err := os.Remove(tempfile.Name()); err != nil {
			logrus.Errorf("Failed to delete tempfile %s: %v", tempfile.Name(), err)
		}
	}()

	ioutil.WriteFile(tempfile.Name(), []byte(content), 0600)

	return d.exec("build", "-t", tag, "-f", tempfile.Name(), ".")
}

func (d *Dapperfile) readEnv(tag string) error {
	var envList []string

	args := []string{"inspect", "-f", "{{json .ContainerConfig.Env}}", tag}

	cmd := exec.Command(d.docker, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("Failed to run docker %v: %v", args, err)
		return err
	}

	if err := json.Unmarshal(output, &envList); err != nil {
		return err
	}

	d.env = map[string]string{}

	for _, item := range envList {
		parts := strings.SplitN(item, "=", 2)
		k, v := parts[0], parts[1]
		logrus.Debugf("Reading Env: %s=%s", k, v)
		d.env[k] = v
	}

	logrus.Debugf("Source: %s", d.env.Source())
	logrus.Debugf("Cp: %s", d.env.Cp())
	logrus.Debugf("Socket: %t", d.env.Socket())
	logrus.Debugf("Mode: %s", d.env.Mode(d.Mode))
	logrus.Debugf("Env: %v", d.env.Env())
	logrus.Debugf("Output: %v", d.env.Output())

	return nil
}

func (d *Dapperfile) tag() string {
	cwd, err := os.Getwd()
	if err == nil {
		cwd = filepath.Base(cwd)
	} else {
		cwd = "dapper-unknown"
	}

	output, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	tag := strings.TrimSpace(string(output))
	if tag == "" {
		tag = randString()
	}

	return fmt.Sprintf("%s:%s", cwd, tag)
}

func (d *Dapperfile) run(args ...string) error {
	return d.exec(append([]string{"run"}, args...)...)
}

func (d *Dapperfile) exec(args ...string) error {
	logrus.Debugf("Running %s %v", d.docker, args)
	cmd := exec.Command(d.docker, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		logrus.Debugf("Failed running %s %v: %v", d.docker, args, err)
	}
	return err
}

func (d *Dapperfile) runExec(args ...string) error {
	logrus.Debugf("Exec %s run %v", d.docker, args)
	return syscall.Exec(d.docker, append([]string{"docker", "run"}, args...), os.Environ())
}

func (d *Dapperfile) execWithOutput(args ...string) ([]byte, error) {
	cmd := exec.Command(d.docker, args...)
	return cmd.CombinedOutput()
}

func (d *Dapperfile) IsBind() bool {
	return d.env.Mode(d.Mode) == "bind"
}
