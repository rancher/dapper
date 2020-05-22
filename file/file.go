package file

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
)

var (
	re           = regexp.MustCompile("[^a-zA-Z0-9]")
	ErrSkipBuild = errors.New("skip build")
)

type Dapperfile struct {
	File        string
	Mode        string
	docker      string
	env         Context
	Socket      bool
	NoOut       bool
	Args        []string
	From        string
	Quiet       bool
	hostArch    string
	Keep        bool
	NoContext   bool
	MountSuffix string
	Target      string
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
	if d.Args, err = d.argsFromEnv(d.File); err != nil {
		return err
	}
	if d.hostArch == "" {
		d.hostArch = d.findHostArch()
	}
	return nil
}

func (d *Dapperfile) argsFromEnv(dockerfile string) ([]string, error) {
	file, err := os.Open(dockerfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	r := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) <= 1 {
			continue
		}

		command := fields[0]
		if command != "ARG" {
			continue
		}

		key := strings.Split(fields[1], "=")[0]
		value := os.Getenv(key)

		if key == "DAPPER_HOST_ARCH" && value == "" {
			value = d.findHostArch()
		}

		if key == "DAPPER_HOST_ARCH" {
			d.hostArch = value
		}

		if value != "" {
			r = append(r, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return r, nil
}

func (d *Dapperfile) Run(commandArgs []string) error {
	tag, err := d.build(nil, true)
	if err != nil {
		return err
	}

	logrus.Debugf("Running build in %s", tag)
	name, args := d.runArgs(tag, "", commandArgs)
	defer func() {
		if d.Keep {
			logrus.Infof("Keeping build container %s", name)
		} else {
			logrus.Debugf("Deleting temp container %s", name)
			if _, err := d.execWithOutput("rm", "-fv", name); err != nil {
				logrus.Debugf("Error deleting temp container: %s", err)
			}
		}
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
	tag, err := d.build(nil, true)
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

	if isatty.IsTerminal(os.Stdout.Fd()) {
		args = append(args, "-t")
	}

	if d.env.Socket() || d.Socket {
		args = append(args, "-v", d.vSocket())
	}

	if d.IsBind() {
		wd, err := os.Getwd()
		if err == nil {
			suffix := ""
			if d.MountSuffix != "" {
				suffix = ":" + d.MountSuffix
			}
			args = append(args, "-v", fmt.Sprintf("%s:%s%s", fmt.Sprintf("%s/%s", wd, d.env.Cp()), d.env.Source(), suffix))
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

func (d *Dapperfile) findHostArch() string {
	output, err := d.execWithOutput("version", "-f", "{{.Server.Arch}}")
	if err != nil {
		return runtime.GOARCH
	}
	return strings.TrimSpace(string(output))
}

func (d *Dapperfile) Build(args []string) error {
	_, err := d.build(args, false)
	return err
}

func (d *Dapperfile) build(args []string, copy bool) (string, error) {
	dapperFile, err := d.dapperFile()
	if err != nil {
		return "", err
	}

	tag := d.tag()
	logrus.Debugf("Building %s using %s", tag, d.File)
	buildArgs := []string{"build"}
	if len(args) == 0 {
		buildArgs = append(buildArgs, "-t", tag)
	}

	if d.Quiet {
		buildArgs = append(buildArgs, "-q")
	}

	if d.Target != "" {
		buildArgs = append(buildArgs, "--target", d.Target)
	}

	for _, v := range d.Args {
		buildArgs = append(buildArgs, "--build-arg", v)
	}

	if d.NoContext {
		buildArgs = append(buildArgs, "-")
		buildArgs = append(buildArgs, args...)
		if err := d.execWithStdin(bytes.NewBuffer(dapperFile), buildArgs...); err != nil {
			return "", err
		}
	} else {
		tempfile, err := d.tempfile(dapperFile)
		if err != nil {
			return "", err
		}
		defer os.Remove(tempfile)

		buildArgs = append(buildArgs, "-f", tempfile)
		if len(args) > 0 {
			buildArgs = append(buildArgs, args...)
		} else {
			buildArgs = append(buildArgs, ".")
		}

		if err := d.exec(buildArgs...); err != nil {
			return "", err
		}
	}

	if !copy {
		return tag, nil
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
	tempfile, err := d.tempfile([]byte(content))
	if err != nil {
		return err
	}

	defer func() {
		logrus.Debugf("Deleting tempfile %s", tempfile)
		if err := os.Remove(tempfile); err != nil {
			logrus.Errorf("Failed to delete tempfile %s: %v", tempfile, err)
		}
	}()

	return d.exec("build", "-t", tag, "-f", tempfile, ".")
}

func (d *Dapperfile) readEnv(tag string) error {
	var envList []string

	args := []string{"inspect", "-f", "{{json .Config.Env}}", tag}

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
	// repository name must be lowercase
	cwd = strings.ToLower(cwd)

	output, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	tag := strings.TrimSpace(string(output))
	if tag == "" {
		tag = randString()
	}
	tag = re.ReplaceAllLiteralString(tag, "-")

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

func (d *Dapperfile) execWithStdin(stdin io.Reader, args ...string) error {
	logrus.Debugf("Running %s %v", d.docker, args)
	cmd := exec.Command(d.docker, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = stdin
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

func (d *Dapperfile) dapperFile() ([]byte, error) {
	var input io.Reader

	if d.NoContext {
		input = os.Stdin
	} else {
		f, err := os.Open(d.File)
		if err != nil {
			return nil, err
		}
		input = f
		defer f.Close()
	}

	buffer := &bytes.Buffer{}
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "FROM ") && len(strings.Fields(line)) == 2 && scanner.Scan() {
			nextLine := scanner.Text()
			if strings.HasPrefix(nextLine, "# FROM") {
				baseImage, ok := toMap(nextLine)[d.hostArch]
				if ok && baseImage == "skip" {
					return nil, ErrSkipBuild
				}
				if ok {
					line = "FROM " + baseImage
				}
			}
			line = line + "\n" + nextLine
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")
	}

	return buffer.Bytes(), scanner.Err()
}
