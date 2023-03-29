package dbx_model

import (
	"github.com/go-web-kits/utils"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

var DBxDefinitions map[string]Definition

type Definition struct {
	Uniqueness   interface{}
	DefaultScope Scope
	Serialization
}

type Serialization struct {
	Add map[string]string
	Rmv []string
}

type Scope struct {
	Order         string
	Where         interface{}
	Preload       interface{}
	Join          []string
	OrderBeJoined string
	WhereBeJoined []interface{}
}

func DefinitionOf(obj interface{}) Definition {
	return DBxDefinitions[NameOf(obj)]
}

func NameOf(obj interface{}) string {
	if name, ok := obj.(string); ok {
		return inflection.Singular(strcase.ToCamel(name))
	}
	return strcase.ToCamel(utils.TypeNameOf(obj))
}

func TableNameOf(obj interface{}) string {
	return strcase.ToSnake(inflection.Plural(NameOf(obj)))
}
