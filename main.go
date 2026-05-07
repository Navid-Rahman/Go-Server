// Every Go program starts with a package declaration.
// "main" is special — it tells Go this file is an executable program, not a library.
package main

import (
	"fmt"      // "fmt" is used for formatted I/O, like printing text to the browser or console
	"log"      // "log" is used for logging errors (prints to the terminal and can exit the program)
	"net/http" // "net/http" provides everything we need to build an HTTP web server
)

// formHandler handles requests sent to the "/form" route.
// In Go, HTTP handler functions always take two parameters:
//   - http.ResponseWriter (w): used to write the response back to the client (the browser)
//   - *http.Request (r):       contains all the data about the incoming request (URL, method, form data, etc.)
func formHandler(w http.ResponseWriter, r *http.Request) {
	// ParseForm() reads the submitted form data from the request body and makes it available via r.FormValue()
	// If something goes wrong, we write the error to the browser and stop execution with "return"
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	// If parsing succeeded, write a success message to the browser
	fmt.Fprint(w, "POST request successful")

	// r.FormValue() retrieves the value of a specific field from the submitted HTML form
	// These names ("name", "address") must match the "name" attributes in your HTML <input> tags
	name := r.FormValue("name")
	address := r.FormValue("address")

	// Write the retrieved form values back to the browser so the user can see what was submitted
	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

// helloHandler handles requests sent to the "/hello" route.
// Unlike formHandler, this one enforces specific URL and HTTP method rules.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Guard clause: if someone visits a path other than exactly "/hello", return a 404 error.
	// This is important because Go's default router can be quite permissive with URL matching.
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	// Guard clause: this endpoint only accepts GET requests.
	// If the client sends a POST, PUT, DELETE, etc., we reject it with an error.
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	// If both checks pass, respond with a simple "Hello, World!" message
	fmt.Fprintf(w, "Hello, World!")
}

// main() is the entry point of the program — Go always starts execution here.
func main() {
	// http.FileServer serves static files (HTML, CSS, JS, images) from a directory on disk.
	// http.Dir("./static") points it to the "static" folder in your project directory.
	fileServer := http.FileServer(http.Dir("./static"))

	// http.Handle registers the file server to handle all requests to the root path "/".
	// This means visiting http://localhost:8080/ will serve files from the ./static folder.
	http.Handle("/", fileServer)

	// http.HandleFunc registers our custom handler functions for specific routes.
	// Visiting /form will trigger formHandler, and /hello will trigger helloHandler.
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	// Print a message to the terminal so we know the server started successfully
	fmt.Printf("Starting server at port 8080\n")

	// http.ListenAndServe starts the web server on port 8080.
	// The "nil" means "use the default router" (which we've been registering routes on above).
	// If the server fails to start (e.g. port already in use), log.Fatal prints the error and exits.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
