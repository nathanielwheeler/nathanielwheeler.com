---
Title: "Website refactored!"
Date: "<POST DATE>"
---

I have known for some time that my little website has been in sore need of a refactor.  The MVC pattern that I used to initially structure this project has proven to be overly cumbersome, and I have been yearning for a solution more idiomatic to Go.

Now, [back in October](#) I posted some thoughts concerning Mat Ryer's [HTTP Web Service](#) patterns, which I have adopted.  However, this required almost a complete rewrite of my services, which is why I haven't posted in quite a long time.

During this rewrite, I did my best to learn new concepts that I could incorporate.  To this end, the lectures from GopherCon 2020 were an invaluable resource, and I did my best to incorporate my favorite tips into this refactor.

## No More ORM
In his talk [_A Journey to Postgres Productivity with Go_](#), Johan Brandhorst-Satzkorn discussed ways to work with PostgreSQL, including his preferred packages.

Previously, I used [GORM](#) for my database interactions.  With the lessons I learned from this talk, I was able to cut the ORM out and make more efficient queries directly.

## Better Logging
I have also decided to use a more robust logger than what the standard library provides.  I settled on [logrus](#) since it is widely used, well-maintained, and is directly supported by the [pgx](#) driver that I use for Postgres.