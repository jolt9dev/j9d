package deployments

import (
	"fmt"

	"github.com/jolt9dev/j9d/pkg/ctxs"
	"github.com/jolt9dev/j9d/pkg/env"
	"github.com/jolt9dev/j9d/pkg/types"
	"github.com/jolt9dev/j9d/pkg/xexec"
	"github.com/jolt9dev/j9d/pkg/xstrings"
)

func RunHooks(ctx *ctxs.ExecContext, tasks []types.Task) error {

	if len(tasks) == 0 {
		return nil
	}

	for _, hook := range tasks {

		if hook.Use == "" {
			hook.Use = "exec"
		}

		vars := map[string]string{}
		for k, v := range ctx.Env {
			vars[k] = v
		}

		envOptions := &env.ExpandOptions{
			Get: func(key string) string {
				if val, ok := vars[key]; ok {
					return val
				}

				return env.Get(key)
			},

			Set: func(key, value string) error {
				vars[key] = value
				return nil
			},
		}

		if len(hook.Env) > 0 {
			for k, v := range hook.Env {
				if xstrings.Contains(v, "$") {
					v, err := env.Expand(v, envOptions)
					if err != nil {
						return err
					}

					vars[k] = v
					continue
				}

				vars[k] = v
			}
		}

		switch hook.Use {

		case "exec":
			run := hook.Run
			if xstrings.Contains(run, "$") {
				r, err := env.Expand(run, envOptions)
				if err != nil {
					return err
				}

				run = r
			}

			cmd := xexec.Command(run)
			cmd.WithEnvMap(vars)
			cmd.WithCwd(ctx.Cwd)
			_, err := cmd.Run()
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown hook use: %s", hook.Use)
		}
	}

	return nil
}
