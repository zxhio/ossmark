package article

import (
	"html/template"
)

const articleListTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>文章列表</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 900px;
            margin: 20px auto;
            padding: 20px;
        }
        .date-header {
            font-weight: bold;
            margin-bottom: 10px;
        }
        .file-list {
            list-style-type: square;
            margin-left: 20px;
        }
        a {
            color: #007bff;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>

<h1>文章列表</h1>

{{- range $index, $month := .MonthArticles }}
<div>
    <div class="date-header">{{$month.Month}}</div>
	<ul class="file-list">
	{{- range $index, $article := .Articles }}
		<li><a href="{{$article.Link}}">{{$article.Name}}</a> <em style="font-size: 12px;"> -- {{$article.LastModify}}</em> </li>
	{{- end }}
	</ul>
</div>
{{- end }}

</body>
</html>
`

var ArticleListTmpl *template.Template

func init() {
	var err error
	ArticleListTmpl, err = template.New("article_list").Parse(articleListTemplate)
	if err != nil {
		panic(err)
	}
}

type article struct {
	Link       string
	Name       string
	LastModify string
}

type monthArticleList struct {
	Month    string
	Articles []article
}

type articleListBody struct {
	MonthArticles []monthArticleList
}

type MonthArticleSlice []monthArticleList

func (m MonthArticleSlice) Len() int           { return len(m) }
func (m MonthArticleSlice) Less(i, j int) bool { return m[i].Month < m[j].Month }
func (m MonthArticleSlice) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

type ArticleSlice []article

func (a ArticleSlice) Len() int           { return len(a) }
func (a ArticleSlice) Less(i, j int) bool { return a[i].LastModify > a[j].LastModify }
func (a ArticleSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
