package syso

import (
	"encoding/json"
	"io"
	"os"

	"github.com/hallazzang/syso/pkg/coff"
	"github.com/hallazzang/syso/pkg/ico"
	"github.com/hallazzang/syso/pkg/rsrc"
	"github.com/pkg/errors"
)

type Icon struct {
	Name string
	ID   int
	Path string
}

func (i *Icon) Validate() error {
	if i.Path == "" {
		return errors.New("no file path given")
	} else if i.ID == 0 && i.Name == "" {
		return errors.New("neither id nor name given")
	} else if i.ID != 0 && i.Name != "" {
		return errors.New("id and name cannot be set together")
	} else if i.ID < 0 {
		return errors.Errorf("invalid id: %d", i.ID)
	}
	return nil
}

// Config is a syso config data.
type Config struct {
	Icons []*Icon
}

// ParseConfig reads JSON-formatted syso config from r and returns Config object.
func ParseConfig(r io.Reader) (*Config, error) {
	var c Config
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return nil, errors.Wrap(err, "failed to decode JSON")
	}
	for i, icon := range c.Icons {
		if icon.Name == "" && icon.ID == 0 {
			return nil, errors.Errorf("icon #%d's name and id are both empty", i)
		} else if icon.Path == "" {
			return nil, errors.Errorf("icon #%d's path is empty", i)
		} else if icon.ID < 0 {
			return nil, errors.Errorf("icon #%d: bad id; %d", i, icon.ID)
		}
		for j, icon2 := range c.Icons[:i] {
			if icon.ID != 0 && icon2.ID != 0 && icon2.ID == icon.ID {
				return nil, errors.Errorf("icon #%d's id and icon #%d's id are same", i, j)
			} else if icon.Name != "" && icon2.Name != "" && icon2.Name == icon.Name {
				return nil, errors.Errorf("icon #%d's name and icon #%d's name are same", i, j)
			}
		}
	}
	return &c, nil
}

// EmbedIcon embeds icon into c.
func EmbedIcon(c *coff.File, icon *Icon) error {
	if err := icon.Validate(); err != nil {
		return errors.Wrap(err, "invalid icon")
	}
	s, err := c.Section(".rsrc")
	if err != nil {
		if err == coff.ErrSectionNotFound {
			s = rsrc.New()
			if err := c.AddSection(s); err != nil {
				return errors.New("failed to add new .rsrc section")
			}
		} else {
			return errors.Wrap(err, "failed to get .rsrc section")
		}
	}
	r, ok := s.(*rsrc.Section)
	if !ok {
		return errors.New(".rsrc section is not a valid rsrc section")
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
		return errors.Wrap(err, "failed to add icon group")
	}
	return nil
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
