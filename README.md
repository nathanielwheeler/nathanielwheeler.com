# README

This is the code for my personal website, designed to be a professional summary of my development skills as well as a blog for both professional and personal topics.

The website, scripted in Go, makes use of Go Templates and Bootstrap for the front-end, and gorilla packages for routing, CSRF protection, and feed generation.

For the database, I use PostgreSQL.  I used to use GORM as an ORM, but I'm currently switching to more direct control so I can make better queries and so I can make use all of Postgres's features.

The project organization is partially inspired by [Mat Ryer's lecture on HTTP Web Services](https://youtu.be/rWBSMsLG8po).  I like to use his handling patterns because of its readability, and I find it easier to maintain than the MVC pattern I used previously.