package page

type TagInfo struct {
	Url   map[string][]string
	StrId map[string]int
}

// Global tags
var Tags *TagInfo

func (s Urls) PageTag(pageTag string) Urls {
	for _, url := range s {
		Tags.Url[pageTag] = append(Tags.Url[pageTag], url)
	}
	return s
}

func (sj *SelectorJob) Tag(tag string) *SelectorJob {
	Tags.StrId[tag] = len(sj.CurrentSel.Task) - 2
	return sj
}

func init() {
	Tags = new(TagInfo)
	Tags.Url = make(map[string][]string)
	Tags.StrId = make(map[string]int)
	/*
		Out.IdMap = make(map[string]int)
		Out.JobMap = make(map[string]*PagingJob)
		Out.Data = make(map[int]OutputListWithTag)
		Out.Tag = make(map[string]string)
		Out.Id = make(map[string]int)
		Out.TagId = make(map[string][]int)
	*/
}
