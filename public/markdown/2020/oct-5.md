---
Title: "Blogtober Day 5: Wireframes"
Date: October 5, 2020
---
After a relaxing weekend, I'm eager to get get back to working on the Theme System app.  Since I'm also taking classes online, I started out my day studying for this week's midterms.  If I don't start my day studying, I find that I get caught up in programming and fall behind in my courses.

After I finished my review work for the day, I started working on wireframes for my Theme System project.  After a few revisions, I decided that my MVP application would work with 5 pages, bare minimum.

My initial thoughts are that this project would do well as an SPA, and I'm thinking that I will build the frontend with Svelte and package it with the server using WAILS.  That said, it's early days yet, so 

I have requirements and wireframes for each page listed below.

<div class="accordion" id="requirements">
	<div class="card border-primary bg-darker">
		<div class="card-header" id="dashboardHeader">
			<button class="btn btn-link text-left btn-block" data-toggle="collapse" data-target="#dashboard"
				aria-expanded="true" aria-controls="dashboard">Dashboard</button>
		</div>
		<div class="collapse show" id="dashboard" aria-labelledby="dashboardHeader" data-parent="#requirements">
			<div class="card-body">
				<img src="/images/posts/themesystem/template-login.jpeg" alt="" class="card-img-top mx-auto pb-3">
				<h6>Requirements</h6>
				<ul>
					<li>Show the user's yearly or seasonal theme.</li>
					<li>If the day has not been logged, prompt the user with user-defined journal questions.</li>
					<li>If the day has been logged, show the input answers.</li>
					<li>Provide navigation elements to review historical data, like an arrow.</li>
					<li>Provide a nav element to go into themes.</li>
					<li>Show a logout button.</li>
				</ul>
			</div>
		</div>
	</div>
	<div class="card border-primary bg-darker">
		<div class="card-header" id="registerHeader">
			<button class="btn btn-link text-left btn-block collapsed" data-toggle="collapse" data-target="#register"
				aria-expanded="false" aria-controls="register">Register</button>
		</div>
		<div class="collapse" id="register" aria-labelledby="registerHeader" data-parent="#requirements">
			<div class="card-body">
				<img src="/images/posts/themesystem/template-register.jpeg" alt="" class="card-img-top mx-auto pb-3">
				<h6>Requirements</h6>
				<ul>
					<li>Take in name, email, and password to create an account.</li>
					<li>Provide a button to go to login if the user already has an account.</li>
					<li>Redirect user to views to setup their dashboard upon successful registration.</li>
					<li>Provide feedback on invalid names, emails, and passwords.</li>
				</ul>
			</div>
		</div>
	</div>
	<div class="card border-primary bg-darker">
		<div class="card-header" id="loginHeader">
			<button class="btn btn-link text-left btn-block collapsed" data-toggle="collapse" data-target="#login"
				aria-expanded="false" aria-controls="login">Login</button>
		</div>
		<div class="collapse" id="login" aria-labelledby="loginHeader" data-parent="#requirements">
			<div class="card-body">
				<img src="/images/posts/themesystem/template-login.jpeg" alt="" class="card-img-top mx-auto pb-3">
				<h6>Requirements</h6>
				<ul>
					<li>Take in email and password to login.</li>
					<li>Provide a button to go to register if the user doesn't have an account.</li>
					<li>Redirect user to their dashboard upon successful login.</li>
					<li>Provide feedback on invalid login information.</li>
				</ul>
			</div>
		</div>
	</div>
	<div class="card border-primary bg-darker">
		<div class="card-header" id="themesHeader">
			<button class="btn btn-link text-left btn-block collapsed" data-toggle="collapse" data-target="#themes"
				aria-expanded="false" aria-controls="themes">Themes</button>
		</div>
		<div class="collapse" id="themes" aria-labelledby="themesHeader" data-parent="#requirements">
			<div class="card-body">
				<img src="/images/posts/themesystem/template-themes.jpeg" alt="" class="card-img-top mx-auto pb-3">
				<h6>Themes Requirements</h6>
				<ul>
					<li>Show the user has not defined their theme, questions, or prompts, display examples. Otherwise, show those
						values.</li>
					<li>If a user changes their themes, prompt them to confirm.</li>
					<li>Show icons that can add or remove questions or prompts.</li>
					<li>Show a nav element that can return the user their dashboard</li>
				</ul>
			</div>
		</div>
	</div>
	<div class="card border-primary bg-darker">
		<div class="card-header" id="historyHeader">
			<button class="btn btn-link text-left btn-block collapsed" data-toggle="collapse" data-target="#history"
				aria-expanded="false" aria-controls="history">History</button>
		</div>
		<div class="collapse" id="history" aria-labelledby="historyHeader" data-parent="#requirements">
			<div class="card-body">
				<img src="/images/posts/themesystem/template-history.jpeg" alt="" class="card-img-top mx-auto pb-3">
				<h6>History Requirements</h6>
				<ul>
					<li>Show the date of the entry, as well as the archived immutable data for that entry.</li>
					<li>Show arrows that allow the user to go back or forth in time.</li>
				</ul>
			</div>
		</div>
	</div>
</div>
<br>

You may have noticed that I don't have a page for analytics.  I have decided that analytics will be a stretch goal as they are not crucial to a minimum viable product.

I also didn't try to do anything drastically different between mobile and desktop.  I didn't feel like different layouts would have been worth the increased complexity of the frontend.  Of course, I may change my mind later, but this gives me a solid starting point.

Tomorrow, I'm going to build an interactive prototype and begin fleshing out what sort of resources I will be using in the app.  The prototype will show me if my wireframing is as intuitive of a user experience as I think it is.  It will also help me visualize what sort of data I will need to be passing between the server and client.

### Inktober: Day 5

<img class="card-img-top" src="/images/posts/inktober20-05.jpeg" alt="ALT TEXT"><br>