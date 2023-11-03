package template

type ImportMap map[string]string

func NewImportMap() ImportMap {
	return make(map[string]string)
}

func (im ImportMap) Add(alias string, importPath string) {
	im[alias] = importPath
}
