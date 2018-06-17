package properties

import (
	"bufio"
	"io"
	"os"
	"strings"
)

//Properties is the structure that represents map of properties
type Properties struct {
	Map      map[string]string
	Defaults *Properties
}

//New creates new instance of Properties
func New() *Properties {
	return &Properties{make(map[string]string), nil}
}

//NewDefault creates new instance of Properties with default values
func NewDefault(d *Properties) *Properties {
	return &Properties{make(map[string]string), d}
}

//Put adds key/value into properties
func (p *Properties) Put(key string, value string) {
	p.Map[key] = value
}

//Get returns value
func (p *Properties) Get(key string) (value string, ok bool) {
	value, ok = p.Map[key]
	if !ok && p.Defaults != nil {
		value, ok = p.Defaults.Get(key)
	}
	return
}

//Filter filters properties
func (p *Properties) Filter(f func(string) bool) (prop *Properties) {
	prop = New()
	keys := p.Keys()
	for _, k := range keys {
		if f(k) {
			v, _ := p.Get(k)
			prop.Put(k, v)
		}
	}
	return
}

//FilterHasPrefix returns all that has prefix
func (p *Properties) FilterHasPrefix(prefix string) *Properties {
	return p.Filter(func(v string) bool {
		return strings.HasPrefix(v, prefix)
	})
}

//Remove removes key
func (p *Properties) Remove(key string) {
	delete(p.Map, key)
}

//GetDefault returns value or defaultValue if key doesn't exists
func (p *Properties) GetDefault(key string, defaultValue string) string {
	value, ok := p.Get(key)
	if !ok {
		return defaultValue
	}
	return value
}

//Keys returns slice of all available keys
func (p *Properties) Keys() []string {
	tmp := make(map[string]bool)

	p.collectKeys(&tmp)

	keys := make([]string, 0, len(tmp))
	for k := range tmp {
		keys = append(keys, k)
	}
	return keys
}

func (p *Properties) collectKeys(m *map[string]bool) {
	if p.Defaults != nil {
		p.Defaults.collectKeys(m)
	}
	for k := range p.Map {
		(*m)[k] = true
	}
}

//ReadFrom reads properties from reader
func (p *Properties) ReadFrom(r io.Reader) (n int64, err error) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		b := s.Bytes()
		n += int64(len(b))
		line := string(b)
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if trimmed[0] == '#' {
			continue
		}
		splitted := strings.SplitN(trimmed, "=", 2)
		key := splitted[0]
		value := splitted[1]
		p.Put(key, value)

	}
	err = s.Err()
	return
}

// LoadFrom loads properties form file
func (p *Properties) LoadFrom(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = p.ReadFrom(file)
	return err
}

//WriteTo writes properties to file
func (p *Properties) WriteTo(w io.Writer) (n int64, err error) {
	nn := 0
	for k, v := range p.Map {

		nn, err = w.Write([]byte(k))
		if err != nil {
			return
		}
		n = n + int64(nn)

		nn, err = w.Write([]byte("="))
		if err != nil {
			return
		}
		n = n + int64(nn)

		nn, err = w.Write([]byte(v))
		if err != nil {
			return
		}
		n = n + int64(nn)

		nn, err = w.Write([]byte("\n"))
		if err != nil {
			return
		}
		n = n + int64(nn)

	}
	return
}

//SaveTo saves propertis to file
func (p *Properties) SaveTo(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = p.WriteTo(file)
	return err
}

//ReadFrom reads from reader
func ReadFrom(r io.Reader) (*Properties, error) {
	p := New()
	_, err := p.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return p, err
}

//LoadFrom loads properties form file
func LoadFrom(fileName string) (*Properties, error) {
	prop := New()
	err := prop.LoadFrom(fileName)
	if err != nil {
		return nil, err
	}
	return prop, nil
}
