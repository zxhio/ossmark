package main

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
		<li><a href="{{$article.Path}}">{{$article.Name}} -- {{$article.LastModify}}</a></li>
	{{- end }}
	</ul>
</div>
{{- end }}

</body>
</html>
`

type article struct {
	Path       string
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

type monthArticles []monthArticleList

func (v monthArticles) Len() int           { return len(v) }
func (v monthArticles) Less(i, j int) bool { return v[i].Month < v[j].Month }
func (p monthArticles) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
