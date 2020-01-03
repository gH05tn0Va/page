# page

Web crawler

## Sample Code

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
