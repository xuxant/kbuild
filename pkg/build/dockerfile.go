package build

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/groupcache/singleflight"
	"github.com/xuxant/kbuild/pkg/options"
	"github.com/xuxant/kbuild/pkg/output/log"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"text/template"
)

type cmdError struct {
	args []string
	stdout []byte
	stderr []byte
	cause error
}

func (e *cmdError) Error() string {
	return fmt.Sprintf("running %s\n - stdout: %q\n - stderr: %q\n - cause: %s", e.args, e.stdout, e.stderr, e.cause)
}

var (
	DefaultExecCommand Command = newCommander()
	dependencyCache = NewSyncStore[[]string]()
)

func newCommander() *Commander {
	return &Commander{
		store: NewSyncStore[[]byte](),
	}
}

func NewSyncStore[T any]() *SyncStore[T] {
	return &SyncStore[T]{
		sf: singleflight.Group{},
		results: syncMap[T]{Map: sync.Map{}},
	}
}

type Commander struct {
	store *SyncStore[[]byte]
}

type Command interface {
	RunCmdOut(cmd *exec.Cmd) ([]byte, error)
}

func RunCmdOut(cmd *exec.Cmd) ([]byte, error) {
	return DefaultExecCommand.RunCmdOut(cmd)
}

func (*Commander) RunCmdOut(cmd *exec.Cmd) ([]byte,error) {
	log.Entry(context.TODO()).Debugf("Running command: %s", cmd.Args)

	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout
	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting command %v: %w", cmd, err)
	}

	if err := cmd.Wait(); err != nil {
		return stdout.Bytes(), &cmdError{
			args: cmd.Args,
			stdout: stdout.Bytes(),
			stderr: stderr.Bytes(),
			cause: err,
		}
	}

	if stderr.Len() > 0 {
		log.Entry(context.TODO()).Debugf("Command Output: [%s], stderr: [%s]", stdout.String(), stderr.String())
	} else {
		log.Entry(context.TODO()).Debugf("Command output: [%s]", stdout.String())
	}

	return stdout.Bytes(), nil
}

var (
	funcsMap = template.FuncMap{
		"cmd": runCmdFunc,
	}
)

func RunCommandOut(cmd *exec.Cmd) ([]byte,error) {
	return DefaultExecCommand.RunCmdOut(cmd)
}

func runCmdFunc(name string, args ...string) (string,error) {
	cmd := exec.Command(name, args...)
	out, err := RunCommandOut(cmd)
}

func createTarContext(w io.Writer, context, dockerfile string) error {
	paths, err :=
}

func getDependenciesCache( context, dockerfile string) ([]string, error) {
	absDockerfilePath, err := normalizeDockerfilePath(context, dockerfile)
	if err != nil {
		return nil, fmt.Errorf("normalizing dockerfile path: %w", err)
	}
	cfg := options.GetConfig()
	return dependencyCache.Exec()

}

func normalizeDockerfilePath(context, dockerfile string) (string, error) {
	rel := filepath.Join(context, dockerfile)
	if _, err := os.Stat(rel); os.IsNotExist(err) {
		if _, err := os.Stat(dockerfile); err == nil || !os.IsNotExist(err) {
			return filepath.Abs(dockerfile)
		}
	}
	if runtime.GOOS == "windows" && (filepath.VolumeName(dockerfile) != "" || filepath.IsAbs(dockerfile)) {
		return dockerfile, nil
	}
	return filepath.Abs(rel)
}


func getDependicies(context, dockerfilePath, absoluteDockerfilePath string) ([]string,error){
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return []string{dockerfilePath}, nil
	}


	return []string{}, nil

}

func ReadCopyCmdsFromDockerfile(onlyLastImage bool, dockerfilePath, context string, buildArgs map[string]*string) ([]string, error) {
	r, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return nil, err
	}

	res, err := parser.Parse(bytes.NewReader(r))
	if err != nil {
		return nil, err
	}

	if err := validateParsedDockerfile(bytes.NewReader(r), res); err != nil {
		return nil, fmt.Errorf("parsing docker file %q: %w", dockerfilePath, err)
	}

	dockerFileLines := res.AST.Children

	if err :=
}

func expandBuildArgs(nodes []*parser.Node, buildArgs map[string]*string) error {
	args, err := utils.Ev
}

func evaluateEnvTemplateMapWithEnv(args map[string]*string) (map[string]*string,error) {
	if args == nil {
		return nil, nil
	}

	evaluated := map[string]*string{}
	for k,v := range args {
		if v == nil {
			evaluated[k] = nil
			continue
		}

		value, err :=
	}
	return evaluated, nil
}

func expandEnvTemplates(s string, envMap map[string]string) (string,error) {
	tmpl, err := template.New("envTemplate").Funcs(funcsMap)
}

func validateParsedDockerfile(r io.Reader, res *parser.Result) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if _,_,_, usesSentax := parser.DetectSyntax(b); usesSentax {
		return nil
	}

	_,_, err = instructions.Parse(res.AST)
	return err
}