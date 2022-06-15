package main

type Statement interface {
	Interpret() (string, error)
}

var (
	_ = Statement(&CreateStatement{})
	_ = Statement(&SelectStatement{})
	_ = Statement(&InsertStatement{})
	_ = Statement(&DeleteStatement{})
)

type CreateStatement struct {
}

func (cs *CreateStatement) Interpret() (string, error) {
	return "", nil
}

type SelectStatement struct {
}

func (ss *SelectStatement) Interpret() (string, error) {
	return "", nil
}

type InsertStatement struct {
}

func (is *InsertStatement) Interpret() (string, error) {
	return "", nil
}

type DeleteStatement struct {
}

func (ds *DeleteStatement) Interpret() (string, error) {
	return "", nil
}

func parse(tokens []*Token) (Statement, error) {
	return nil, nil
}
