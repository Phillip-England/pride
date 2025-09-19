package cmd

type CmdHelp struct {
}

func NewCmdHelp() (*CmdHelp, error) {
	cmd := &CmdHelp{}
	return cmd, nil
}

func (cmd CmdHelp) GetOpCode() (int, error) {
	return 0, nil
}
