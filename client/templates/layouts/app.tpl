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
<nav class="container-flex bg-dark">
  <h1 class="row">
    <a class="col-12 text-center" href="/">Nathaniel Wheeler</a>
  </h1>
  <div class="row">
    <div class="col-12 text-center">
      <a href="#">Portfolio</a> | 
      <a href="#">Contact</a>
    </div>
  </div>
</nav>
{{end}}