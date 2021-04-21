![Selection_097](https://user-images.githubusercontent.com/4579573/115503925-77398b80-a277-11eb-9dbf-704d8c233b8e.png)
# go-tracing

go-tracing is a lightweight, easy-to-use request tracing library. You can use it to see what parts of your system take a long time, and which parts are quicker.

It consists of 2 parts:

- The `tracing` package at `github.com/jamesrr39/go-tracing`. You call this part in your application. This then writes the traces to a file.
- The `tracingviz/cmd/tracingviz-main.go` application. This reads the file, and converts it into an html file. This html file can then be opened and you can see all of your traces. This file is designed to be portable, so you can copy it between computers.

Tracing using the `tracing` package relies heavily on the context, so you must have added the tracer and that specific trace to the context. The provided `Middleware` function will do that for you, if you want to use it for an HTTP server. If you want to use it for something, you can look at the `Middleware` function and adapt it for what you want.

To run the example:

- `go run example/example_http_server.go` (this will print out the location of the newly-creating tracer file)
- `go run tracingviz/cmd/tracingviz-main.go <path-to-tracer-file> <desired-html-output-file-path>`
- open `<desired-html-output-file-path>` with a modern web browser.

What you should see:
![Selection_097](https://user-images.githubusercontent.com/4579573/115503945-7b65a900-a277-11eb-8698-97edf3c74cfc.png)
