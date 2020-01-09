# page

Web crawler

## Sample Code

```
func f(i int) {
	format := "http://www.h28o.com/cn/vl_newrelease.php?&mode=&page=%d"
	page.FromRange(format, 1, i, 1).As("films").
		Select(".videos .video").
		Find(".id").Text().As("id").
		Children().Attr("href").As("href").
		Children().Attr("title").As("title").
		Run()

	commentsFinder := page.New().
		Select(".page.last").
		Attr("href").As("last").
		Select(".comment").
		Text().As("comment").
		Run()

	for _, t := range page.Tag("href").List() {
		go func(pgTag string) {
			url := "http://www.h28o.com/cn/videocomments.php?mode=2&v=" + pgTag[5:]
			page.From(url).As(pgTag).
				AddToWorker(commentsFinder)
			lastPgStr := page.From(url).As(pgTag).
				AddToWorker(commentsFinder).
				Out().WaitFirst(pgTag).Get("last")
			last, _ := strconv.Atoi(
				regexp.MustCompile(`\d*$`).FindString(lastPgStr))
			page.FromRange(url+"&page=%d", 2, last, 1).As(pgTag).
				AddToWorker(commentsFinder)
		}(t)
	}

	for _, film := range page.PageTag("films", 0) {
		var down string
		var img []string
		pgTag := film.Get("href")
		for _, c := range page.PageTag(pgTag,1).Get("comment") {
			matches := regexp.MustCompile(
				`\[url=([^\]]+)`).FindAllStringSubmatch(c, -1)
			for _, match := range matches {
				if strings.Contains(match[1], "yimuhe") {
					down = match[1]
				}
				if len(img) < 10 && strings.Contains(match[1], "jpg") {
					img = append(img, match[1])
				}
			}
		}

		fmt.Println(film.Get("title"))
		fmt.Println("Download", down)
		fmt.Println("Image", img)
	}
}
```

```
func main() {
	ctxFinder := page.New().
		Selector("[name=pnum]").
		Attr("value").

		Selector(".i").
		Find(".g").Text().
		Contents().First().Text().
		Run()

	url := "http://m.tieba.com/m?kz=6413784508"
	Add(url, ctxFinder)
	ctxFinder.Wait()

	for _, v := range page.PageTag(url, 1) {
		fmt.Println("楼主", v[0])
		fmt.Println(v[1])
	}
}

func Add(url string, w *page.Worker) {
	pg, _ := strconv.Atoi(
		page.OnOne(url).AddToWorker(w).
			Out().First(url).GetIdOr(0, "1"))
	u := strings.ReplaceAll(url, "%", "%%") + "&pn=%d"

	page.OnRange(u, 30, 30*(pg-1), 30).PageTag(url).
		AddToWorker(w)
}
```
