package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/cueblox/blox/internal/cuedb"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func (s *Service) prepGraphQL() error {
	if !s.built {
		err := s.build()
		if err != nil {
			return err
		}
	}

	dag := s.engine.GetDataSetsDAG()
	nodes, _ := dag.GetDescendants("root")

	// GraphQL API
	graphqlObjects := map[string]cuedb.GraphQlObjectGlue{}
	graphqlFields := graphql.Fields{}
	keys := make([]string, 0, len(nodes))
	for k := range nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	vertexComplete := map[string]bool{}

	// iterate through all the root nodes
	for _, k := range keys {
		node := nodes[k]

		chNode, _, err := dag.DescendantsWalker(k)
		cobra.CheckErr(err)

		// first iterate through all the children
		// so dependencies are registered
		for nd := range chNode {
			n, err := dag.GetVertex(nd)
			cobra.CheckErr(err)
			if dg, ok := n.(*cuedb.DagNode); ok {
				_, ok := vertexComplete[dg.Name]
				if !ok {
					err := s.translateNode(n, graphqlObjects, graphqlFields)
					if err != nil {
						return err
					}
					vertexComplete[dg.Name] = true
				}
			}
		}
		// now process the parent node
		_, ok := vertexComplete[node.(*cuedb.DagNode).Name]
		if !ok {
			err := s.translateNode(node, graphqlObjects, graphqlFields)
			if err != nil {
				return err
			}
			vertexComplete[node.(*cuedb.DagNode).Name] = true
		}

	}

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphqlFields,
		})

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	if err != nil {
		return err
	}
	s.schema = &schema
	return nil
}

func (s *Service) translateNode(node interface{}, graphqlObjects map[string]cuedb.GraphQlObjectGlue, graphqlFields map[string]*graphql.Field) error {
	dataSet, _ := s.engine.GetDataSet(node.(*cuedb.DagNode).Name)

	var objectFields graphql.Fields
	objectFields, err := cuedb.CueValueToGraphQlField(graphqlObjects, dataSet.GetSchemaCue())
	if err != nil {
		cobra.CheckErr(err)
	}

	// Inject ID field into each object
	objectFields["id"] = &graphql.Field{
		Type: &graphql.NonNull{
			OfType: graphql.String,
		},
	}

	objType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   dataSet.GetExternalName(),
			Fields: objectFields,
		},
	)

	resolver := func(p graphql.ResolveParams) (interface{}, error) {
		dataSetName := p.Info.ReturnType.Name()

		id, ok := p.Args["id"].(string)
		if ok {
			data := s.engine.GetAllData(fmt.Sprintf("#%s", dataSetName))

			records := make(map[string]interface{})
			if err = data.Decode(&records); err != nil {
				return nil, err
			}

			for recordID, record := range records {
				if string(recordID) == id {
					return record, nil
				}
			}
		}
		return nil, nil
	}

	graphqlObjects[dataSet.GetExternalName()] = cuedb.GraphQlObjectGlue{
		Object:   objType,
		Resolver: resolver,
		Engine:   s.engine,
	}

	graphqlFields[dataSet.GetExternalName()] = &graphql.Field{
		Name: dataSet.GetExternalName(),
		Type: objType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: resolver,
	}

	graphqlFields[fmt.Sprintf("all%vs", dataSet.GetExternalName())] = &graphql.Field{
		Type: graphql.NewList(objType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			dataSetName := p.Info.ReturnType.Name()

			data := s.engine.GetAllData(fmt.Sprintf("#%s", dataSetName))

			records := make(map[string]interface{})
			if err = data.Decode(&records); err != nil {
				return nil, err
			}

			values := []interface{}{}
			for _, value := range records {
				values = append(values, value)
			}

			return values, nil
		},
	}
	return nil
}

// GQLHandlerFunc returns a stand alone graphql handler for use
// in netlify/aws/azure serverless scenarios
func (s *Service) GQLHandlerFunc() (http.HandlerFunc, error) {
	if s.schema == nil {
		err := s.prepGraphQL()
		if err != nil {
			return nil, err
		}
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		result := s.executeQuery(r.URL.Query().Get("query"))
		err := json.NewEncoder(w).Encode(result)
		if err != nil {
			pterm.Warning.Printf("failed to encode: %v", err)
		}
	}
	return hf, nil
}

// GQLPlaygroundHandler returns a stand alone graphql playground handler for use
// in netlify/aws/azure serverless scenarios
func (s *Service) GQLPlaygroundHandler() (http.Handler, error) {
	if s.schema == nil {
		err := s.prepGraphQL()
		if err != nil {
			return nil, err
		}
	}

	h := handler.New(&handler.Config{
		Schema:     s.schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})
	return h, nil
}

func (s *Service) executeQuery(query string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        *s.schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		pterm.Error.Printf("errors: %v\n", result.Errors)
	}

	return result
}
