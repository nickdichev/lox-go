package valuer

type Environment struct {
	Values    map[string]Valuer
	Enclosing *Environment
}

func (env *Environment) Define(key string, v Valuer) {
	env.Values[key] = v
}

func (env *Environment) Get(key string) (Valuer, bool) {
	v, ok := env.Values[key]
	return v, ok
}

func (env *Environment) Assign(key string, v Valuer) bool {
	if _, ok := env.Values[key]; ok {
		env.Values[key] = v
		return true
	}
	return false
}

func New() *Environment {
	return &Environment{
		Values: make(map[string]Valuer),
	}
}
