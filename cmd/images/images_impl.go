package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cueblox/blox/plugins"
	"github.com/disintegration/imaging"
	"github.com/goccy/go-yaml"
	"github.com/h2non/filetype"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/pterm/pterm"
)

// Here is a real implementation of Greeter
type ImageScanner struct {
	logger hclog.Logger
}

func (g *ImageScanner) Process() error {
	pterm.Info.Println("PROCESS()")
	g.logger.Debug("message from ImageScanner.Process")
	return g.processImages()
}

func (g *ImageScanner) processImages() error {
	// hard coded, need to pass in config? read config?
	staticDir := "static"

	g.logger.Debug("processing images", "dir", staticDir)
	fi, err := os.Stat(staticDir)
	if errors.Is(err, os.ErrNotExist) {
		g.logger.Debug("no image directory found, skipping")
		return nil
	}
	if !fi.IsDir() {
		return errors.New("given static directory is not a directory")
	}
	imagesDirectory := filepath.Join(staticDir, "images")

	fi, err = os.Stat(imagesDirectory)
	if errors.Is(err, os.ErrNotExist) {
		g.logger.Debug("no image directory found, skipping")
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

			g.logger.Debug("Processing", "path", path)
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
					g.logger.Debug("File is an image", "path", relpath)
					kind, err := filetype.Match(buf)
					if err != nil {
						return err
					}
					g.logger.Debug("\t\tFile type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
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
					dataDir := "data"

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
							g.logger.Debug(err.Error())
							return err
						}
					}
				} else {
					g.logger.Debug("File is not an image",
						"path", path)
				}
			}

			return nil
		})
	return err
}

type BloxImage struct {
	FileName string `yaml:"file_name"`
	Height   int    `yaml:"height"`
	Width    int    `yaml:"width"`
	CDN      string `yaml:"cdn"`
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BLOX_PLUGIN",
	MagicCookieValue: "image",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	imageScanner := &ImageScanner{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"images": &plugins.PrebuildPlugin{Impl: imageScanner},
	}

	logger.Debug("message from plugin", "hello", "images")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
