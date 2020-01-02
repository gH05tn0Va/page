// Functions handling output
package page

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

type OutputList []string

var (
	IdMap  map[string]int
	JobMap map[string]*PagingJob
)

func (sj *SelectorJob) Tag(tag string) *SelectorJob {
	IdMap[tag] = len(sj.CurrentSel.Task) - 2
	JobMap[tag] = &sj.PagingJob
	return sj
}

// Global output

func OutputBy(tag string) map[string]OutputList {
	j, ok := JobMap[tag]
	if !ok {
		return nil
	}
	return j.Out().MapBy(tag)
}

func OutputListBy(tag string) OutputList {
	j, ok := JobMap[tag]
	if !ok {
		return nil
	}
	return j.Out().ListBy(tag)
}

func OutputFilterBy(pageTag string) OutputListWithTag {
	j, ok := JobMap[pageTag]
	if !ok {
		return nil
	}
	return j.Out()[pageTag]
}

// map[string][]OutputWithTag -> map[string]OutputWithTag

func (o PagingOutput) MapBy(tag string) map[string]OutputList {
	i, ok := IdMap[tag]
	if !ok {
		return nil
	}
	return o.Map(i)
}

func (o PagingOutput) Map(selectorId int) map[string]OutputList {
	out := make(map[string]OutputList)
	for k, v := range o {
		// var v []OutputWithTag
		for _, vv := range v {
			out[k] = append(out[k], vv[selectorId])
		}
	}
	return out
}

// map[string][]OutputWithTag -> []OutputWithTag

func (o PagingOutput) ListAll() OutputListWithTag {
	var out OutputListWithTag
	for _, v := range o {
		// var v []OutputWithTag
		out = append(out, v...)
	}
	return out
}

func (o PagingOutput) FilterBy(pageTags []string) OutputListWithTag {
	var out OutputListWithTag
	for _, pageTag := range pageTags {
		out = append(out, o[pageTag]...)
	}
	return out
}

// map[string][]OutputWithTag -> OutputWithTag

func (o PagingOutput) FirstOf(pageTag string) OutputWithTag {
	return o[pageTag].First()
}

// map[string]OutputListWithTag -> OutputList

func (o PagingOutput) ListBy(tag string) OutputList {
	i, ok := IdMap[tag]
	if !ok {
		return nil
	}
	return o.List(i)
}

func (o PagingOutput) List(selectorId int) OutputList {
	var out OutputList
	for _, v := range o {
		// var v []OutputWithTag
		out = append(out, v.Get(selectorId)...)
	}
	return out
}

// OutputListWithTag -> OutputList

func (o OutputListWithTag) GetBy(tag string) OutputList {
	i, ok := IdMap[tag]
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
	i, ok := IdMap[tag]
	if !ok {
		return ""
	}
	return o.Get(i)
}

func (o OutputWithTag) GetByOr(tag string, defaultStr string) string {
	i, ok := IdMap[tag]
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
	IdMap = make(map[string]int)
	JobMap = make(map[string]*PagingJob)
}
