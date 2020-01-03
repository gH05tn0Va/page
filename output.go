// Functions handling output
package page

import "log"

/*
OutMap{

	"Tag-1": OutList{

		MultiOut{

			"value1-1", // ID 0
			"value1-2", // ID 1
			...

		}, // Child 1

		MultiOut{

			"value2-1", // ID 0
			"value2-2", // ID 1
			...

		}, // Child 2

		...

	},

	"Tag-2": OutList{...},

	...

}
*/

type (
	MultiOut []string
	OutList  []MultiOut
	OutMap   map[string]OutList
	TaskMap  map[string][]int

	SingleOutList []string
	SingleOutMap  map[string]SingleOutList
)

var (
	Out   OutMap
	Tasks TaskMap
)

// map[string][]string -> []string

func (o SingleOutMap) List() SingleOutList {
	var out SingleOutList
	for _, v := range o {
		// var v []string
		out = append(out, v...)
	}
	return out
}

// map[string][]MultiOut -> []MultiOut

func (o OutMap) ListAll() OutList {
	var out OutList
	for _, v := range o {
		// var v PagingOutputValue
		out = append(out, v...)
	}
	return out
}

func (o OutMap) List(pageTag string, task int) OutList {
	Tags.WaitFor(pageTag)
	var out OutList
	for _, url := range Tags.GetUrl(pageTag) {
		out = append(out, o[url].Task(url, task)...)
	}
	return out
}

// map[string][]MultiOut -> MultiOut

func (o OutMap) First(pageTag string) MultiOut {
	Tags.WaitFor(pageTag)
	for _, url := range Tags.GetUrl(pageTag) {
		return o[url].First()
	}
	return nil
}

func (o OutMap) FirstOfTask(pageTag string, task int) MultiOut {
	Tags.WaitFor(pageTag)
	for _, url := range Tags.GetUrl(pageTag) {
		return o[url].Task(url, task).First()
	}
	return nil
}

// map[string]OutList -> map[string]SingleOutList

func (o OutMap) Select(tag string, task int) map[string]SingleOutList {
	Tags.WaitFor(tag)
	out := make(map[string]SingleOutList)
	for k, v := range o {
		// var v OutList
		out[k] = v.Task(k, task).Get(tag)
	}
	return out
}

func (o OutMap) SelectId(selectorId int) map[string]SingleOutList {
	out := make(map[string]SingleOutList)
	for k, v := range o {
		// var v OutList
		out[k] = v.GetId(selectorId)
	}
	return out
}

// []MultiOut -> []MultiOut

func (o OutList) Task(url string, task int) OutList {
	begin := Tasks[url][task]
	end := Tasks[url][task+1]
	if len(o) < end {
		log.Printf("[ERROR] Task %d (begin:%d, end:%d) of %s out of range: %d",
			task, begin, end, url, len(o))
		return nil
	}
	return o[begin:end]
}

// OutList -> SingleOutList

func (o OutList) Get(tag string) SingleOutList {
	return o.GetId(Tags.GetId(tag))
}

func (o OutList) GetId(selectorId int) SingleOutList {
	var out SingleOutList
	for _, v := range o {
		vv := v.GetId(selectorId)
		if vv != "" {
			out = append(out, vv)
		}
	}
	return out
}

// []MultiOut -> MultiOut

func (o OutList) First() MultiOut {
	if len(o) > 0 {
		return o[0]
	}
	return nil
}

// []string -> string

func (o MultiOut) Get(tag string) string {
	return o.GetId(Tags.GetId(tag))
}

func (o MultiOut) GetOr(tag string, defaultStr string) string {
	return o.GetIdOr(Tags.GetId(tag), defaultStr)
}

func (o MultiOut) GetId(selectorId int) string {
	return o.GetIdOr(selectorId, "")
}

func (o MultiOut) GetIdOr(selectorId int, defaultStr string) string {
	if len(o) > selectorId {
		return o[selectorId]
	}
	return defaultStr
}

// []string -> string

func (o SingleOutList) First() string {
	return o.FirstOr("")
}

func (o SingleOutList) FirstOr(defaultStr string) string {
	if len(o) > 0 {
		return o[0]
	}
	return defaultStr
}

func init() {
	Out = make(OutMap)
	Tasks = make(TaskMap)
}
