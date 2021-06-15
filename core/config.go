package core

import (
	"bytes"
	"fmt"
	"github.com/creasty/defaults"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
)

// Config json configuration
var Conf Config

type Log struct {
	Path     string `yaml:"path"`
	Filename string `yaml:"filename"`
	Level    string `yaml:"level"`
	Age      struct {
		Max int `yaml:"max"`
	} `yaml:"age"`
	Size struct {
		Max int `yaml:"max"`
	} `yaml:"size"`
	Backup struct {
		Max int `yaml:"max"`
	} `yaml:"backup"`
}

// Configuration json configuration struct
type Config struct {
	Db struct {
		Type     string `yaml:"type"`
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
	} `yaml:"db"`
	HTTPAPI struct {
		// IP   string `yaml:"ip" `
		Port int `yaml:"port"`
	} `yaml:"http"`
	Logger Log `yaml:"logger"`
}

func setDefaults(v reflect.Value) error {
	tmp := reflect.New(v.Type())
	tmp.Elem().Set(v)
	err := SetDefaults(tmp.Interface())
	if err != nil {
		return err
	}
	v.Set(tmp.Elem())
	return nil
}
func SetDefaults(ptr interface{}) error {
	err := defaults.Set(ptr)
	if err != nil {
		return fmt.Errorf("%v: %s", ptr, err.Error())
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		vf := v.Field(i)
		if tf.Type.Kind() == reflect.Slice {
			for j := 0; j < vf.Len(); j++ {
				item := vf.Index(j)
				if item.Kind() != reflect.Struct {
					continue
				}
				err := setDefaults(item)
				if err != nil {
					return err
				}
			}
		}
		if tf.Type.Kind() == reflect.Map {
			for _, k := range vf.MapKeys() {
				item := vf.MapIndex(k)
				if item.Kind() != reflect.Struct {
					continue
				}
				tmp := reflect.New(item.Type())
				tmp.Elem().Set(item)
				err := setDefaults(tmp.Elem())
				vf.SetMapIndex(k, tmp.Elem())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func unmarshalYAML(in []byte, out interface{}) error {
	err := yaml.Unmarshal(in, out)
	if err != nil {
		return err
	}
	err = SetDefaults(out)
	if err != nil {
		return err
	}
	err = validator.Validate(out)
	if err != nil {
		return err
	}
	return nil
}
func parseEnv(data []byte) ([]byte, error) {
	text := string(data)
	envs := os.Environ()
	envMap := make(map[string]string)
	for _, s := range envs {
		t := strings.Split(s, "=")
		envMap[t[0]] = t[1]
	}
	tmpl, err := template.New("template").Option("missingkey=error").Parse(text)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(nil)
	err = tmpl.Execute(buffer, envMap)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func loadYAML(path string, out interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	res, err := parseEnv(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config parse error: %s", err.Error())
		res = data
	}
	return unmarshalYAML(res, out)
}
func loadConfig(cfg interface{}, confPathFile string) error {
	if confPathFile == "" {
		confPathFile = "./service.yml"
	}
	return loadYAML(confPathFile, cfg)
}

// // Loads load json configuration
// func (v *Configuration) Loads(conf string) error {
// 	if f, err := io.ReadFile(conf); err != nil {
// 		return err
// 	} else {
// 		d := []byte(f)
// 		if err := json.Unmarshal(d, v); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
//
// // NewConfig create new configuration
// func NewConfig() *Configuration {
// 	c := &Configuration{}
// 	return c
// }
