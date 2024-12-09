package deployments

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jolt9dev/j9d/pkg/cps"
	"github.com/jolt9dev/j9d/pkg/ctxs"
	"github.com/jolt9dev/j9d/pkg/env"
	"github.com/jolt9dev/j9d/pkg/platform"
	exec "github.com/jolt9dev/j9d/pkg/xexec"
	fs "github.com/jolt9dev/j9d/pkg/xfs"
	"github.com/jolt9dev/j9d/pkg/xstrings"
)

type CommonDeploymentParams struct {
	Workspace string
	Project   string
	File      string
	Target    string
}

type DeployParams struct {
	CommonDeploymentParams
}

func Deploy(params DeployParams) error {

	file, err := getFile(params.CommonDeploymentParams)
	if err != nil {
		return err
	}

	ctx, err := ctxs.Load(file)
	if err != nil {
		return err
	}

	j9d := ctx.Jolt9

	if j9d.Compose != nil {
		err = deployCompose(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFile(params CommonDeploymentParams) (string, error) {
	cwd, err := cps.Cwd()
	if err != nil {
		return "", err
	}

	if params.Project == "" && params.File == "" {
		return filepath.Join(cwd, "j9d.yaml"), nil
	}

	if params.Project != "" {
		if xstrings.Contains(params.Project, "/") {
			parts := strings.Split(params.Project, "/")
			if len(parts) == 2 {
				// special project handling for local projects
				// e.g. @cwd/project -> ./project/j9d.yaml
				// j9d deploy -p @cwd/project -t dev => ./project/dev.j9d.yaml
				if parts[0] == "." || parts[0] == "@cwd" {
					if params.Target != "" {
						f := fmt.Sprintf("%s.jd9.yaml", params.Target)
						return fs.Resolve(filepath.Join(parts[1], f), cwd)
					}

					return fs.Resolve(filepath.Join(parts[1], "j9d.yaml"), cwd)
				} else {
					return "", fmt.Errorf("workspaces are not supported yet")
				}
			}
		} else {
			return "", fmt.Errorf("default workspaces are not supported yet")
		}

	}

	if params.File != "" {
		if params.File == "." {
			return fs.Resolve("j9d.yaml", cwd)
		}

		// handle the case where you're in a root folder of many projects
		// and you simply want to deploy the project using the subfolder name
		// e.g.  j9d deploy -f myproject
		if !filepath.IsAbs(params.File) {
			ext := filepath.Ext(params.File)
			if ext == "" {
				base := filepath.Base(params.File)
				if !strings.HasSuffix(base, "j9d") {
					return fs.Resolve(filepath.Join(params.File, "j9d.yaml"), cwd)
				} else {
					return fs.Resolve(params.File+".yaml", cwd)
				}
			}
		}

		return fs.Resolve(params.File, cwd)
	}

	return fs.Resolve("j9d.yaml", cwd)
}

func deployCompose(ctx *ctxs.ExecContext) error {

	hooks := ctx.Jolt9.Hooks

	if hooks != nil {
		if len(hooks.Before) > 0 {
			err := RunHooks(ctx, hooks.Before)
			if err != nil {
				return err
			}
		}

		if len(hooks.BeforeDeploy) > 0 {
			err := RunHooks(ctx, hooks.BeforeDeploy)
			if err != nil {
				return err
			}
		}
	}

	j9d := ctx.Jolt9
	context := j9d.Compose.Context
	if context == "" {
		context = "default"
	}

	if context != "default" {
		if context != "default" {
			err := ensureContext(context, ctx)
			if err != nil {
				return err
			}
		}
	}

	if j9d.Compose.Mode == "swarm" || j9d.Compose.Mode == "stack" {
		args := []string{"stack", "rm", "--context", context}
		for _, f := range j9d.Compose.Include {
			args = append(args, "-c", f)
		}

		proc := "docker"
		if j9d.Compose.Sudo && !platform.IsWindows() && !cps.IsElevated() {
			args = append([]string{"-E", "docker"}, args...)
			proc = "sudo"
		}

		cmd := exec.New(proc, args...)
		cmd.WithEnvMap(j9d.Env)

		out, err := cmd.Run()

		if err != nil {
			return err
		}

		if out.Code != 0 {
			return fmt.Errorf("docker stack deploy failed: %s", out.ErrorText())
		}
	} else {
		args := []string{"--context", context, "compose", "--project-name", j9d.Name}
		for _, f := range j9d.Compose.Include {
			n, err := fs.Resolve(f, ctx.Cwd)

			if err != nil {
				return err
			}

			args = append(args, "-f", n)
		}

		args = append(args, "up", "-d")

		proc := "docker"
		if j9d.Compose.Sudo && !platform.IsWindows() && !cps.IsElevated() {
			args = append([]string{"-E", "docker"}, args...)
			proc = "sudo"
		}

		cmd := exec.New(proc, args...)
		cmd.WithEnvMap(ctx.Env)

		println(proc, strings.Join(args, " "))
		for k, v := range ctx.Env {
			println(k, v)

			println(k, env.Get(k))
		}

		for _, l := range cmd.Cmd.Env {
			println(l)
		}

		out, err := cmd.Run()

		if err != nil {
			return err
		}

		if out.Code != 0 {
			return fmt.Errorf("docker compose up failed: %s", out.ErrorText())
		}
	}

	if hooks != nil {
		if len(hooks.After) > 0 {
			err := RunHooks(ctx, hooks.After)
			if err != nil {
				return err
			}
		}

		if len(hooks.AfterDeploy) > 0 {
			err := RunHooks(ctx, hooks.AfterDeploy)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type RemoveParams struct {
	CommonDeploymentParams
}

func Remove(params RemoveParams) error {

	file, err := getFile(params.CommonDeploymentParams)
	if err != nil {
		return err
	}

	ctx, err := ctxs.Load(file)
	if err != nil {
		return err
	}

	if ctx.Jolt9.Compose != nil {
		removeCompose(ctx)
	} else {
		return fmt.Errorf("no deployment found e.g. compose block")
	}

	return nil
}

func ensureContext(context string, ctx *ctxs.ExecContext) error {
	j9d := ctx.Jolt9
	out, err := exec.Command("docker context ls --format '{{.Name}}'").Output()
	if err != nil {
		return err
	}

	lines := out.Lines()
	if !slices.Contains(lines, ctx.Jolt9.Name) {
		if j9d.Ssh == nil {
			return fmt.Errorf("context %s not found", ctx.Jolt9.Name)
		}

		host := j9d.Ssh.Host
		port := j9d.Ssh.Port
		user := j9d.Ssh.User

		if port < 1 {
			port = 22
		}

		if strings.Contains(host, "$") {
			host = env.ExpandSafe(host)
		}

		if strings.Contains(user, "$") {
			user = env.ExpandSafe(user)
		}

		format := fmt.Sprintf("ssh://%s@%s", user, host)
		if port != 22 {
			format = fmt.Sprintf("%s:%d", format, port)
		}

		out, err := exec.Command(fmt.Sprintf("docker context create %s --docker %s", ctx.Jolt9.Name, format)).Output()
		if err != nil {
			return err
		}

		if out.Code != 0 {
			return fmt.Errorf("docker context create failed: %s", out.ErrorText())
		}
	}

	return nil
}

func removeCompose(ctx *ctxs.ExecContext) error {

	hooks := ctx.Jolt9.Hooks
	if hooks != nil {
		if len(hooks.Before) > 0 {
			err := RunHooks(ctx, hooks.Before)
			if err != nil {
				return err
			}
		}

		if len(hooks.BeforeRemove) > 0 {
			err := RunHooks(ctx, hooks.BeforeRemove)
			if err != nil {
				return err
			}
		}
	}

	j9d := ctx.Jolt9
	context := j9d.Compose.Context
	if context == "" {
		context = "default"
	}

	if context != "default" {
		err := ensureContext(context, ctx)
		if err != nil {
			return err
		}
	}

	if j9d.Compose.Mode == "swarm" || j9d.Compose.Mode == "stack" {
		args := []string{"stack", "deploy", "--context", context}
		for _, f := range j9d.Compose.Include {
			args = append(args, "-c", f)
		}

		proc := "docker"
		if j9d.Compose.Sudo && !platform.IsWindows() && !cps.IsElevated() {
			args = append([]string{"-E", "docker"}, args...)
			proc = "sudo"
		}

		cmd := exec.New(proc, args...)
		cmd.WithEnvMap(ctx.Env)

		out, err := cmd.Run()

		if err != nil {
			return err
		}

		if out.Code != 0 {
			return fmt.Errorf("docker stack deploy failed: %s", out.ErrorText())
		}
	} else {
		args := []string{"--context", context, "compose", "--project-name", j9d.Name}
		for _, f := range j9d.Compose.Include {
			n, err := fs.Resolve(f, ctx.Cwd)

			if err != nil {
				return err
			}

			args = append(args, "-f", n)
		}

		args = append(args, "down")

		proc := "docker"
		if j9d.Compose.Sudo && !platform.IsWindows() && !cps.IsElevated() {
			args = append([]string{"-E", "docker"}, args...)
			proc = "sudo"
		}

		cmd := exec.New(proc, args...)
		cmd.WithEnvMap(ctx.Env)

		println(proc, strings.Join(args, " "))
		for k, v := range ctx.Env {
			println(k, v)

			println(k, env.Get(k))
		}

		for _, l := range cmd.Cmd.Env {
			println(l)
		}

		out, err := cmd.Run()

		if err != nil {
			return err
		}

		if out.Code != 0 {
			return fmt.Errorf("docker compose down failed: %s", out.ErrorText())
		}
	}

	if hooks != nil {
		if len(hooks.After) > 0 {
			err := RunHooks(ctx, hooks.After)
			if err != nil {
				return err
			}
		}

		if len(hooks.AfterRemove) > 0 {
			err := RunHooks(ctx, hooks.AfterRemove)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
