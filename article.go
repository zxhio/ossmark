package main

const articleContentTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Title}}</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/atom-one-light.css">
	<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/x86asm.min.js"></script>
	<script>hljs.highlightAll();</script>
	<style>
		body {
			font-family: Arial, sans-serif;
			max-width: 900px;
			margin: 0 auto;
			padding: 20px;
		}

		pre {
			border: 0.1px solid #ddd;
		}

		a {
			color: #61afef;
			text-decoration: none;
		}

		a:hover {
			text-decoration: underline;
		}
	</style>
</head>

<body>
	<h1>{{.Title}}</h1>

	{{.Content}}

</body>

</html>
`

type articleContentBody struct {
	Title   string
	Content string
}
