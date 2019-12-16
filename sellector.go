package page

import (
	"github.com/PuerkitoBio/goquery"
)

type SelectorCallBack func(*goquery.Selection) interface{}

type SelectorJob struct {
	Job
	Selector    string
	SelCallBack SelectorCallBack
}

func (s *Pages) Selector(sel string) *SelectorJob {
	j := new(SelectorJob)
	j.OutFormat = SubStringOut
	j.Set = s
	j.Selector = sel
	return j
}

func (j *SelectorJob) Do(f SelectorCallBack) *SelectorJob {
	j.SelCallBack = f
	j.CallBack = func(doc *goquery.Document) interface{} {
		var res []interface{}
		doc.Find(j.Selector).Each(func(i int, s *goquery.Selection) {
			out := j.SelCallBack(s)
			if out != nil {
				res = append(res, out)
			}
		})
		return res
	}
	j.setCallback()
	return j
}

func (j *SelectorJob) SetOutput(out interface{}) *SelectorJob {
	switch out.(type) {
	case []int:
		j.OutFormat = SubIntOut
		j.Out.IntSlice = out.([]int)
	case []string:
		j.OutFormat = SubStringOut
		j.Out.Slice = out.([]string)
	case map[string][]int:
		j.OutFormat = SubIntMapOut
		j.Out.IntMapSlice = out.(map[string][]int)
	case map[string][]string:
		j.OutFormat = SubStringMapOut
		j.Out.MapSlice = out.(map[string][]string)
	case []interface{}:
		j.OutFormat = SubInterfaceOut
		j.Out.CustomSlice = out.([]interface{})
	case map[string][]interface{}:
		j.OutFormat = SubInterfaceMapOut
		j.Out.CustomMapSlice = out.(map[string][]interface{})
	case nil:
		break
	default:
		return nil
	}
	return j
}

func (j *SelectorJob) Text() []string {
	j.OutFormat = SubStringOut
	j.Do(
		func(s *goquery.Selection) interface{} {
			return s.Text()
		})
	return j.Run().Wait().Job.Out.Slice
}
