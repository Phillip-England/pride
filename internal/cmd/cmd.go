package cmd

type Cmd interface {
	GetOpCode() (int, error)
}

func Extract() (Cmd, error) {
	id := NewCmdId()
	switch id {
	case FlagNewNew:
		cmd, err := NewCmdNew()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	case FlagEmit:
		cmd, err := NewCmdEmit()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	case FlagServe:
		cmd, err := NewCmdServe()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	case FlagBuild:
		cmd, err := NewCmdBuild()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	default:
		cmd, err := NewCmdHelp()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}
}
