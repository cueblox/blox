package content

import (
	"encoding/json"
	"errors"
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
	"github.com/disintegration/imaging"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/h2non/filetype"
	"github.com/spf13/cobra"

	"github.com/cueblox/blox/internal/cuedb"
	"github.com/cueblox/blox/internal/encoding/markdown"
	"github.com/cueblox/blox/internal/repository"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
)

type Service struct {
	engine *cuedb.Engine
	Cfg    *blox.Config
	ri     bool
	schema *graphql.Schema
	built  bool
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

	err := s.parseRemotes()
	if err != nil {
		return err
	}
	err = s.processImages()
	if err != nil {
		return err
	}

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
	s.built = true
	return nil
}

// RestHandlerFunc returns a handler function that will render
// the dataset specified as the last path parameter.
func (s *Service) RestHandlerFunc() (http.HandlerFunc, error) {
	if !s.built {
		err := s.build()
		if err != nil {
			return nil, err
		}
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.Split(path, "/")
		pterm.Debug.Println(parts, len(parts))
		if len(parts) == 0 {
			pterm.Warning.Println("No dataset specified")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		dataset := parts[len(parts)-1]

		ds, err := s.engine.GetDataSetByPlural(dataset)
		if err != nil {
			pterm.Warning.Println("Requested dataset not found", parts, len(parts))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		data := s.engine.GetAllData(ds.GetExternalName())
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			pterm.Warning.Printf("failed to encode: %v", err)
		}
	}
	return hf, nil
}

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

func (s *Service) parseRemotes() error {
	remotes, err := s.Cfg.GetList("remotes")
	if err != nil {
		return err
	}
	iter, err := remotes.List()
	if err != nil {
		return err
	}
	for iter.Next() {
		val := iter.Value()

		//nolint
		name, err := val.FieldByName("name", false)
		if err != nil {
			return err
		}
		n, err := name.Value.String()
		if err != nil {
			return err
		}
		//nolint
		version, err := val.FieldByName("version", false)
		if err != nil {
			return err
		}
		v, err := version.Value.String()
		if err != nil {
			return err
		}
		//nolint
		repository, err := val.FieldByName("repository", false)
		if err != nil {
			return err
		}
		r, err := repository.Value.String()
		if err != nil {
			return err
		}
		err = s.ensureRemote(n, v, r)
		if err != nil {
			return err
		}
	}
	return nil
}

// processImages scans the static dir for images
// when it finds an image, it reads the image metadata
// and saves a corresponding YAML file describing the image
// in the 'images' data directory.
func (s *Service) processImages() error {
	staticDir, err := s.Cfg.GetString("static_dir")
	if err != nil {
		pterm.Info.Printf("no static directory present, skipping image linking")
		return nil
	}

	pterm.Info.Printf("processing images in %s\n", staticDir)
	fi, err := os.Stat(staticDir)
	if errors.Is(err, os.ErrNotExist) {
		pterm.Info.Println("no image directory found, skipping")
		return nil
	}
	if !fi.IsDir() {
		return errors.New("given static directory is not a directory")
	}
	imagesDirectory := filepath.Join(staticDir, "images")

	fi, err = os.Stat(imagesDirectory)
	if errors.Is(err, os.ErrNotExist) {
		pterm.Info.Println("no image directory found, skipping")
		return nil
	}
	if !fi.IsDir() {
		return errors.New("given images directory is not a directory")
	}
	err = filepath.Walk(imagesDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			pterm.Debug.Printf("\t\tProcessing %s\n", path)
			if !info.IsDir() {
				buf, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				if filetype.IsImage(buf) {

					src, err := imaging.Open(path)
					if err != nil {
						return err
					}

					relpath, err := filepath.Rel(staticDir, path)
					if err != nil {
						return err
					}
					pterm.Debug.Printf("\t\tFile is an image: %s\n", relpath)
					kind, err := filetype.Match(buf)
					if err != nil {
						return err
					}
					pterm.Debug.Printf("\t\tFile type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
					if err != nil {
						return err
					}
					/*	if cloud {
							pterm.Info.Println("Synchronizing images to cloud provider")
							bucketURL := os.Getenv("IMAGE_BUCKET")
							if bucketURL == "" {
								return errors.New("image sync enabled (-s,--sync), but no IMAGE_BUCKET environment variable set")
							}

							ctx := context.Background()
							// Open a connection to the bucket.
							b, err := blob.OpenBucket(ctx, bucketURL)
							if err != nil {
								return fmt.Errorf("failed to setup bucket: %s", err)
							}
							defer b.Close()

							w, err := b.NewWriter(ctx, relpath, nil)
							if err != nil {
								return fmt.Errorf("cloud sync failed to obtain writer: %s", err)
							}
							_, err = w.Write(buf)
							if err != nil {
								return fmt.Errorf("cloud sync failed to write to bucket: %s", err)
							}
							if err = w.Close(); err != nil {
								return fmt.Errorf("cloud sync failed to close: %s", err)
							}
						}
					*/
					cdnEndpoint := os.Getenv("CDN_URL")

					bi := &BloxImage{
						FileName: relpath,
						Height:   src.Bounds().Dy(),
						Width:    src.Bounds().Dx(),
						CDN:      cdnEndpoint,
					}
					bytes, err := yaml.Marshal(bi)
					if err != nil {
						return err
					}
					dataDir, err := s.Cfg.GetString("data_dir")
					if err != nil {
						return err
					}

					ext := strings.TrimPrefix(filepath.Ext(relpath), ".")
					slug := strings.TrimSuffix(relpath, "."+ext)

					outputPath := filepath.Join(dataDir, slug+".yaml")
					err = os.MkdirAll(filepath.Dir(outputPath), 0o755)
					if err != nil {
						pterm.Error.Println(err)
						return err
					}
					// only write the yaml file if it doesn't exist.
					// don't overwrite existing records.
					_, err = os.Stat(outputPath)
					if err != nil && errors.Is(err, os.ErrNotExist) {
						err = os.WriteFile(outputPath, bytes, 0o755)
						if err != nil {
							pterm.Error.Println(err)
							return err
						}
					}
				} else {
					pterm.Debug.Printf("File is not an image: %s\n", path)
				}
			}

			return nil
		})
	return err
}

// ensureRemote downloads a remote schema at a specific
// version if it doesn't exist locally
func (s *Service) ensureRemote(name, version, repo string) error {
	// Load Schemas!
	schemataDir, err := s.Cfg.GetString("schemata_dir")
	if err != nil {
		return err
	}

	schemaFilePath := path.Join(schemataDir, fmt.Sprintf("%s_%s.cue", name, version))
	_, err = os.Stat(schemaFilePath)
	if os.IsNotExist(err) {
		pterm.Info.Printf("Schema does not exist locally: %s_%s.cue\n", name, version)
		manifest := fmt.Sprintf("https://%s/manifest.json", repo)
		res, err := http.Get(manifest)
		if err != nil {
			return err
		}

		var repos repository.Repository
		err = json.NewDecoder(res.Body).Decode(&repos)
		if err != nil {
			return err
		}

		var selectedVersion *repository.Version
		for _, s := range repos.Schemas {
			if s.Name == name {
				for _, v := range s.Versions {
					if v.Name == version {
						selectedVersion = v
						pterm.Debug.Println(v.Name, v.Schema, v.Definition)
					}
				}
			}
		}

		// make schemata directory
		err = os.MkdirAll(schemataDir, 0o755)
		if err != nil {
			return err
		}

		// TODO: don't overwrite each time
		filename := fmt.Sprintf("%s_%s.cue", name, version)
		filePath := path.Join(schemataDir, filename)
		err = os.WriteFile(filePath, []byte(selectedVersion.Definition), 0o755)
		if err != nil {
			return err
		}
		pterm.Info.Printf("Schema downloaded: %s_%s.cue\n", name, version)
		return nil
	}
	pterm.Info.Println("Schema already exists locally, skipping download.")
	return nil
}

type BloxImage struct {
	FileName string `yaml:"file_name"`
	Height   int    `yaml:"height"`
	Width    int    `yaml:"width"`
	CDN      string `yaml:"cdn"`
}
