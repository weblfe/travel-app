package libs

func ArgsInt(args []int, defs int) int {
		if len(args) == 0 {
				return defs
		}
		return args[0]
}

func ArgsStr(args []string,defs string) string  {
		if len(args) == 0 {
				return defs
		}
		return args[0]
}

func ArgsIntN(args []int64,defs int64) int64  {
		if len(args) == 0 {
				return defs
		}
		return args[0]
}

func ArgsAny(args []interface{},defs interface{}) interface{}  {
		if len(args) == 0 {
				return defs
		}
		return args[0]
}

func ArgsBool(args []bool,defs bool) bool  {
		if len(args) == 0 {
				return defs
		}
		return args[0]
}
