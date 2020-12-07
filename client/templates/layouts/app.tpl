{{define "app"}}
<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" href="/assets/main.css">
	<title>NJW blog</title>
</head>

<body>

  {{template "navbar"}}

	{{template "yield" .}}

</body>

</html>
{{end}}

{{define "navbar"}}
<nav class="navbar navbar-expand-md navbar-secondary bg-secondary text-dark">
  <a href="/" class="navbar-brand">
    <img src="images/logoface.svg" alt="Nathan's face" width="100px" class="d-inline-block align-top">Nathaniel Wheeler
  </a>
</nav>
{{end}}