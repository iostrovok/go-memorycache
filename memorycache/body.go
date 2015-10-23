package memorycache

type Press func(interface{}) (interface{}, error)
type UnPress func(interface{}) (interface{}, error)
