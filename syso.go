package syso

import (
	"encoding/json"
	"io"
	"os"

	"github.com/hallazzang/syso/pkg/coff"
	"github.com/hallazzang/syso/pkg/common"
	"github.com/hallazzang/syso/pkg/ico"
	"github.com/hallazzang/syso/pkg/rsrc"
	"github.com/pkg/errors"
)

// FileResource represents a file resource that can be found at Path.
type FileResource struct {
	ID   int
	Name string
	Path string
}

// Validate returns an error if the resource is invalid.
func (r *FileResource) Validate() error {
	if r.Path == "" {
		return errors.New("no file path given")
	} else if r.ID == 0 && r.Name == "" {
		return errors.New("neither id nor name given")
	} else if r.ID != 0 && r.Name != "" {
		return errors.New("id and name cannot be set together")
	} else if r.ID < 0 {
		return errors.Errorf("invalid id: %d", r.ID)
	}
	return nil
}

// Config is a syso config data.
type Config struct {
	Icons    []*FileResource
	Manifest *FileResource
}

// ParseConfig reads JSON-formatted syso config from r and returns Config object.
func ParseConfig(r io.Reader) (*Config, error) {
	var c Config
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return nil, errors.Wrap(err, "failed to decode JSON")
	}
	for i, icon := range c.Icons {
		if err := icon.Validate(); err != nil {
			return nil, errors.Wrapf(err, "failed to validate icon #%d", i)
		}
		for j, icon2 := range c.Icons[:i] {
			if icon.ID != 0 && icon2.ID != 0 && icon2.ID == icon.ID {
				return nil, errors.Errorf("icon #%d's id and icon #%d's id are same", i, j)
			} else if icon.Name != "" && icon2.Name != "" && icon2.Name == icon.Name {
				return nil, errors.Errorf("icon #%d's name and icon #%d's name are same", i, j)
			}
		}
	}
	if c.Manifest != nil {
		if err := c.Manifest.Validate(); err != nil {
			return nil, errors.Wrap(err, "failed to validate manifest")
		}
	}
	return &c, nil
}

// EmbedIcon embeds icon into c.
func EmbedIcon(c *coff.File, icon *FileResource) error {
	if err := icon.Validate(); err != nil {
		return errors.Wrap(err, "invalid icon")
	}
	r, err := getOrCreateRSRCSection(c)
	if err != nil {
		return errors.Wrap(err, "failed to get or create .rsrc section")
	}
	f, err := os.Open(icon.Path)
	if err != nil {
		return errors.Wrap(err, "failed to open icon file")
	}
	defer f.Close()
	icons, err := ico.DecodeAll(f)
	if err != nil {
		return errors.Wrap(err, "failed to decode icon file")
	}
	for i, img := range icons.Images {
		img.ID = findPossibleID(r, 1000)
		if err := r.AddIconByID(img.ID, img); err != nil {
			return errors.Wrapf(err, "failed to add icon image #%d", i)
		}
	}
	if icon.ID != 0 {
		err = r.AddIconGroupByID(icon.ID, icons)
	} else {
		err = r.AddIconGroupByName(icon.Name, icons)
	}
	if err != nil {
		return errors.Wrap(err, "failed to add icon group resource")
	}
	return nil
}

func EmbedManifest(c *coff.File, manifest *FileResource) error {
	if err := manifest.Validate(); err != nil {
		return errors.Wrap(err, "invalid manifest")
	}
	r, err := getOrCreateRSRCSection(c)
	if err != nil {
		return errors.Wrap(err, "failed to get or create .rsrc section")
	}
	f, err := os.Open(manifest.Path)
	if err != nil {
		return errors.Wrap(err, "failed to open manifest file")
	}
	defer f.Close()
	b, err := common.NewBlob(f)
	if err != nil {
		return err
	}
	if manifest.ID != 0 {
		err = r.AddManifestByID(manifest.ID, b)
	} else {
		err = r.AddManifestByName(manifest.Name, b)
	}
	if err != nil {
		return errors.Wrap(err, "failed to add manifest resource")
	}
	return nil
}

func getOrCreateRSRCSection(c *coff.File) (*rsrc.Section, error) {
	s, err := c.Section(".rsrc")
	if err != nil {
		if err == coff.ErrSectionNotFound {
			s = rsrc.New()
			if err := c.AddSection(s); err != nil {
				return nil, errors.New("failed to add new .rsrc section")
			}
		} else {
			return nil, errors.Wrap(err, "failed to get .rsrc section")
		}
	}
	r, ok := s.(*rsrc.Section)
	if !ok {
		return nil, errors.New("the .rsrc section is not a valid rsrc section")
	}
	return r, nil
}

func findPossibleID(r *rsrc.Section, from int) int {
	// TODO: is 65535 a good limit for resource id?
	for ; from < 65536; from++ {
		if !r.ResourceIDExists(from) {
			break
		}
	}
	return from
}
