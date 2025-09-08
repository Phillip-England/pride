#!/bin/bash

rm -r ./tmp/test-site
go build -o pride main.go
./pride new site test-site
mv pride ./test-site
mv test-site ./tmp
cd ./tmp/test-site
sleep 2
./pride new content ./content/about.md
./pride new content ./content/contact.md
./pride new content ./content/posts/1.md
./pride new content ./content/posts/2.md
./pride new content ./content/posts/3.md
./pride new content ./content/docs/1.md
./pride new content ./content/docs/2.md
./pride new content ./content/docs/3.md
./pride new content ./content/deep/dark/deep_file.md
./pride emit nav ./content ./templates/nav.html
./pride emit nav ./content/posts ./templates/nav-posts.html
./pride emit nav ./content/docs ./templates/nav-docs.html
# ./pride serve 8080
./pride build ./out
