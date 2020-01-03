package page

import "log"

type TagInfo struct {
	Url    []string
	StrId  int
	TaskId int
	Job    *PagingJob
}

type TagInfoMap map[string]TagInfo

// Global tags
var Tags TagInfoMap

func Tag(tag string) SingleOutMap {
	Tags.WaitFor(tag)
	return Out.Select(tag, Tags.GetTask(tag))
}

func PageTag(pageTag string, task int) OutList {
	Tags.WaitFor(pageTag)
	return Out.List(pageTag, task)
}

func (t TagInfoMap) WaitFor(pageTag string) {
	pj := t[pageTag].Job
	if pj != nil {
		if pj.Worker == nil {
			log.Fatalf("Job for '%s' is not started", pageTag)
		}
		pj.Worker.Wait()
	}
}

func (t TagInfoMap) GetUrl(pageTag string) []string {
	info, ok := t[pageTag]
	if ok {
		return info.Url
	}

	if DebugWorker {
		log.Printf("'%s' is not a page tag", pageTag)
	}
	return []string{pageTag}
}

func (t TagInfoMap) GetId(tag string) int {
	info, ok := t[tag]
	if ok && info.StrId > 0 {
		return info.StrId - 1
	}

	log.Fatalf("'%s' is not an id tag", tag)
	return -1
}
func (t TagInfoMap) GetTask(tag string) int {
	info, ok := t[tag]
	if ok && info.TaskId > 0 {
		return info.TaskId - 1
	}

	log.Fatalf("'%s' is not an id tag", tag)
	return -1
}

func (t TagInfoMap) AddJob(tag string, pj *PagingJob) {
	info, ok := t[tag]
	if ok {
		info.Job = pj
		t[tag] = info
	}
}

func (s Urls) PageTag(pageTag string) Urls {
	s.Tag = pageTag

	info, ok := Tags[pageTag]
	if ok {
		info.Url = append(info.Url, s.Data...)
		Tags[pageTag] = info
	} else {
		Tags[pageTag] = TagInfo{
			s.Data, -1, -1,nil}
	}
	return s
}

func (sj *SelectorJob) Tag(tag string) *SelectorJob {
	_, ok := Tags[tag]
	if ok {
		log.Fatalf("'%s' Already exsists!", tag)
	}
	Tags[tag] = TagInfo{
		sj.Input,
		len(sj.CurrentSel.Task) - 1,
		len(sj.Tasks),
		&sj.PagingJob,
	}
	return sj
}

func init() {
	Tags = make(TagInfoMap)
}
