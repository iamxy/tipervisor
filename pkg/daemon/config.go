package daemon

import (
	"os/user"
)

// Config maintains the configurations for daemon to run process
type Config struct {
	Name      string
	Cmd       string
	Args      []string
	Cwd       string
	Env       map[string]string
	User      string
	StatusDir string

	pidfile string
	user    *user.User
}
