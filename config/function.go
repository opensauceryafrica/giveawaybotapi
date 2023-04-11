package config

import (
	"fmt"
	"reflect"

	"github.com/opensaucerers/giveawaybot/typing"
)

// VerifyEnvironment checks if all environment variables are set
func VerifyEnvironment(env typing.Env) error {
	// get the type of argument
	t := reflect.TypeOf(env)
	if t == nil {
		return fmt.Errorf("env is nil")
	}
	// only allow struct type
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("env is not a struct")
	}
	// verify each struct field tag
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// get the field tag value
		tag := field.Tag.Get(envTagName)
		if tag == "" {
			continue
		}
		// check if environment variable is set
		if MustGet(tag, " ") == " " {
			return fmt.Errorf("environment variable %s is not set", tag)
		}
	}
	return nil
}

// AppendEnvironment appends environment variables to constant.Env
func AppendEnvironment(env *typing.Env) {
	// get the type of argument
	t := reflect.TypeOf(*env)
	if t == nil {
		return
	}
	// only allow struct type
	if t.Kind() != reflect.Struct {
		return
	}
	// append each struct field tag
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// get the field tag value
		tag := field.Tag.Get(envTagName)
		if tag == "" {
			continue
		}
		// append environment variable to constant.Env
		reflect.ValueOf(env).Elem().Field(i).SetString(MustGet(tag, ""))
	}
}
