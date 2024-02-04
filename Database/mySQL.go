package database

type MySQL struct {
}

func NewMySQL() *MySQL {
	//TODO:SingleTone Pattern
	return &MySQL{}
}

func (obj *MySQL) Insert(tableName string, argsKeys []string, argsVals []string) error {
	//TODO:MYSQL Insert statement
	return nil
}

func (obj *MySQL) start() error {
	//TODO:MYSQL Start statement
	return nil
}
