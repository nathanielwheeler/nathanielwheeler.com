---
Title: "GopherCon"
Date: November 11, 2020
---

It's been a while!  It's the first day of GopherCon, and I have to say I'm very impressed with how the virtual convention has gone!  It's nice to connect to people with similar interests to mine and talk about technical stuff.

So first, I should probably talk about my abrupt stop to Blogtober last month.  Well, I really should have known better.  Daily challenges have just _never_ worked out for me.  I goes great in the beginning, but the first time something goes wrong, it comes crashing down.  In fact, that why I'm so interested in developing the Theme System into an app.  This notion is _built into_ the whole idea of the system: to avoid a New Year resolution and focus on a theme and to frequently reflect.

With that out of the way, I can get back on the wagon and get to back to writing.

## Website plans

I have been thinking about doing a complete rewrite of this website.  After thinking about it, namely about how much work it would be, I have decided to just refactor the sections that are bothering me.  Namely, the way that blog posts themselves are handled.

Right now, blog metadata is stored in PostgreSQL, but refers to a local markdown file for its contents.  This is a bit silly, so I want to rewrite my blog model to simply look up local files and get any metadata from the YAML that prefaces these files.  Doing this will greatly simplify the blog controller, which will simply be serving the files, with no API needed for blog creation, deletion, or modification.

Unfortunately, this will likely have to happen next weekend, since this week I'm focused on GopherCon and with my classes at Fort Hays.