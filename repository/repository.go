package repository

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cueblox/blox"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/spf13/cobra"

	"github.com/cueblox/blox/internal/cuedb"
	"github.com/cueblox/blox/internal/encoding/markdown"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
)

type Service struct {
	engine *cuedb.Engine
	Cfg    *blox.Config
	ri     bool
	schema *graphql.Schema
}

func NewService(bloxConfig string, referentialIntegrity bool) (*Service, error) {
	cfg, err := blox.NewConfig(BaseConfig)
	if err != nil {
		return nil, err
	}

	err = cfg.LoadConfigString(bloxConfig)
	if err != nil {
		return nil, err
	}

	engine, err := cuedb.NewEngine()
	if err != nil {
		return nil, err
	}

	schemataDir, err := cfg.GetString("schemata_dir")
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(schemataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			bb, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			pterm.Debug.Printf("\t\tAttempting to register schema: %s\n", path)
			err = engine.RegisterSchema(string(bb))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		engine: engine,
		Cfg:    cfg,
		ri:     referentialIntegrity,
	}, nil
}

const BaseConfig = `{
	#Remote: {
	    name: string
	    version: string
	    repository: string
	}
	build_dir:    string | *"_build"
	data_dir:     string | *"data"
	schemata_dir: string | *"schemata"
	    static_dir: string | *"static"
	    template_dir: string | *"templates"
	    remotes: [ ...#Remote ]
    
    }`

func (s *Service) RenderJSON() ([]byte, error) {
	err := s.build()
	if err != nil {
		return nil, err
	}
	pterm.Debug.Println("Building output data blox")
	output, err := s.engine.GetOutput()
	if err != nil {
		return nil, err
	}

	pterm.Debug.Println("Rendering data blox to JSON")
	return output.MarshalJSON()
}

func (s *Service) RenderAndSave() error {
	err := s.build()
	if err != nil {
		return err
	}

	bb, err := s.RenderJSON()
	if err != nil {
		return err
	}
	buildDir, err := s.Cfg.GetString("build_dir")
	if err != nil {
		return err
	}
	err = os.MkdirAll(buildDir, 0o755)
	if err != nil {
		return err
	}
	filename := "data.json"
	filePath := path.Join(buildDir, filename)
	err = os.WriteFile(filePath, bb, 0o755)
	if err != nil {
		return err
	}
	pterm.Success.Printf("Data blox written to '%s'\n", filePath)
	return nil
}

func (s *Service) build() error {
	var errors error

	for _, dataSet := range s.engine.GetDataSets() {
		pterm.Debug.Printf("\t\tBuilding Dataset: %s\n", dataSet.ID())

		// We're using the Or variant of GetString because we know this call can't
		// fail, as the config isn't valid without.
		dataSetDirectory := fmt.Sprintf("%s/%s", s.Cfg.GetStringOr("data_dir", ""), dataSet.GetDataDirectory())

		err := os.MkdirAll(dataSetDirectory, 0o755)
		if err != nil {
			errors = multierror.Append(err)
			continue
		}

		err = filepath.Walk(dataSetDirectory,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return err
				}

				ext := strings.TrimPrefix(filepath.Ext(path), ".")

				if !dataSet.IsSupportedExtension(ext) {
					return nil
				}
				relative, err := filepath.Rel(dataSetDirectory, path)
				if err != nil {
					return err
				}
				slug := strings.TrimSuffix(relative, "."+ext)
				pterm.Debug.Println(slug)
				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					return multierror.Append(err)
				}

				// Loaders to get to YAML
				// We should offer various, simple for now with markdown
				mdStr := ""
				if ext == "md" || ext == "mdx" {
					mdStr, err = markdown.ToYAML(string(bytes))
					if err != nil {
						return err
					}

					bytes = []byte(mdStr)
				}

				istruct := make(map[string]interface{})

				err = yaml.Unmarshal(bytes, &istruct)
				if err != nil {
					return multierror.Append(err)
				}

				record := make(map[string]interface{})
				record[slug] = istruct

				err = s.engine.Insert(dataSet, record)
				if err != nil {
					return multierror.Append(err)
				}

				return err
			},
		)

		if err != nil {
			errors = multierror.Append(err)
		}
	}

	if s.ri {
		err := s.engine.ReferentialIntegrity()
		if err != nil {
			errors = multierror.Append(err)
		}
	}

	if errors != nil {
		pterm.Error.Println("Validation Failed")
		return errors
	}

	pterm.Success.Println("Validation Complete")
	return nil
}

func (s *Service) prepGraphQL() error {

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

	for _, k := range keys {
		node := nodes[k]

		chNode, _, err := dag.DescendantsWalker(k)
		cobra.CheckErr(err)

		for nd := range chNode {
			n, err := dag.GetVertex(nd)
			cobra.CheckErr(err)
			if dg, ok := n.(*cuedb.DagNode); ok {
				_, ok := vertexComplete[dg.Name]
				if !ok {
					dataSet, _ := s.engine.GetDataSet(dg.Name)

					var objectFields graphql.Fields
					objectFields, err = cuedb.CueValueToGraphQlField(graphqlObjects, dataSet.GetSchemaCue())

					if err != nil {
						cobra.CheckErr(err)
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
					vertexComplete[dg.Name] = true
				}
			}
		}
		_, ok := vertexComplete[node.(*cuedb.DagNode).Name]
		if !ok {
			dataSet, _ := s.engine.GetDataSet(node.(*cuedb.DagNode).Name)

			var objectFields graphql.Fields
			objectFields, err = cuedb.CueValueToGraphQlField(graphqlObjects, dataSet.GetSchemaCue())

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