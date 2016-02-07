// Code generated by protoc-gen-gogo.
// source: dessert.proto
// DO NOT EDIT!

/*
Package graphql is a generated protocol buffer package.

It is generated from these files:
	dessert.proto

It has these top-level messages:
	Dessert
*/
package graphql

import github_com_graphql_go_graphql "github.com/graphql-go/graphql"
import fmt "fmt"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/opsee/protobuf/gogogqlproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var DessertType *github_com_graphql_go_graphql.Object

func init() {
	DessertType = github_com_graphql_go_graphql.NewObject(github_com_graphql_go_graphql.ObjectConfig{
		Name:        "dessert",
		Description: "A delicious dessert dish on the menu",
		Fields: (github_com_graphql_go_graphql.FieldsThunk)(func() github_com_graphql_go_graphql.Fields {
			return github_com_graphql_go_graphql.Fields{
				"name": &github_com_graphql_go_graphql.Field{
					Type:        github_com_graphql_go_graphql.String,
					Description: "The name of the dish",
					Resolve: func(p github_com_graphql_go_graphql.ResolveParams) (interface{}, error) {
						switch obj := p.Source.(type) {
						case *Dessert:
							return obj.Name, nil
						}
						return nil, fmt.Errorf("field name not resolved")
					},
				},
				"sweetness": &github_com_graphql_go_graphql.Field{
					Type:        github_com_graphql_go_graphql.Int,
					Description: "How sweet is the dish, an integer between 0 and 10",
					Resolve: func(p github_com_graphql_go_graphql.ResolveParams) (interface{}, error) {
						switch obj := p.Source.(type) {
						case *Dessert:
							return obj.Sweetness, nil
						}
						return nil, fmt.Errorf("field sweetness not resolved")
					},
				},
			}
		}),
	})
}
