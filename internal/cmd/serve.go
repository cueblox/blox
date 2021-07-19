package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

			// GraphQL API
			dataSets := engine.GetDataSets()
			graphqlObjects := graphql.Fields{}

			for _, dataSet := range dataSets {
				fmt.Println("Setting up %v dataset", dataSet.GetExternalName())

				graphqlFields := graphql.Fields{}
				graphqlFields, err := cuedb.CueValueToGraphQlField(dataSet.GetSchemaCue())
				if err != nil {
					cobra.CheckErr(err)
				}

				objType := graphql.NewObject(
					graphql.ObjectConfig{
						Name:   dataSet.GetExternalName(),
						Fields: graphqlFields,
					},
				)

				graphqlObjects[strings.ToLower(dataSet.GetExternalName())] = &graphql.Field{
					Type: objType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["id"].(string)

						if ok {
							dataSetName := p.Info.ReturnType.Name()

							fmt.Println("Fetching data for %v", dataSetName)
							data := engine.GetAllData(fmt.Sprintf("#%s", dataSetName))

							records := make(map[string]interface{})
							if err = data.Decode(&records); err != nil {
								fmt.Printf("FAILED: %v\n", err)
								return nil, err
							}

							for recordID, record := range records {
								if string(recordID) == id {
									return record, nil
								}
							}
						}

						return nil, nil
					},
				}

				graphqlObjects[fmt.Sprintf("all%vs", dataSet.GetExternalName())] = &graphql.Field{
					Type: graphql.NewList(objType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						dataSetName := p.Info.ReturnType.Name()

						fmt.Println("Fetching data for %v", dataSetName)
						data := engine.GetAllData(fmt.Sprintf("#%s", dataSetName))

						records := make(map[string]interface{})
						if err = data.Decode(&records); err != nil {
							return nil, err
						}

						values := []interface{}{}
						for _, value := range records {
							fmt.Println(value)
							values = append(values, value)
						}

						return values, nil
					}}
			}

			var queryType = graphql.NewObject(
				graphql.ObjectConfig{
					Name:   "Query",
					Fields: graphqlObjects,
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
				Schema:   &schema,
				Pretty:   true,
				GraphiQL: true,
			})

			http.Handle("/graphiql", h)

			fmt.Println("Server is running on port 8080")
			http.ListenAndServe(":8080", nil)
		},
	}

	root.cmd = cmd
	return root
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}
