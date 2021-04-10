package schema

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/otiai10/copy"
	"github.com/pterm/pterm"
)

//go:embed config.cue
var BaseConfig string

// Database stores information about the
// repository
type Database struct {
	runtime *cue.Runtime
	db      cue.Value
	config  *cue.Value
}

// NewDatabase creates a "world" struct to store
// records
func NewDatabase() (Database, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", "")

	if nil != err {
		return Database{}, err
	}

	database := Database{
		runtime: &cueRuntime,
		db:      cueInstance.Value(),
	}

	err = database.LoadConfig()
	if nil != err {
		return Database{}, err
	}

	return database, nil
}

func (d *Database) LoadConfig() error {
	configInstance, err := d.runtime.Compile("", BaseConfig)
	if err != nil {
		return err
	}

	configValue := configInstance.Value()

	localConfig, err := ioutil.ReadFile("repository.cue")
	if err != nil {
		return err
	}

	localConfigInstance, err := d.runtime.Compile("", localConfig)
	if err != nil {
		return err
	}

	mergedConfig := configValue.Unify(localConfigInstance.Value())
	if err = mergedConfig.Validate(cue.Concrete(true)); err != nil {
		return err
	}

	d.config = &mergedConfig

	return nil
}

func (d *Database) GetConfigString(key string) (string, error) {
	value, err := d.config.LookupField(key)
	if err != nil {
		return "", err
	}

	str, err := value.Value.String()
	if err != nil {
		return "", err
	}

	return str, nil
}

// Repository is a group of schemas
type Repository struct {
	Root      string
	Namespace string
	Output    string
	Schemas   []*Schema
}

// GetRepository returns the Repository
// described by the repository.cue file in the
// current directory
func GetRepository() (*Repository, error) {
	sdb, err := NewDatabase()
	if err != nil {
		return nil, err
	}
	build_dir, err := sdb.GetConfigString("output_dir")
	if err != nil {
		return nil, err
	}
	namespace, err := sdb.GetConfigString("namespace")
	if err != nil {
		return nil, err
	}
	reporoot, err := sdb.GetConfigString("repository_root")
	if err != nil {
		return nil, err
	}
	r := &Repository{
		Namespace: namespace,
		Root:      reporoot,
		Output:    build_dir,
	}
	err = r.load()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NewRepository creates a new repository root and writes
// the metadata information
func NewRepository(namespace, output, root string) (*Repository, error) {
	r := &Repository{
		Root:      root,
		Namespace: namespace,
		Output:    output,
	}
	// create the repository directory
	err := r.createRoot()
	if err != nil {
		return nil, err
	}
	// write the config file
	err = r.writeConfig()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Repository) load() error {
	// load schemas and versions recursively
	r.Schemas = make([]*Schema, 0)
	schemaPath := r.Root
	err := filepath.WalkDir(schemaPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			pterm.Error.Printf("failure accessing a path %q: %v\n", path, err)
			return err
		}
		// be friendly to our Windows neighbors :)
		paths := strings.Split(path, string(os.PathSeparator))
		if d.IsDir() {
			if d.Name() == r.Root {
				return nil
			}
			if d.Name() == r.Output {
				return nil
			}
			if len(paths) == 2 {
				// this is a schema
				// process
				s := &Schema{
					Namespace: r.Namespace,
					Name:      d.Name(),
				}
				r.Schemas = append(r.Schemas, s)
				return nil
			}

			if len(paths) == 3 {
				// this is a version

				// process
				v := &Version{
					Namespace: r.Namespace,
					Name:      d.Name(),
				}
				for _, s := range r.Schemas {
					if s.Name == paths[1] {
						v.Schema = paths[1]
						s.Versions = append(s.Versions, v)
					}
				}
				return nil
			}

		} else {
			// not a dir, must be file
			// we only care about files in
			// version directories
			fmt.Println("file", len(paths), path)
			if len(paths) == 4 {
				if d.Name() == "schema.cue" {
					bb, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					for _, s := range r.Schemas {
						if s.Name == paths[1] {
							for _, v := range s.Versions {
								if v.Name == paths[2] {
									buf := bytes.NewBuffer([]byte{})
									json.HTMLEscape(buf, bb)
									v.Definition = buf.String()
								}
							}
						}
					}
				}
			}
		}

		return nil
	})

	return err
}
func (r *Repository) writeConfig() error {
	b := &Config{
		RepositoryRoot:  r.Root,
		Namespace:       r.Namespace,
		OutputDirectory: r.Output,
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	configPath := path.Join(wd, "repository.cue")
	return os.WriteFile(configPath, bb, 0755)
}

func (r *Repository) createRoot() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	repoPath := path.Join(wd, r.Root)
	err = os.MkdirAll(repoPath, 0755)
	return err
}

// AddSchema creates a new directory for a schema
// and creates the first version of the schema.
func (r *Repository) AddSchema(name string) error {

	// create the schema directory
	schemaPath := path.Join(r.Root, name)
	err := os.MkdirAll(schemaPath, 0744)
	if err != nil {
		return err
	}

	// create the first version
	versionPath := path.Join(schemaPath, "v1")
	err = os.MkdirAll(versionPath, 0744)
	if err != nil {
		return err
	}

	// write the schema
	schemaFile := path.Join(versionPath, "schema.cue")
	return os.WriteFile(schemaFile, schemaCue, 0755)

}

// AddVersion creates a new version for the specified
// schema
func (r *Repository) AddVersion(schema string) error {

	var sch *Schema
	for _, s := range r.Schemas {
		if s.Name == schema {
			sch = s
		}
	}
	if sch == nil {
		return fmt.Errorf("schema %s not found", schema)
	}
	versions := len(sch.Versions)
	prevVersionPath := path.Join(r.Root, sch.Name, fmt.Sprintf("v%d", versions))
	pterm.Info.Printf("Schema %s has %d version(s)\n", sch.Name, versions)
	nextVersion := versions + 1
	nextVersionPath := path.Join(r.Root, sch.Name, fmt.Sprintf("v%d", nextVersion))
	err := os.MkdirAll(nextVersionPath, 0755)
	if err != nil {
		return err
	}
	err = copy.Copy(prevVersionPath, nextVersionPath)
	if err != nil {
		return err
	}
	return nil
}

// Build serializes the Repository object
// into a json file in the `Output` directory.
func (r *Repository) Build() error {

	buildDir := path.Join(r.Root, r.Output)
	buildFile := path.Join(buildDir, "manifest.json")

	err := os.MkdirAll(buildDir, 0755)
	if err != nil {
		return err
	}

	bb, err := json.Marshal(r)

	if err != nil {
		return err
	}
	err = os.WriteFile(buildFile, bb, 0755)
	if err != nil {
		return err
	}
	return nil
}
