package rules

import (
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/wmatex/automux/internal/tmux"
)

func LoadFromConfig() (*Rules, error) {
	var rules Rules
	err := viper.Unmarshal(&rules, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			paneHookFunc(), windowHookFunc(), ruleHookFunc(),
		),
	))

	return &rules, err
}

func ruleHookFunc() mapstructure.DecodeHookFuncType {
	// Wrapped in a function call to add optional input parameters (eg. separator)
	return func(
		f reflect.Type, // data type
		t reflect.Type, // target data type
		data interface{}, // raw data
	) (interface{}, error) {
		// Check if the data type matches the expected one
		if f.Kind() != reflect.Map || f.Key().Kind() != reflect.String {
			return data, nil
		}

		// Check if the target type matches the expected one
		if t.Kind() != reflect.Slice || !t.Elem().Implements(reflect.TypeOf((*RuleCheck)(nil)).Elem()) {
			return data, nil
		}

		// Format/decode/parse the data and return the new value
		m := data.(map[string]interface{})
		var rules []RuleCheck
		for key := range m {
			r, err := ruleCheckFactory(key, m[key].(string))
			if err != nil {
				log.Fatal(err)
			}
			rules = append(rules, r)
		}

		return rules, nil
	}
}

func firstMapKey(m map[string]interface{}) string {
	for k := range m {
		return k
	}

	return ""
}

func windowHookFunc() mapstructure.DecodeHookFuncType {
	// Wrapped in a function call to add optional input parameters (eg. separator)
	return func(
		f reflect.Type, // data type
		t reflect.Type, // target data type
		data interface{}, // raw data
	) (interface{}, error) {
		// Check if the data type matches the expected one
		if f.Kind() == reflect.Map && f.Key().Kind() == reflect.String && t == reflect.TypeOf(tmux.Window{}) {
			w := data.(map[string]interface{})
			name := firstMapKey(w)

			window := w[name].(map[string]interface{})

			panesConfig := window["panes"].([]interface{})
			var panes []*tmux.Pane
			for _, cmd := range panesConfig {
				panes = append(panes, &tmux.Pane{
					Cmd: cmd.(string),
				})
			}

			return tmux.Window{
				Name:  name,
				Panes: panes,
			}, nil

		}

		return data, nil
	}
}

func paneHookFunc() mapstructure.DecodeHookFuncType {
	// Wrapped in a function call to add optional input parameters (eg. separator)
	return func(
		f reflect.Type, // data type
		t reflect.Type, // target data type
		data interface{}, // raw data
	) (interface{}, error) {
		if f.Kind() == reflect.String && t == reflect.TypeOf(tmux.Pane{}) {
			return tmux.Pane{
				Cmd: data.(string),
			}, nil
		}

		return data, nil
	}
}
