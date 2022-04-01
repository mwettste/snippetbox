# Learnings

## Windows Firewall Popups
To avoid the issue of the Windows Firewall having a heart attack every time I start an HTTP listener, it is best to start it with `http.ListenAndServe("localhost:4000", r)`, instead of just setting the port. By doing it on localhost, the firewall won't complain.

## Fixed Path & Subtree Patterns
Go's servemux supports two URL patterns: fixed paths and subtree patterns.
* fixed paths: no trailing backslash, e.g. /foo/bar - this only matches the exact path /foo/bar
* subtree patterns: have a trailing backslash, e.g. /static/ - this matches all paths starting with /static/

# Do not use http.HandleFunc
Although a little bit shorter, it is not recommended to use this in production apps. The reason being that it instantiates a DefaultServeMux in the background, which is a global variable. This means that any 3rd party package that my code uses, could access this ServeMux instance as well and change the registered routes. This could lead to a broken application or in the worst case, it could be used to redirect to malicious websites.

# Go project template
Some well known go project directory templates:
* https://github.com/thockin/go-build-template
* https://peter.bourgon.org/go-best-practices-2016/#repository-structure

# Making logging etc. available in handlers
To make custom loggers available in other code besides the main function, it is idiomatic to define an application struct and add the loggers there. All the methods that need access to these loggers (or configuration etc.) need then to be written as methods against that struct like so `func (app *application) home(...)`.

# Formatting Date & Time
Golang is amazingly flexible with formatting date and time. Below some examples which result in the expected output.

```
func nicerDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
```

```
func nicerDate(t time.Time) string {
	return t.Format("02 January 2006 at 15:04")
}
```