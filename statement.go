package main

import (
	"reflect"

	"github.com/pkg/errors"
)

type CreateStatement struct {
	table      string
	schema     map[string]reflect.Kind
	primaryKey string
}

type SelectStatement struct {
	table  string
	fields []string
	where  *WhereClause
}

type WhereClause struct {
	field    string
	operator *Operator
	value    any
}

type Operator int

const equal Operator = iota

type InsertStatement struct {
	table  string
	values map[string]any
}

type DeleteStatement struct {
	table string
	keys  []any
	where *WhereClause
}

type DropStatement struct {
	table string
}

func parseSelectStatement(tokens []*Token) (*SelectStatement, error) {
	// skip the first token, i.e., "SELECT"
	i := 1
	if i == len(tokens) {
		return nil, errors.New("incomplete SELECT statement")
	}
	fields := []string{}
	if !cmpTks(*tokens[i], TokenStar) {
		for ; !cmpTks(*tokens[i], TokenFrom); i++ {
			if i%2 == 0 {
				// must be a comma
				if !cmpTks(*tokens[i], TokenComma) {
					return nil, errors.New("the desired field must follow a comma")
				}
				continue
			}
			if !isUnquoteStringToken(tokens[i]) {
				return nil, errors.Errorf("invalid token (%s) before FROM",
					tokens[i].String())
			}
			fields = append(fields, tokens[i].String())
		}
		// check the previous token, in case the FROM follows a comma
		if !isUnquoteStringToken(tokens[i-1]) {
			return nil, errors.Errorf("FROM must follow " +
				"a unquote string token")
		}
	} else {
		// '*' means we will list all fields
		i++
		if i == len(tokens) {
			return nil, errors.New("incomplete SELECT statement")
		}
		if !cmpTks(*tokens[i], TokenFrom) {
			return nil, errors.Errorf("invalid token: got(%s), expect(%s)",
				*tokens[i], TokenFrom)
		}
	}

	i++
	if !isUnquoteStringToken(tokens[i]) {
		return nil, errors.Errorf("FROM must be followed " +
			"by a unquote string token")
	}
	table := tokens[i].String()
	i++
	var where *WhereClause
	if i == len(tokens) {
		goto RETURN
	}

	where = &WhereClause{}
	// parse the where clause if exist
	if !cmpTks(*tokens[i], TokenWhere) {
		// TODO(charleszheng44) more detailed error message
		return nil, errors.New("invalid token")
	}
	i++
	if i == len(tokens) {
		return nil, errors.New("incomplete statement")
	}

	if !isUnquoteStringToken(tokens[i]) {
		return nil, errors.New("invalid token")
	}
	where.field = tokens[i].StringVal
	i++
	if i == len(tokens) {
		return nil, errors.New("incomplete statement")
	}

	if !cmpTks(*tokens[i], TokenEqual) {
		return nil, errors.New("invalid token")
	}
	i++
	if i == len(tokens) {
		return nil, errors.New("incomplete statement")
	}

	switch tokens[i].Type {
	case StringToken:
		where.value = tokens[i].StringVal
	case IntegerToken:
		where.value = tokens[i].IntegerVal
	case FloatToken:
		where.value = tokens[i].FloatVal
	case BoolToken:
		where.value = tokens[i].BoolVal
	default:
		return nil, errors.New("invalid token")
	}

RETURN:
	return &SelectStatement{
		table:  table,
		fields: fields,
		where:  where,
	}, nil
}

func parseField(
	tokens []*Token,
	schema map[string]reflect.Kind,
	i *int) (string, error) {
	if *i >= len(tokens) {
		return "", errors.New("incomplete statement")
	}

	// skip the comma
	if cmpTks(*tokens[*i], TokenComma) {
		*i++
	}

	if tokens[*i].Type != UnquoteStringToken {
		return "", errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[*i].Type, UnquoteStringToken)
	}
	colName := tokens[*i].StringVal
	*i++

	if *i >= len(tokens) {
		return "", errors.New("incomplete statement")
	}
	if tokens[*i].Type != UnquoteStringToken {
		return "", errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[*i].Type, UnquoteStringToken)
	}
	dataTypeStr := tokens[*i].StringVal

	kind, err := stringToKind(dataTypeStr)
	if err != nil {
		return "", err
	}

	schema[colName] = kind
	*i++

	if *i >= len(tokens) || cmpTks(*tokens[*i], TokenComma) {
		return "", nil
	}

	if !cmpTks(*tokens[*i], TokenPrimary) {
		return "", errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[*i], TokenPrimary)
	}
	*i++
	if *i >= len(tokens) {
		return "", errors.New("invalid create statement")
	}
	if !cmpTks(*tokens[*i], TokenKey) {
		return "", errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[*i], TokenKey)
	}
	*i++

	return colName, nil
}

func genSchema(tokens []*Token) (map[string]reflect.Kind, string, error) {
	var (
		schema = make(map[string]reflect.Kind)
		pk     string
	)

	// parse the token and generate the schema
	for i := 0; i < len(tokens); {
		tpk, err := parseField(tokens, schema, &i)
		if err != nil {
			return nil, "",
				errors.Errorf("failed to parse the field definition: %v", err)
		}

		if len(pk) == 0 {
			pk = tpk
			continue
		}

		// if primary key has been set
		if len(tpk) != 0 {
			return nil, "",
				errors.Errorf("duplicate primary key: %s and %s", pk, tpk)
		}
	}

	if pk == "" {
		return nil, "", errors.New("primary key is not set")
	}

	return schema, pk, nil
}

func parseCreateStatement(tokens []*Token) (*CreateStatement, error) {
	// skip the first token, i.e., "Create"
	i := 1
	if i == len(tokens) {
		return nil, errors.New("incomplete create statement")
	}
	if !cmpTks(*tokens[i], TokenTable) {
		return nil, errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[i], TokenTable)
	}
	i++

	// get the table name
	if i == len(tokens) {
		return nil, errors.New("incomplete create statement")
	}
	if tokens[i].Type != UnquoteStringToken {
		return nil, errors.Errorf("invalid token type: "+
			"got(%s), expect(%s)", tokens[i].Type, UnquoteStringToken)
	}
	table := tokens[i].StringVal
	i++

	// parse and get the schema
	if i == len(tokens) {
		return nil, errors.New("incomplete create statement")
	}
	if !cmpTks(*tokens[i], TokenLeftParen) {
		return nil, errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[i], TokenLeftParen)
	}
	i++
	start := i

	if i == len(tokens) {
		return nil, errors.New("incomplete create statement")
	}
	for !cmpTks(*tokens[i], TokenRightParen) && i != len(tokens) {
		i++
	}
	if i == len(tokens) {
		return nil, errors.New("incomplete create statement")
	}
	end := i

	schema, primaryKey, err := genSchema(tokens[start:end])
	if err != nil {
		return nil, errors.Errorf("failed to get schema %v", err)
	}

	return &CreateStatement{
		table:      table,
		schema:     schema,
		primaryKey: primaryKey,
	}, nil
}

func getColumnNames(tokens []*Token, i *int) ([]string, error) {
	*i++
	ret := []string{}
	isColumnName := true
	for ; !cmpTks(*tokens[*i], TokenRightParen); *i++ {
		if isColumnName {
			if !isUnquoteStringToken(tokens[*i]) {
				return nil, errors.Errorf("invalid token (%s)",
					tokens[*i].String())
			}
			ret = append(ret, tokens[*i].StringVal)
		} else {
			// check if the token is a comma
			if !cmpTks(*tokens[*i], TokenComma) {
				return nil, errors.Errorf("invalid token (%s)"+
					" expect(%s)",
					tokens[*i].String(),
					TokenComma.String())
			}
		}
		isColumnName = !isColumnName
	}
	return ret, nil
}

func getValues(tokens []*Token, i *int) ([]any, error) {
	*i++
	ret := []any{}
	isValue := true
	for ; !cmpTks(*tokens[*i], TokenRightParen); *i++ {
		if isValue {
			if isUnquoteStringToken(tokens[*i]) {
				return nil, errors.Errorf("invalid token (%s)",
					tokens[*i].String())
			}
			switch tokens[*i].Type {
			case IntegerToken:
				ret = append(ret, tokens[*i].IntegerVal)
			case FloatToken:
				ret = append(ret, tokens[*i].FloatVal)
			case BoolToken:
				ret = append(ret, tokens[*i].BoolVal)
			case StringToken:
				ret = append(ret, tokens[*i].StringVal)
			default:
				return nil, errors.Errorf("invalid token type(%s)",
					tokens[*i].Type.String())
			}
		} else {
			// check if the token is a comma
			if !cmpTks(*tokens[*i], TokenComma) {
				return nil, errors.Errorf("invalid token (%s)"+
					" expect(%s)",
					tokens[*i].String(),
					TokenComma.String())
			}
		}
		isValue = !isValue
	}
	return ret, nil
}

func parseInsertStatement(tokens []*Token) (*InsertStatement, error) {
	// skip the first token, i.e., INSERT
	i := 1
	if i == len(tokens) {
		return nil, errors.New("incomplete insert statement")
	}
	if !cmpTks(*tokens[i], TokenInto) {
		return nil, errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[i], TokenInto)
	}
	i++

	// get the table name
	if i == len(tokens) {
		return nil, errors.New("incomplete insert statement")
	}
	if tokens[i].Type != UnquoteStringToken {
		return nil, errors.Errorf("invalid token type: "+
			"got(%s/'%s'), expect(%s)",
			tokens[i].Type, tokens[i], UnquoteStringToken)
	}
	table := tokens[i].StringVal
	i++

	// get column names if specified
	cns := []string{}
	if i == len(tokens) {
		return nil, errors.New("incomplete insert statement")
	}

	if cmpTks(*tokens[i], TokenLeftParen) {
		var err error
		cns, err = getColumnNames(tokens, &i)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get column names")
		}
	}
	i++

	// check the VALUES keyword
	if i == len(tokens) {
		return nil, errors.New("incomplete insert statement")
	}
	if !cmpTks(*tokens[i], TokenValues) {
		return nil, errors.Errorf("invalid token: "+
			"got(%s), expect(%s)", tokens[i], TokenValues)
	}
	i++

	// get the values
	if i == len(tokens) {
		return nil, errors.New("incomplete insert statement")
	}
	if !cmpTks(*tokens[i], TokenLeftParen) {
		return nil, errors.Errorf("invalid token: "+
			"got(%s), expect(%s)", tokens[i], TokenLeftParen)
	}
	vs, err := getValues(tokens, &i)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get values")
	}

	if len(cns) != len(vs) {
		return nil, errors.Errorf("number columns(%d) "+
			"not equal to number values(%d)", len(cns), len(vs))
	}

	values := make(map[string]any)
	for i, c := range cns {
		values[c] = vs[i]
	}
	return &InsertStatement{
		table:  table,
		values: values,
	}, nil
}

func parseDeleteStatement(tokens []*Token) (*DeleteStatement, error) {
	// skip the first token, i.e., DELETE
	i := 1
	if i == len(tokens) {
		return nil, errors.New("incomplete delete statement")
	}
	if !cmpTks(*tokens[i], TokenFrom) {
		return nil, errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[i], TokenFrom)
	}
	i++

	// get the table name
	if i == len(tokens) {
		return nil, errors.New("incomplete delete statement")
	}
	if tokens[i].Type != UnquoteStringToken {
		return nil, errors.Errorf("invalid token type: "+
			"got(%s/'%s'), expect(%s)",
			tokens[i].Type, tokens[i], UnquoteStringToken)
	}
	table := tokens[i].StringVal
	i++

	// parse the WHERE clause
	var where *WhereClause
	if i == len(tokens) {
		return nil, errors.New("incomplete delete statement")
	}

	where = &WhereClause{}
	if !cmpTks(*tokens[i], TokenWhere) {
		return nil, errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[i], TokenWhere)
	}
	i++

	if i == len(tokens) {
		return nil, errors.New("incomplete delete statement")
	}

	if !isUnquoteStringToken(tokens[i]) {
		return nil, errors.New("invalid token")
	}
	where.field = tokens[i].StringVal
	i++

	if i == len(tokens) {
		return nil, errors.New("incomplete delete statement")
	}

	if !cmpTks(*tokens[i], TokenEqual) {
		return nil, errors.Errorf("invalid token: got(%s), expect(%s)",
			tokens[i], TokenEqual)
	}
	i++
	if i == len(tokens) {
		return nil, errors.New("incomplete statement")
	}

	switch tokens[i].Type {
	case StringToken:
		where.value = tokens[i].StringVal
	case IntegerToken:
		where.value = tokens[i].IntegerVal
	case FloatToken:
		where.value = tokens[i].FloatVal
	case BoolToken:
		where.value = tokens[i].BoolVal
	default:
		return nil, errors.New("invalid token")
	}

	return &DeleteStatement{
		table: table,
		where: where,
	}, nil
}

func parseDropStatement(tokens []*Token) (*DropStatement, error) {
	// skip the first token, i.e., "DROP"
	i := 1
	if i == len(tokens) {
		return nil, errors.New("incomplete statement")
	}
	if !cmpTks(*tokens[i], TokenTable) {
		return nil, errors.Errorf("invalid token: got(%s), expect(%s)",
			*tokens[i], TokenTable)
	}
	i++

	if i == len(tokens) {
		return nil, errors.New("incomplete statement")
	}
	if !isUnquoteStringToken(tokens[i]) {
		return nil, errors.Errorf("invalid token (%s)",
			tokens[i].String())
	}
	table := tokens[i].StringVal
	return &DropStatement{
		table: table,
	}, nil
}

func parse(tokens []*Token) (any, error) {
	if len(tokens) == 0 {
		return nil, errors.New("cannot parse an empty token slice")
	}

	if tokens[0].Type != KeyWordToken {
		return nil, errors.New("invalid input format: the first token is not a keyword")
	}

	switch tokens[0].KeyWordVal {
	case Select:
		return parseSelectStatement(tokens)
	case Create:
		return parseCreateStatement(tokens)
	case Insert:
		return parseInsertStatement(tokens)
	case Delete:
		return parseDeleteStatement(tokens)
	case Drop:
		return parseDropStatement(tokens)
	default:
		return nil, errors.Errorf("invalid input format: unsupported keyword %s",
			tokens[0].KeyWordVal.String())
	}
}
