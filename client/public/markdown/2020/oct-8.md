---
Title: "Blogtober Day 8: Midterms and Lectures"
Date: October 8, 2020
---
My classes this semester haven't been really engaging.  They're all core curriculum classes, and none of them really pertain to my field.  Next semester though, I'll be able to take as much Computer Science as I can handle, and I can hardly wait.

Well today, I got one step closer to that desire.  I have finished my midterms!

Afterwards, I didn't really want to do any coding, but I still wanted to do something.  So, I binged a bunch of different lectures.  Y'know, like normal people do.

### Programming Structure Lectures: is there a better way to unwind?

The most interesting in a theoretical sense was Kevlin Henney's [The Forgotten Art of Structured Programming](https://www.youtube.com/watch?v=SFv8Wm2HdNM).  Kevlin runs through a history of programming patterns and paradigms and concludes that a lot of insight was lost in the quest for strict object-oriented programming.  

I vaguely remember watching this lecture at sometime earlier this year, but I definitely got a lot more out of it this time around.  I hope that means I'm getting better at this code stuff.  One can only hope.

Another video I liked was Brian Will's [Object-Oriented Programming is Bad](https://www.youtube.com/watch?v=QM1iUe6IofM).  As you might imagine from the title, Brian Will is much more dismissive concerning OOP than Henney was, advocating for what he calls Procedural Programming, not to be confused with Functional Programming.

His strongest argument against OOP was that, in many cases, the implementation of OOP patterns does the opposite of the desired effect.  His claim was that, where patterns should make code easier to read and maintain, OOP makes it more difficult.  Brian doubled down in his subsequent video [Object-Oriented Programming is Embarrassing: 4 Short Examples](https://www.youtube.com/watch?v=IRTfhkiAqPw), where one of those four examples was taken from Uncle Bob, generally agreed to be an authority on OOP.

The last video I watched was Mat Ryer's lecture at GopherCon 2019: [How I Write HTTP Web Services after Eight Years](https://www.youtube.com/watch?v=rWBSMsLG8po&t=1668s).  Of all the lectures, this one was the most applicable to the real-world.  In it, Mat Ryer lays out some simple patterns he uses to organize his HTTP service code in a clean, maintainable, and testable way.

The thing that struck me the most about his pattern was how flat it was: one package.  Everything centered around a single `server struct`.  Almost all business was conducted by handlers, which utilized anonymous functions that returned `http.HandlerFunc`.  Middleware?  Just go functions that took in these handlers.  Both handlers and middleware were attached to the `server struct` as methods.

What did all this make?  An HTTP codebase that was absurdly easy to navigate and understand.  What more could a developer ask for?

### Moving forward

These were cool lectures, but how do I want to apply these concepts in my own code?  For sure, I want to write flatter code.  No more forcing Go types to behave like C# classes.  With Mat Ryer's pattern on using a single server struct to hold the router, config, logger, and other variables, I can avoid having to endlessly nest the same objects in structs within structs.

With the Theme System project, I did play around a little bit with Mat Ryer's structure for a [WAILS](https://wails.app) app.  Unfortunately, this seems at first pass to be incompatible with the method binding that WAILS incorporates.  In order to have the frontend access the api handlers, I would have to export them, which wouldn't be very secure.  

I will probably have to forego the method binding that WAILS offers and make API calls directly from the frontend.  But then, is WAILS worth implementing?  _Can_ it be implemented without binds?  It's hard to discount the allure of that small cross-platform binary it compiles.  I'll have to explore more on this tomorrow.

### Inktober Day 8

It's the gopher again!  This time, wearing a spacesuit.  Something about it looks familiar...

<img class="card-img-top" src="/images/posts/inktober20-08.jpeg" alt="A gopher wearing a spacesuit from the hit Indie game Among Us."><br>
