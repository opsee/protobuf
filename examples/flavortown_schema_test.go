package flavortown

import (
	"github.com/graphql-go/graphql"
	dessert "github.com/opsee/protobuf/examples/dessert"
	opsee_types "github.com/opsee/protobuf/opseeproto/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchema(t *testing.T) {
	populatedMenu := &Menu{
		Items: []*LineItem{
			{
				Dish: &LineItem_Lunch{&Lunch{
					Name:        "hogslop",
					Description: []byte("disgusting"),
				}},
				PriceCents: 100,
				CreatedAt:  &opsee_types.Timestamp{100, 100},
				UpdatedAt:  &opsee_types.Timestamp{200, 200},
			},
			{
				Dish: &LineItem_TastyDessert{&dessert.Dessert{
					Name:      "coolwhip",
					Sweetness: 9,
				}},
				PriceCents: 50,
				CreatedAt:  &opsee_types.Timestamp{100, 100},
				UpdatedAt:  &opsee_types.Timestamp{200, 200},
				Nothing:    nil,
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
					... on flavortownLunch {
						name
						description
					}
					... on flavortown_dessertDessert {
						name
						sweetness
					}
				}
				price_cents
				created_at
				updated_at
				nothing {
					void
				}
			}
		}
	}`})

	if queryResponse.HasErrors() {
		t.Fatalf("graphql query errors: %#v\n", queryResponse.Errors)
	}

	lunchitem := populatedMenu.Items[0]
	assert.Equal(t, lunchitem.GetLunch().Name, getProp(queryResponse.Data, "menu", "items", 0, "dish", "name"))
	assert.Equal(t, string(lunchitem.GetLunch().Description), getProp(queryResponse.Data, "menu", "items", 0, "dish", "description"))
	assert.EqualValues(t, lunchitem.PriceCents, getProp(queryResponse.Data, "menu", "items", 0, "price_cents"))
	assert.EqualValues(t, lunchitem.CreatedAt.Millis(), getProp(queryResponse.Data, "menu", "items", 0, "created_at"))
	assert.EqualValues(t, lunchitem.UpdatedAt.Millis(), getProp(queryResponse.Data, "menu", "items", 0, "updated_at"))

	dessertitem := populatedMenu.Items[1]
	assert.Equal(t, dessertitem.GetTastyDessert().Name, getProp(queryResponse.Data, "menu", "items", 1, "dish", "name"))
	assert.EqualValues(t, dessertitem.GetTastyDessert().Sweetness, getProp(queryResponse.Data, "menu", "items", 1, "dish", "sweetness"))
	assert.EqualValues(t, dessertitem.PriceCents, getProp(queryResponse.Data, "menu", "items", 1, "price_cents"))
	assert.EqualValues(t, dessertitem.CreatedAt.Millis(), getProp(queryResponse.Data, "menu", "items", 0, "created_at"))
	assert.EqualValues(t, dessertitem.UpdatedAt.Millis(), getProp(queryResponse.Data, "menu", "items", 0, "updated_at"))
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
