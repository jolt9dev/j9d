package deployments

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jolt9dev/j9d/pkg/ctxs"
	"github.com/jolt9dev/j9d/pkg/env"
	"github.com/jolt9dev/j9d/pkg/platform"
	exec "github.com/jolt9dev/j9d/pkg/xexec"
	fs "github.com/jolt9dev/j9d/pkg/xfs"
)

type DeployParams struct {
	Workspace string
	Project   string
	File      string
}

func Deploy(params DeployParams) error {

	ctx, err := ctxs.Load(params.File)
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

func deployCompose(ctx *ctxs.ExecContext) error {

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

		cmd := exec.New("docker", args...)
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
			n, err := fs.Resolve(ctx.Cwd, f)

			if err != nil {
				return err
			}

			args = append(args, "-f", n)
		}

		args = append(args, "up", "-d")

		proc := "docker"
		if j9d.Compose.Sudo && !platform.IsWindows() {
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

	return nil
}

type RemoveParams struct {
	Workspace string
	Project   string
	File      string
}

func Remove(params RemoveParams) error {

	ctx, err := ctxs.Load(params.File)
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
		if j9d.Compose.Sudo && !platform.IsWindows() {
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
			n, err := fs.Resolve(ctx.Cwd, f)

			if err != nil {
				return err
			}

			args = append(args, "-f", n)
		}

		args = append(args, "down")

		proc := "docker"
		if j9d.Compose.Sudo && !platform.IsWindows() {
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

	return nil
}
