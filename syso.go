package syso

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"

	"github.com/hallazzang/syso/pkg/coff"
	"github.com/hallazzang/syso/pkg/common"
	"github.com/hallazzang/syso/pkg/ico"
	"github.com/hallazzang/syso/pkg/rsrc"
	"github.com/hallazzang/syso/pkg/versioninfo"
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

type VersionInfoResource struct {
	ID      *int
	Name    *string
	Fixed   *FixedFileInfoResource
	Strings *StringFileInfoResource
	// Vars     *VarFileInfoResource // TODO: support it
}

type FixedFileInfoResource struct {
	FileVersion    *string
	ProductVersion *string
	// TODO: add other fields
}

type StringFileInfoResource struct {
	Comments         *string
	CompanyName      *string
	FileDescription  *string
	FileVersion      *string
	InternalName     *string
	LegalCopyright   *string
	LegalTradeMarks  *string
	OriginalFileName *string
	PrivateBuild     *string
	ProductName      *string
	ProductVersion   *string
	SpecialBuild     *string
}

func (res *StringFileInfoResource) fields() map[string]string {
	result := make(map[string]string)
	target := reflect.ValueOf(res).Elem()
	for i := 0; i < target.NumField(); i++ {
		field := target.Type().Field(i)
		value := target.Field(i)
		if !value.IsNil() {
			result[field.Name] = value.Elem().String()
		}
	}
	return result
}

type VarFileInfoResource struct {
}

// Config is a syso config data.
type Config struct {
	Icons       []*FileResource
	Manifest    *FileResource
	VersionInfo *VersionInfoResource
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
	// TODO: validate version info resource
	return &c, nil
}

// EmbedIcon embeds an icon into c.
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
		if err := r.AddResourceByID(rsrc.IconResource, img.ID, img); err != nil {
			return errors.Wrapf(err, "failed to add icon image #%d", i)
		}
	}
	if icon.ID != 0 {
		err = r.AddResourceByID(rsrc.IconGroupResource, icon.ID, icons)
	} else {
		err = r.AddResourceByName(rsrc.IconGroupResource, icon.Name, icons)
	}
	if err != nil {
		return errors.Wrap(err, "failed to add icon group resource")
	}
	return nil
}

// EmbedManifest embeds a manifest into c.
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
		err = r.AddResourceByID(rsrc.ManifestResource, manifest.ID, b)
	} else {
		err = r.AddResourceByName(rsrc.ManifestResource, manifest.Name, b)
	}
	if err != nil {
		return errors.Wrap(err, "failed to add manifest resource")
	}
	return nil
}

func EmbedVersionInfo(c *coff.File, v *VersionInfoResource) error {
	r, err := getOrCreateRSRCSection(c)
	if err != nil {
		return errors.Wrap(err, "failed to get or create .rsrc section")
	}
	vi := versioninfo.New()

	if v.Fixed != nil {
		if v.Fixed.FileVersion != nil {
			if err := vi.SetFileVersionString(*v.Fixed.FileVersion); err != nil {
				return errors.Wrap(err, "failed to set file version string")
			}
		}
		if v.Fixed.ProductVersion != nil {
			if err := vi.SetProductVersionString(*v.Fixed.ProductVersion); err != nil {
				return errors.Wrap(err, "failed to set product version string")
			}
		}
	}
	if v.Strings != nil {
		fs := v.Strings.fields()
		for k, v := range fs {
			vi.SetString(0x0409, 0x04b0, k, v)
		}
	}
	vi.AddTranslation(0x0409, 0x04b0)

	// TODO: need more efficient way
	buf := &bytes.Buffer{}
	if _, err := vi.WriteTo(buf); err != nil {
		return errors.Wrap(err, "failed to write version info data")
	}

	b, err := common.NewBlob(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return errors.Wrap(err, "failed to create version info blob")
	}

	if v.ID != nil {
		err = r.AddResourceByID(rsrc.VersionInfoResource, *v.ID, b)
	} else {
		err = r.AddResourceByName(rsrc.VersionInfoResource, *v.Name, b)
	}
	if err != nil {
		return errors.Wrap(err, "failed to add version info resource")
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
