package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/cueblox/blox"
	"github.com/cueblox/blox/repository"

	"github.com/disintegration/imaging"
	"github.com/goccy/go-yaml"
	"github.com/h2non/filetype"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"gocloud.dev/blob"

	// Import the blob packages we want to be able to open.
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

type bloxBuildCmd struct {
	cmd *cobra.Command
}

func newBloxBuildCmd() *bloxBuildCmd {
	root := &bloxBuildCmd{}
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Validate & Build dataset",
		Long: `The build command will ensure that your dataset is correct by
	validating it against your schemata. Once validated, it will render all
	your content into a single JSON file, which can be consumed by your tooling
	of choice.
	
	Referential Integrity can be enforced with -i. This ensures that any fields
	ending with _id are valid references to identifiers within the other content type.
	
	The build process will create an 'image' record for images in your 'static_dir' if you use the -g or --images flag.
	
	Images will be pushed to blob storage if you use -s/--sync and set the appropriate environment variables. 
	Currently only Azure blob storate is implemented. See https://gocloud.dev/howto/blob/#services for required 
	environment variables and setup information.
	`,
		Run: func(cmd *cobra.Command, args []string) {
			userConfig, err := ioutil.ReadFile("blox.cue")

			pterm.Debug.Printf("loading user config")

			cobra.CheckErr(err)

			repo, err := repository.NewService(string(userConfig), referentialIntegrity)
			cobra.CheckErr(err)

			remotes, err := repo.Cfg.GetList("remotes")
			if err == nil {
				cobra.CheckErr(parseRemotes(remotes))
			}
			if images {
				err = processImages(repo.Cfg)
				if err != nil {
					cobra.CheckErr(err)
				}
			}

			err = repo.RenderAndSave()
			cobra.CheckErr(err)

		},
	}
	cmd.Flags().BoolVarP(&referentialIntegrity, "referential-integrity", "i", false, "Verify referential integrity")
	cmd.Flags().BoolVarP(&images, "images", "g", false, "Create 'image' records for images in static directory")
	cmd.Flags().BoolVarP(&cloud, "sync", "s", false, "Sync images to blob storage")
	root.cmd = cmd
	return root
}

var (
	referentialIntegrity bool
	images               bool
	cloud                bool
)

const DefaultConfigName = "blox.cue"

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

func parseRemotes(value cue.Value) error {
	iter, err := value.List()
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
		err = ensureRemote(n, v, r)
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
func processImages(cfg *blox.Config) error {
	staticDir, err := cfg.GetString("static_dir")
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
					if cloud {
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
					dataDir, err := cfg.GetString("data_dir")
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
	cobra.CheckErr(err)
	return nil
}

type BloxImage struct {
	FileName string `yaml:"file_name"`
	Height   int    `yaml:"height"`
	Width    int    `yaml:"width"`
	CDN      string `yaml:"cdn"`
}
