# Pride
Simple static sites.

## About
Pride is a static site generator intended to make building static sites as simple as possible.

## Installation
Pride may be built from [source](https://github.com/phillip-england/pride) with `go` version `1.24.4` or later.
```bash
git clone https://github.com/phillip-england/pride
cd pride
go install
```
Assuming your `$GOPATH` is set correct, you should be able to view the help display using:
```bash
pride help
```

## Initalize a Project
```bash
pride new site my-site
cd my-site
```

## Create New Content
```bash
pride new content ./content/about.md
```

## Generate Navigation
`pride emit nav` enables us to generate navigation based on the structure of the `./content` dir. Using `pride emit nav`, we can point to any dir within the content dir (including the content dir itself) and output an html file containing the navigation.
At the top of the `.html` file, you will see the name of the navigation wrapped in `{{ }}` braces.
```bash
pride emit nav ./content ./templates/nav.html
```

## Serve Locally
To serve on `localhost:8080`:
```bash
pride serve 8080
```
Pride enables hot-reload by default for any file changed within the project.

## Build Static Site
To generate within `./out`:
```bash
pride build ./out
```