package chirouter_test

import (
	"testing"

	"github.com/bdoerfchen/webcmd/src/common/config"
	"github.com/bdoerfchen/webcmd/src/services/chirouter"
	"github.com/stretchr/testify/assert"
)

func TestParseParameters(t *testing.T) {
	testCases := []struct {
		Name               string
		Route              string
		ExpectedParamNames []string
	}{
		{Name: "none", Route: "/test", ExpectedParamNames: []string{}},
		{Name: "simple", Route: "/hello/{Id}", ExpectedParamNames: []string{"Id"}},
		{Name: "with config", Route: "/bonjour/{id:number}/something", ExpectedParamNames: []string{"id"}},
		{Name: "regex", Route: "/bonjour/{äöaSp:[A-Za-z]}/something", ExpectedParamNames: []string{"äöaSp"}},
		{Name: "broken", Route: "/bonjour/{äöaSp:as/something", ExpectedParamNames: []string{}},
		{Name: "multiple", Route: "/multiple/{one}/something/{two}", ExpectedParamNames: []string{"one", "two"}},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			route := chirouter.OptimizeRoute(config.Route{Route: tc.Route})
			assert.ElementsMatch(t, route.ParamNames, tc.ExpectedParamNames)
		})
	}
}
