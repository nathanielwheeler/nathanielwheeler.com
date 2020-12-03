{{define "markdown"}}
<section class="markdown">
	{{if ("isMarkdown" .)}}

	{{else}}
		{{template "alert" "Post content failed to load."}}
	{{end}}
</section>
{{end}}