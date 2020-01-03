// Functions handling output
package page

import "log"

/*
Output PagingOutput{

	"Tag-1": OutputListWithTag{

		OutputWithTag{

			"value1-1", // TaskN 0
			"value1-2", // TaskN 1
			...

		}, // Child 1

		OutputWithTag{

			"value2-1", // TaskN 0
			"value2-2", // TaskN 1
			...

		}, // Child 2

		...

	},

	"Tag-2": OutputListWithTag{...},

	...

}
*/

type (
	OutputWithTag     []string
	OutputListWithTag []OutputWithTag
	PagingOutput      map[string]OutputListWithTag

	OutputList        []string
	OutputListMap     map[string]OutputList
)

var Out PagingOutput

func (o PagingOutput) GetByTag(tag string) OutputListMap {
	return o.SelectBy(tag)
}

func (o PagingOutput) Page(pageTag string) OutputListWithTag {
	return o.ListOf(pageTag)
}

// map[string][]string -> []string

func (o OutputListMap) List() OutputList {
	var out OutputList
	for _, v := range o {
		// var v []string
		out = append(out, v...)
	}
	return out
}

// map[string][]OutputWithTag -> []OutputWithTag

func (o PagingOutput) List() OutputListWithTag {
	var out OutputListWithTag
	for _, v := range o {
		// var v PagingOutputValue
		out = append(out, v...)
	}
	return out
}

func (o PagingOutput) ListOf(pageTag string) OutputListWithTag {
	var out OutputListWithTag
	for _, url := range Tags.Url[pageTag] {
		out = append(out, o[url]...)
	}
	return out
}

// map[string][]OutputWithTag -> OutputWithTag

func (o PagingOutput) FirstOf(pageTag string) OutputWithTag {
	for _, url := range Tags.Url[pageTag] {
		return o[url].First()
	}
	log.Printf("No pageTag '%s'", pageTag)
	return nil
}

// map[string]OutputListWithTag -> map[string]OutputList

func (o PagingOutput) SelectBy(tag string) map[string]OutputList {
	i, ok := Tags.StrId[tag]
	if !ok {
		return nil
	}
	return o.SelectById(i)
}

func (o PagingOutput) SelectById(selectorId int) map[string]OutputList {
	out := make(map[string]OutputList)
	for k, v := range o {
		// var v OutputListWithTag
		out[k] = v.Get(selectorId)
	}
	return out
}

// OutputListWithTag -> OutputList

func (o OutputListWithTag) GetBy(tag string) OutputList {
	i, ok := Tags.StrId[tag]
	if !ok {
		return nil
	}
	return o.Get(i)
}

func (o OutputListWithTag) Get(selectorId int) OutputList {
	var out OutputList
	for _, v := range o {
		vv := v.Get(selectorId)
		if vv != "" {
			out = append(out, vv)
		}
	}
	return out
}

// []OutputWithTag -> OutputWithTag

func (o OutputListWithTag) First() OutputWithTag {
	if len(o) > 0 {
		return o[0]
	}
	return nil
}

// []string -> string

func (o OutputWithTag) GetBy(tag string) string {
	i, ok := Tags.StrId[tag]
	if !ok {
		return ""
	}
	return o.Get(i)
}

func (o OutputWithTag) GetByOr(tag string, defaultStr string) string {
	i, ok := Tags.StrId[tag]
	if !ok {
		return defaultStr
	}
	return o.GetOr(i, defaultStr)
}

func (o OutputWithTag) Get(selectorId int) string {
	return o.GetOr(selectorId, "")
}

func (o OutputWithTag) GetOr(selectorId int, defaultStr string) string {
	if len(o) > selectorId {
		return o[selectorId]
	}
	return defaultStr
}

// []string -> string

func (o OutputList) First() string {
	return o.FirstOr("")
}

func (o OutputList) FirstOr(defaultStr string) string {
	if len(o) > 0 {
		return o[0]
	}
	return defaultStr
}

func init() {
	Out = make(PagingOutput)
}
