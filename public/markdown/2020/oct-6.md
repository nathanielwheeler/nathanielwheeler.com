---
Title: "Blogtober Day 6: Prototyping"
Date: October 6, 2020
---
I spent my afternoon today making a prototype of the pages I wireframed yesterday, [which you can see here.](/prototypes/theme-system)  Note that UX, position, and form elements were at the front of my mind, not the colors.

This prototype was tricky to make!  I always forget how fiddly forms are, and they make up the majority of the elements in this application.

I am definitely going to want a proper javascript framework for this, because I'm going to need to give the users a lot of feedback based on these forms.

Having made this prototype, I now have to consider what kind of models I want to use.  Here's what I've come up with:

```go
type user struct {
	id uint
	name          string
	email         string
	password      hash
	currentTheme  uint
	activePrompts []uint
}

type theme struct {
	id    uint
	theme string
}

type prompt struct {
	id       uint
	question string
	qtype    string
}

type response struct {
	id         uint
	date       time.Time
	questionID uint
	answer     JSON
}

type page struct {
	id      uint
	themeID uint
	date    time.Time
	data    JSON
}
```

Note that prompt has the field `qtype.`  To start with, this will have two possible values: text and rating, and will be used to check if a value for `response.answer` is valid.

Pages are more of an archival model.  They will be immutable, so even if a user changes their prompt's wording, they will still hold a historical view of that day.  Later on, when I add local storage, I can simply export the pages JSON for this purpose.