package valuer

type Environment struct {
	Values    map[string]Valuer
	Enclosing *Environment
}

func (env *Environment) Define(key string, v Valuer) {
	env.Values[key] = v
}

func (env *Environment) Get(key string) (Valuer, bool) {
	if v, ok := env.Values[key]; ok {
		return v, true
	}
	if env.Enclosing != nil {
		return env.Enclosing.Get(key)
	}
	return nil, false
}

func (env *Environment) Assign(key string, v Valuer) bool {
	if _, ok := env.Values[key]; ok {
		env.Values[key] = v
		return true
	}
	if env.Enclosing != nil {
		return env.Enclosing.Assign(key, v)
	}
	return false
}

func New() *Environment {
	return &Environment{
		Values: make(map[string]Valuer),
	}
}

func NewEnclosing(env *Environment) *Environment {
	closing := New()
	closing.Enclosing = env
	return closing
}
