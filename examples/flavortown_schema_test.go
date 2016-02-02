package flavortown

import (
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	google_protobuf "google/protobuf"
	"testing"
)

func TestSchema(t *testing.T) {
	populatedMenu := &Menu{
		Items: []*LineItem{
			{
				Dish: &Dish{
					Name:        "hogslop",
					Description: []byte("disgusting"),
				},
				PriceCents: 100,
				CreatedAt:  &google_protobuf.Timestamp{100, 100},
				UpdatedAt:  &google_protobuf.Timestamp{200, 200},
			},
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"menu": &graphql.Field{
					Type: GraphQLMenuType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return populatedMenu, nil
					},
				},
			},
		}),
	})

	if err != nil {
		t.Fatal(err)
	}

	queryResponse := graphql.Do(graphql.Params{Schema: schema, RequestString: `query goodQuery {
		menu {
			items {
				dish {
					name
					description
				}
				price_cents
				created_at
				updated_at
			}
		}
	}`})

	if queryResponse.HasErrors() {
		t.Fatalf("graphql query errors: %#v\n", queryResponse.Errors)
	}

	item := populatedMenu.Items[0]
	assert.Equal(t, item.Dish.Name, getProp(queryResponse.Data, "menu", "items", 0, "dish", "name"))
	assert.Equal(t, string(item.Dish.Description), getProp(queryResponse.Data, "menu", "items", 0, "dish", "description"))
	assert.EqualValues(t, item.PriceCents, getProp(queryResponse.Data, "menu", "items", 0, "price_cents"))
	assert.EqualValues(t, item.CreatedAt.String(), getProp(queryResponse.Data, "menu", "items", 0, "created_at"))
	assert.EqualValues(t, item.UpdatedAt.String(), getProp(queryResponse.Data, "menu", "items", 0, "updated_at"))
}

func getProp(i interface{}, path ...interface{}) interface{} {
	cur := i

	for _, s := range path {
		switch cur.(type) {
		case map[string]interface{}:
			cur = cur.(map[string]interface{})[s.(string)]
			continue
		case []interface{}:
			cur = cur.([]interface{})[s.(int)]
			continue
		default:
			return cur
		}
	}

	return cur
}
