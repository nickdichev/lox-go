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

func (env *Environment) GetAt(distance int, key string) (Valuer, bool) {
	return env.ancestor(distance).Get(key)
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

func (env *Environment) AssignAt(distance int, key string, v Valuer) bool {
	return env.ancestor(distance).Assign(key, v)
}

func (env *Environment) ancestor(distance int) *Environment {
	cur := env
	for i := 0; i < distance; i++ {
		cur = cur.Enclosing
	}
	return cur
}

func NewEnv() *Environment {
	return &Environment{
		Values: make(map[string]Valuer),
	}
}

func NewEnclosing(env *Environment) *Environment {
	closing := NewEnv()
	closing.Enclosing = env
	return closing
}
