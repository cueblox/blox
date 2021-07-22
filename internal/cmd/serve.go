package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/cuedb"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	// Import the blob packages we want to be able to open.
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

type bloxServeCmd struct {
	cmd *cobra.Command
}

func newBloxServeCmd() *bloxServeCmd {
	root := &bloxServeCmd{}
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a GraphQL API",

		Run: func(cmd *cobra.Command, args []string) {
			userConfig, err := ioutil.ReadFile("blox.cue")

			pterm.Debug.Printf("loading user config")

			cobra.CheckErr(err)

			engine, err := cuedb.NewEngine()

			pterm.Debug.Printf("new engine")
			cobra.CheckErr(err)

			cfg, err := blox.NewConfig(BaseConfig)

			pterm.Debug.Printf("newConfig")
			cobra.CheckErr(err)

			err = cfg.LoadConfigString(string(userConfig))
			cobra.CheckErr(err)

			// Load Schemas!
			schemataDir, err := cfg.GetString("schemata_dir")
			cobra.CheckErr(err)

			remotes, err := cfg.GetList("remotes")
			if err == nil {
				cobra.CheckErr(parseRemotes(remotes))
			}
			if images {
				err = processImages(cfg)
				if err != nil {
					cobra.CheckErr(err)
				}
			}
			pterm.Debug.Printf("\t\tUsing schemata from: %s\n", schemataDir)

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
			cobra.CheckErr(err)

			pterm.Debug.Println("\t\tBuilding DataSets")
			cobra.CheckErr(buildDataSets(engine, cfg))

			if referentialIntegrity {
				pterm.Info.Println("Verifying Referential Integrity")
				cobra.CheckErr(engine.ReferentialIntegrity())
				pterm.Success.Println("Referential Integrity OK")
			}

			if static {
				staticDir, err := cfg.GetString("static_dir")
				pterm.Info.Println("Serving static files from", staticDir)

				cobra.CheckErr(err)
				http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(".", staticDir)))))
			}

			dag := engine.GetDataSetsDAG()
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
							dataSet, _ := engine.GetDataSet(dg.Name)

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
									data := engine.GetAllData(fmt.Sprintf("#%s", dataSetName))

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
								Engine:   engine,
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

									data := engine.GetAllData(fmt.Sprintf("#%s", dataSetName))

									records := make(map[string]interface{})
									if err = data.Decode(&records); err != nil {
										return nil, err
									}

									values := []interface{}{}
									for _, value := range records {
										values = append(values, value)
									}

									return values, nil
								}}
							vertexComplete[dg.Name] = true
						}
					}
				}
				_, ok := vertexComplete[node.(*cuedb.DagNode).Name]
				if !ok {
					dataSet, _ := engine.GetDataSet(node.(*cuedb.DagNode).Name)

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
							data := engine.GetAllData(fmt.Sprintf("#%s", dataSetName))

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
						Engine:   engine,
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

							data := engine.GetAllData(fmt.Sprintf("#%s", dataSetName))

							records := make(map[string]interface{})
							if err = data.Decode(&records); err != nil {
								return nil, err
							}

							values := []interface{}{}
							for _, value := range records {
								values = append(values, value)
							}

							return values, nil
						}}
				}

			}

			var queryType = graphql.NewObject(
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
				cobra.CheckErr(err)
			}

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				result := executeQuery(r.URL.Query().Get("query"), schema)
				json.NewEncoder(w).Encode(result)
			})

			h := handler.New(&handler.Config{
				Schema:     &schema,
				Pretty:     true,
				GraphiQL:   false,
				Playground: true,
			})

			http.Handle("/ui", h)

			pterm.Info.Printf("Server is running at %s\n", address)
			http.ListenAndServe(address, nil)
		},
	}
	cmd.Flags().BoolVarP(&static, "static", "s", true, "Serve static files")
	cmd.Flags().StringVarP(&address, "address", "a", ":8080", "Listen address")

	root.cmd = cmd
	return root
}

var static bool
var address string

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		pterm.Error.Printf("errors: %v\n", result.Errors)
	}

	return result
}
