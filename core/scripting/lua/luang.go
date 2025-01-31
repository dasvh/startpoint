package luang

import (
	"fmt"

	"github.com/susiteemu/startpoint/core/model"

	"github.com/rs/zerolog/log"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type result struct {
	Url     string
	Method  string
	Headers map[string]interface{}
	Body    interface{}
	Auth    map[string]interface{}
	Options map[string]interface{}
	Output  string
}

func RunLuaScript(request model.RequestMold, previousResponse *model.Response) (map[string]interface{}, error) {
	L := lua.NewState()
	defer L.Close()

	prevResponseMap := map[string]interface{}{}
	if previousResponse != nil {
		headers := previousResponse.HeadersAsMapString()
		prevResponseMap["headers"] = headers
	}
	prevResponse := luar.New(L, prevResponseMap)
	L.SetGlobal("prevResponse", prevResponse)

	if err := L.DoString(request.Scriptable.Script); err != nil {
		log.Error().Err(err).Msg("Running Lua script resulted to error")
		return nil, err
	}
	lv := L.Get(-1)
	values := map[string]interface{}{}
	if lv.Type() == lua.LTTable {
		res := result{}
		err := gluamapper.NewMapper(gluamapper.Option{
			NameFunc: gluamapper.Id,
		}).Map(lv.(*lua.LTable), &res)
		if err != nil {
			return map[string]interface{}{}, err
		}

		// Convert map of interface{} to map of string
		bodyAsMapInterface, isMapInterface := res.Body.(map[interface{}]interface{})
		if isMapInterface {
			bodyAsMapString := map[string]interface{}{}
			for k, v := range bodyAsMapInterface {
				bodyAsMapString[fmt.Sprintf("%v", k)] = v
			}
			res.Body = bodyAsMapString
		}

		log.Debug().Msgf("Received from Lua: %v", res)
		values["url"] = res.Url
		values["method"] = res.Method
		values["headers"] = res.Headers
		values["body"] = res.Body
		values["auth"] = res.Auth
		values["options"] = res.Options
		values["output"] = res.Output
	} else {
		return map[string]interface{}{}, fmt.Errorf("Expected table, got %v", lv.Type())
	}

	return values, nil
}
