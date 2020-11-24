---
Title: "Blogtober Day 9: Architecture"
Date: October 9, 2020
---
I spent some more time studying up on WAILS, Svelte, and some golang antipatterns to avoid.

It dawned on me sometime early afternoon that the way I've been thinking about WAILS is flawed.  Just because I _can_ serve a WAILS binary as website doesn't mean I _should_.  WAILS is, first and foremost, a desktop app compiler.  Trying to force it to work as a website just so I don't have less code to write isn't going to work out.

So, I've decided to make my codebase adhere to Mat Ryer's pattern.  I will still use Svelte as a frontend for the actual app, and I've figured out how to serve it.  Since it's an SPA, I can scope the frontend logic to the root directory of the site, and the only other endpoints I need define would be API.

