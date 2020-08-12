# README

This is the code for my personal website, designed to be a professional summary of my development skills as well as a blog for both professional and personal topics.

The website, scripted in Go, makes use of Go Templates and Bootstrap for the front-end, gorilla/mux for routing, GORM for database interactions, and PostgreSQL for the database itself.

The project is organized using the Model-View-Controller pattern.  Models have service and validation layers, as well as ORM interfaces.  User middleware allows for authentication and authorization control.