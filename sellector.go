// Functions handling html selections
package page

import (
	"github.com/PuerkitoBio/goquery"
)

type (
	SelectorTask struct {
		Name string
		Task []struct {
			SubSel []SubSel
			Func   SelTask
		}
	}

	SelectorJob struct {
		PagingJob
		CurrentSel SelectorTask
	}

	Sel     = *goquery.Selection
	SelTask func(Sel) string
	SubSel  func(Sel) Sel
)

func (s *Urls) Selector(sel string) *SelectorJob {
	sj := new(SelectorJob)
	sj.PagingJob = *New()
	sj.Add(s)
	sj.Selector(sel)
	return sj
}

func (pj *PagingJob) Selector(sel string) *SelectorJob {
	sj := new(SelectorJob)
	sj.PagingJob = *pj
	sj.Selector(sel)
	return sj
}

func (sj *SelectorJob) Selector(sel string) *SelectorJob {
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
	sj.Tasks[len(sj.Tasks)-1] = func(doc Doc) (res []OutputWithTag) {
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

func (sj *SelectorJob) First() *SelectorJob {
	return sj.AddSubTask(
		func(s Sel) Sel {
			return s.First()
		})
}

func (sj *SelectorJob) Contents() *SelectorJob {
	return sj.AddSubTask(
		func(s Sel) Sel {
			return s.Contents()
		})
}

func (sj *SelectorJob) Children() *SelectorJob {
	return sj.AddSubTask(
		func(s Sel) Sel {
			return s.Children()
		})
}

func (sj *SelectorJob) Parents() *SelectorJob {
	return sj.AddSubTask(
		func(s Sel) Sel {
			return s.Parents()
		})
}

func (sj *SelectorJob) ChildrenFiltered(str string) *SelectorJob {
	return sj.AddSubTask(
		func(s Sel) Sel {
			return s.ChildrenFiltered(str)
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

func (sj *SelectorJob) AttrOr(str, defaultStr string) *SelectorJob {
	return sj.AddSelectorTask(
		func(s Sel) string {
			return s.AttrOr(str, defaultStr)
		})
}