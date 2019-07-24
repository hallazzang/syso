package syso

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/hallazzang/syso/pkg/coff"
	"github.com/hallazzang/syso/pkg/common"
	"github.com/hallazzang/syso/pkg/rsrc"
	"github.com/hallazzang/syso/pkg/versioninfo"
	"github.com/pkg/errors"
)

// VersionInfoResource represents a version info resource.
type VersionInfoResource struct {
	ID           *int
	Name         *string
	Fixed        *VersionInfoFixed
	StringTables []*VersionInfoStringTable
	Translations []*VersionInfoTranslation
}

// Validate returns an error if the resource is invalid.
func (r *VersionInfoResource) Validate() error {
	if r.ID == nil && r.Name == nil {
		return errors.New("resource id or name must be given")
	} else if r.ID != nil && r.Name != nil {
		return errors.New("resource id and name cannot be given at same time")
	} else if r.ID != nil && *r.ID < 1 {
		return errors.Errorf("invalid resource id; %d", *r.ID)
	} else if r.Name != nil && *r.Name == "" {
		return errors.New("resource name cannot be empty")
	}

	if r.Fixed != nil {
		if err := r.Fixed.Validate(); err != nil {
			return err
		}
	}

	for _, st := range r.StringTables {
		if err := st.Validate(); err != nil {
			return err
		}
	}

	for _, t := range r.Translations {
		if err := t.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// VersionInfoFixed holds fixed information that is language and codepage
// independent, like file version or product version.
type VersionInfoFixed struct {
	FileVersion    *string
	ProductVersion *string
	// TODO: add other fields
}

// Validate returns data validation result.
func (f *VersionInfoFixed) Validate() error {
	if f.FileVersion != nil {
		if _, err := common.ParseVersionString(*f.FileVersion); err != nil {
			return errors.Wrap(err, "failed to parse file version string")
		}
	}
	if f.ProductVersion != nil {
		if _, err := common.ParseVersionString(*f.ProductVersion); err != nil {
			return errors.Wrap(err, "failed to parse product version string")
		}
	}
	return nil
}

// VersionInfoStringTable holds string table associated with language and charset
// pair.
// See https://docs.microsoft.com/en-us/windows/win32/menurc/versioninfo-resource#remarks
// for details about language and charset.
type VersionInfoStringTable struct {
	Language *string
	Charset  *string
	Strings  *VersionInfoStrings
}

// Validate returns data validation result.
func (st *VersionInfoStringTable) Validate() error {
	if _, err := st.languageID(); err != nil {
		return errors.Wrap(err, "failed to parse language identifier")
	}
	if _, err := st.charsetID(); err != nil {
		return errors.Wrap(err, "failed to parse charset identifier")
	}
	if st.Strings == nil {
		return errors.New("strings should present")
	}

	return nil
}

func (st *VersionInfoStringTable) languageID() (uint16, error) {
	return parseUint16ID(st.Language, 0x0409) // default English
}

func (st *VersionInfoStringTable) charsetID() (uint16, error) {
	return parseUint16ID(st.Charset, 0x04b0) // default unicode
}

// VersionInfoStrings holds strings which describes the application.
type VersionInfoStrings struct {
	Comments         *string
	CompanyName      *string
	FileDescription  *string
	FileVersion      *string
	InternalName     *string
	LegalCopyright   *string
	LegalTradeMarks  *string
	OriginalFilename *string
	PrivateBuild     *string
	ProductName      *string
	ProductVersion   *string
	SpecialBuild     *string
}

func (res *VersionInfoStrings) fields() [][2]string {
	var result [][2]string
	target := reflect.ValueOf(res).Elem()
	for i := 0; i < target.NumField(); i++ {
		field := target.Type().Field(i)
		value := target.Field(i)
		if !value.IsNil() {
			result = append(result, [2]string{field.Name, value.Elem().String()})
		}
	}
	return result
}

// VersionInfoTranslation holds language-codepage pairs that application
// supports.
type VersionInfoTranslation struct {
	Language *string
	Charset  *string
}

// Validate returns data validation result.
func (t *VersionInfoTranslation) Validate() error {
	if _, err := t.languageID(); err != nil {
		return errors.Wrap(err, "failed to parse language identifier")
	}
	if _, err := t.charsetID(); err != nil {
		return errors.Wrap(err, "failed to parse charset identifier")
	}

	return nil
}

func (t *VersionInfoTranslation) languageID() (uint16, error) {
	if t.Language == nil {
		return 0, errors.New("language identifier must be set")
	}
	return parseUint16ID(t.Language, 0)
}

func (t *VersionInfoTranslation) charsetID() (uint16, error) {
	if t.Charset == nil {
		return 0, errors.New("charset identifier must be set")
	}
	return parseUint16ID(t.Charset, 0)
}

// EmbedVersionInfo embeds a version info resource.
func EmbedVersionInfo(c *coff.File, v *VersionInfoResource) error {
	if err := v.Validate(); err != nil {
		return errors.Wrap(err, "invalid version info")
	}

	r, err := getOrCreateRSRCSection(c)
	if err != nil {
		return errors.Wrap(err, "failed to get or create .rsrc section")
	}

	vi := versioninfo.New()

	if v.Fixed != nil {
		if v.Fixed.FileVersion != nil {
			if err := vi.SetFileVersionString(*v.Fixed.FileVersion); err != nil {
				return errors.Wrap(err, "failed to set file version")
			}
		}
		if v.Fixed.ProductVersion != nil {
			if err := vi.SetProductVersionString(*v.Fixed.ProductVersion); err != nil {
				return errors.Wrap(err, "failed to set product version")
			}
		}
	}

	for _, st := range v.StringTables {
		lang, _ := st.languageID()
		charset, _ := st.charsetID()
		for _, kv := range st.Strings.fields() {
			vi.SetString(lang, charset, kv[0], kv[1])
		}
	}

	for _, t := range v.Translations {
		lang, _ := t.languageID()
		charset, _ := t.charsetID()
		vi.AddTranslation(lang, charset)
	}

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

func parseUint16ID(s *string, defaultID uint16) (uint16, error) {
	if s == nil {
		return defaultID, nil
	}
	id, err := strconv.ParseUint(*s, 16, 16)
	if err != nil {
		return 0, err
	}
	return uint16(id), nil
}
