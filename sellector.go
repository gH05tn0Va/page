package page

import (
	"github.com/PuerkitoBio/goquery"
	"regexp"
)

type (
	Sel     = *goquery.Selection
	RegTask func(Sel) []string
	SelTask func(Sel) string
	SubSel  func(Sel) Sel
)

type SelectorTask struct {
	Name string
	Task []struct {
		SubSel []SubSel
		Func   SelTask
	}
	RegExpr RegTask
}

type SelectorJob struct {
	Job
	CurrentSel SelectorTask
}

func (s *Urls) AddSelector(sel string) *SelectorJob {
	sj := new(SelectorJob)
	sj.Set = *s
	sj.Output = make(map[string]PageOutput)
	sj.WorkerFunc = func(url string, j *Job) error {
		doc, err := GetPageBody(url)
		if err != nil {
			return err
		}
		var out []PageTaskOutput
		for _, taskFunc := range j.Tasks {
			v := taskFunc(doc)
			if v != nil {
				out = append(out, v)
			}
		}
		j.Output[url] = out
		return nil
	}
	sj.AddSelector(sel)
	return sj
}

func (j *Job) AddSelector(sel string) *SelectorJob {
	sj := new(SelectorJob)
	sj.Job = *j
	sj.AddSelector(sel)
	return sj
}

func (sj *SelectorJob) AddSelector(sel string) *SelectorJob {
	sj.CurrentSel = SelectorTask{
		Name: sel, Task: []struct {
			SubSel []SubSel
			Func   SelTask
		}{{}}}
	sj.AddTask(nil)
	return sj
}

func (sj *SelectorJob) AddSelectorTask(f SelTask) *SelectorJob {
	sel := sj.CurrentSel
	sel.Task[len(sel.Task)-1].Func = f

	// Update the Job.Tasks
	sj.Tasks[len(sj.Tasks)-1] = func(doc Doc) (res [][]string) {
		doc.Find(sel.Name).Each(func(i int, s Sel) {
			var selOut []string
			selSkip := true
			for _, selTask := range sel.Task {
				tmp := s
				for _, sub := range selTask.SubSel {
					tmp = sub(tmp)
				}
				f := selTask.Func
				if f != nil {
					v := f(tmp)
					if v != "" {
						selSkip = false
					}
					selOut = append(selOut, f(tmp))
				}
			}
			if !selSkip {
				res = append(res, selOut)
			}
		})
		return
	}

	sel.Task = append(sel.Task, struct {
		SubSel []SubSel
		Func   SelTask
	}{})

	sj.CurrentSel = sel
	return sj
}

func (sj *SelectorJob) AddSubTask(f SubSel) *SelectorJob {
	sel := sj.CurrentSel
	tsk := sel.Task[len(sel.Task)-1].SubSel

	tsk = append(tsk, f)
	sel.Task[len(sel.Task)-1].SubSel = tsk

	sj.CurrentSel = sel
	return sj
}

func (sj *SelectorJob) Children() *SelectorJob {
	return sj.AddSubTask(
		func(s Sel) Sel {
			return s.Children()
		})
}

func (sj *SelectorJob) Find(str string) *SelectorJob {
	return sj.AddSubTask(
		func(s Sel) Sel {
			return s.Find(str)
		})
}

func (sj *SelectorJob) Text() *SelectorJob {
	return sj.AddSelectorTask(
		func(s Sel) string {
			return s.Text()
		})
}

func (sj *SelectorJob) Attr(str string) *SelectorJob {
	return sj.AddSelectorTask(
		func(s Sel) string {
			out, _ := s.Attr(str)
			return out
		})
}

func (sj *SelectorJob) Match(str string) *SelectorJob {
	return sj.AddSelectorTask(
		func(s Sel) string {
			regexp := regexp.MustCompile(str)
			return regexp.FindString(s.Text())
		})
}

func (sj *SelectorJob) SubMatch(str string) *SelectorJob {
	return sj.AddSelectorTask(
		func(s Sel) string {
			regexp := regexp.MustCompile(str)
			res := regexp.FindStringSubmatch(s.Text())
			if len(res) >= 2 {
				return res[1]
			}
			return ""
		})
}

func (sj *SelectorJob) AddRegularExprTask(f RegTask) *SelectorJob {
	sel := sj.CurrentSel
	sel.RegExpr = f

	// Update the Job.Tasks
	sj.Tasks[len(sj.Tasks)-1] = func(doc Doc) (res [][]string) {
		doc.Find(sel.Name).Each(func(i int, s Sel) {
			var out []string
			regF := sel.RegExpr
			if regF != nil {
				v := regF(s)
				if len(v) > 0 {
					res = append(res, v)
				}
			}

			for _, selTask := range sel.Task {
				tmp := s
				for _, sub := range selTask.SubSel {
					tmp = sub(tmp)
				}
				f := selTask.Func
				if f != nil {
					out = append(out, f(tmp))
				}
			}
			if len(out) > 0 {
				res = append(res, out)
			}
		})
		return
	}

	sel.Task = append(sel.Task, struct {
		SubSel []SubSel
		Func   SelTask
	}{})

	sj.CurrentSel = sel
	return sj
}
