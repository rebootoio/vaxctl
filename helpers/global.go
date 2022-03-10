package helpers

import (
	"bytes"
	"errors"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v3"
)

func CallFunc(function interface{}, params ...interface{}) (result interface{}, err interface{}) {
	f := reflect.ValueOf(function)
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is out of index.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	var res []reflect.Value
	res = f.Call(in)
	result = res[0].Interface()
	err = res[1].Interface()
	return
}

func ReadFileToJSON(filename string) ([]byte, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data, err := ToJSON(buf)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func EncodeToYaml(data interface{}) ([]byte, error) {
	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	e.SetIndent(2)

	if err := e.Encode(&data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
