package cmd

import (
	"errors"
	"offergen/logging"
)

type Cmd interface {
	Execute()
}

type RootCmd struct {
	cmds map[string]Cmd
}

func NewRootCmd() *RootCmd {
	cmds := make(map[string]Cmd)
	cmds["serve"] = &ServeCmd{}
	cmds["migrate"] = &MigrateCmd{}

	return &RootCmd{
		cmds: cmds,
	}
}

func (rc *RootCmd) Execute(name string) {
	for cmdName, cmd := range rc.cmds {
		if name == cmdName {
			logger.Info("running cmd", "cmdName", name)
			cmd.Execute()
			return
		}
	}
	logger.Error("cmd not found", "cmdName", name)
	panic(errors.New("cmd not found"))
}

var logger = logging.GetLogger()
