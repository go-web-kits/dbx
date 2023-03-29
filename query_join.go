package dbx

import (
	"reflect"
	"strings"

	"github.com/go-web-kits/utils"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/inflection"
)

func (c *Chain) joins(opt Opt) *Chain {
	if len(opt.Join) > 0 {
		return c.Joins(opt.Join[0], opt.Join[1:]...)
	}
	return c
}

// @doc c.4.1 Query Chain #8
func (c *Chain) Joins(assocFieldName string, nestedAssocFieldNames ...string) *Chain {
	obj := c.GetModel()
	c, field := c.doJoins(obj, assocFieldName)
	for _, fieldName := range nestedAssocFieldNames {
		c, field = c.doJoins(reflect.New(field.Struct.Type).Interface(), fieldName)
	}

	return c
}

func (c *Chain) doJoins(source interface{}, assocFieldName string) (*Chain, *gorm.Field) {
	scope := c.InstantSet("gorm:association:source", source).NewScope(source)
	field, ok := scope.FieldByName(assocFieldName)
	if ok {
		if field.Relationship == nil || len(field.Relationship.ForeignFieldNames) == 0 {
			panic("Join: invalid association " + assocFieldName)
		}
	} else {
		panic("Join: no such column " + assocFieldName)
	}

	sourceTable := strcase.ToSnake(inflection.Plural(utils.TypeNameOf(source)))
	joinedModel := utils.TypeNameOf(field.Struct.Type.String())
	destTable := strcase.ToSnake(inflection.Plural(joinedModel))

	if field.Relationship.PolymorphicType != "" {
		c = c.joinsByPolymorphic(sourceTable, destTable, field)
	} else if field.Relationship.JoinTableHandler != nil {
		c = c.joinsByJoinTable(sourceTable, destTable, field)
	} else {
		c = c.joinsNormally(sourceTable, destTable, field)
	}

	return c.ScopingJoined(joinedModel), field
}

func (c *Chain) joinsByPolymorphic(sourceTable, destTable string, field *gorm.Field) *Chain {
	r := field.Relationship
	sourceKey := r.AssociationForeignDBNames[0]
	foreignKey := r.ForeignDBNames[0]

	return &Chain{c.DB.Joins("/*\n*/ LEFT JOIN \""+destTable+"\" ON \""+
		destTable+"\".\""+foreignKey+"\" = \""+sourceTable+"\".\""+sourceKey+"\" AND \""+
		destTable+"\".\""+r.PolymorphicDBName+"\" = ?", r.PolymorphicValue)}
}

func (c *Chain) joinsByJoinTable(sourceTable, destTable string, field *gorm.Field) *Chain {
	jth := field.Relationship.JoinTableHandler
	joinTable := jth.Table(c.DB)
	source := jth.SourceForeignKeys()[0]
	dest := jth.DestinationForeignKeys()[0]

	return &Chain{c.DB.
		Joins("/*\n*/ LEFT JOIN \"" + joinTable + "\" ON \"" + joinTable + "\".\"" + source.DBName + "\" = \"" + sourceTable + "\".\"" + source.AssociationDBName + "\"").
		Joins("/*\n*/ LEFT JOIN \"" + destTable + "\" ON \"" + joinTable + "\".\"" + dest.DBName + "\" = \"" + destTable + "\".\"" + dest.AssociationDBName + "\""),
	}
}

func (c *Chain) joinsNormally(sourceTable, destTable string, field *gorm.Field) *Chain {
	r := field.Relationship
	assocKey := r.AssociationForeignDBNames[0]
	foreignKey := r.ForeignDBNames[0]

	if strings.Contains(r.Kind, "has") {
		return &Chain{c.DB.Joins("/*\n*/ LEFT JOIN \"" + destTable + "\" ON \"" + sourceTable + "\".\"" + assocKey + "\" = \"" + destTable + "\".\"" + foreignKey + "\"")}
	} else if strings.Contains(r.Kind, "belongs") {
		return &Chain{c.DB.Joins("/*\n*/ LEFT JOIN \"" + destTable + "\" ON \"" + sourceTable + "\".\"" + foreignKey + "\" = \"" + destTable + "\".\"" + assocKey + "\"")}
	}
	return c
}
