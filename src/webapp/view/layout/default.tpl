<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
	<link rel="stylesheet" href="/assets/css/style.css" media="all">
	<title>{{template "title" .}}</title>
</head>
<body>
{{template "content" .}}
</body>
</html>
{{/* 
   * 空のblockを定義しておくと、コンテンツ側で未定義でもエラー回避できる
   */}}
{{define "title"}}{{end}}

