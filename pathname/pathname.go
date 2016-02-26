package pathname

import (
	"os"
	"path"
	"path/filepath"
	"time"
)

type Pathname struct {
	path string
}

func NewPathname(p ...string) *Pathname {
	fullpath := path.Join(p...)
	return &Pathname{fullpath}
}

func Getwd() *Pathname {
	pwd, _ := os.Getwd()
	return NewPathname(pwd)
}

func TempDir() *Pathname {
	return NewPathname(os.TempDir())
}

func FindInPath(name string, path []string) []Pathname {
	results := []Pathname{}

	for _, p := range path {
		dir := NewPathname(p)
		if dir.IsBlank() {
			continue
		}
		for _, cmd := range dir.EntriesMatching(name) {
			if cmd.IsExecutable() {
				results = append(results, *cmd.Abs())
			}
		}
	}

	return results
}

func (p *Pathname) String() string {
	return p.path
}

func (p *Pathname) Dir() *Pathname {
	return NewPathname(path.Dir(p.path))
}

func (p *Pathname) Base() string {
	return path.Base(p.path)
}

func (p *Pathname) Abs() *Pathname {
	abs, err := filepath.Abs(p.path)
	if err == nil {
		return NewPathname(abs)
	} else {
		return p
	}
}

func (p *Pathname) Join(names ...string) *Pathname {
	components := []string{p.path}
	components = append(components, names...)
	return NewPathname(path.Join(components...))
}

func (p *Pathname) IsBlank() bool {
	return p.path == ""
}

func (p *Pathname) IsRoot() bool {
	return p.path == "/"
}

func (p *Pathname) IsExecutable() bool {
	fileInfo, err := os.Stat(p.path)
	return err == nil && (fileInfo.Mode()&0111) != 0
}

func (p *Pathname) IsFile() bool {
	fileInfo, err := os.Stat(p.path)
	return err == nil && !fileInfo.IsDir()
}

func (p *Pathname) ModTime() (time.Time, error) {
	fileInfo, err := os.Stat(p.path)
	if err == nil {
		return fileInfo.ModTime(), nil
	} else {
		epoch, _ := time.Parse("2006-Jan-02", "1970-01-01")
		return epoch, err
	}
}

func (p *Pathname) Exists() bool {
	_, err := os.Stat(p.path)
	return err == nil
}

func (p *Pathname) Equal(other *Pathname) bool {
	return p.String() == other.String()
}

func (p *Pathname) Entries() []Pathname {
	if file, err := os.Open(p.path); err == nil {
		entries, err := file.Readdirnames(0)
		if err == nil {
			results := make([]Pathname, len(entries))
			for i, entry := range entries {
				results[i] = *p.Join(entry)
			}
			return results
		}
	}
	return []Pathname{}
}

func (p *Pathname) EntriesMatching(pattern string) []Pathname {
	entries, _ := filepath.Glob(p.path + "/" + pattern)
	results := make([]Pathname, len(entries))
	for i, entry := range entries {
		results[i] = *NewPathname(entry)
	}
	return results
}

func (p *Pathname) MkdirAll() error {
	return os.MkdirAll(p.path, 0755)
}

func (p *Pathname) Create() (*os.File, error) {
	err := p.Dir().MkdirAll()
	if err != nil {
		return nil, err
	}

	return os.Create(p.path)
}
