Markdown Go: Literate Programming for Go
========================================

Get, install and run test:

	go get -u github.com/ungerik/mdgo
	cd $GOPATH/src/github.com/ungerik/mdgo
	go install && mdgo
	// mdgo called in this directory creates test.go from test.go.md

mdgo creates .go files from .go.md files by commenting out markdown
and using indented code blocks as go source.

Arguments passed to mdgo will be parsed directly if they are files
or traversed recursively if they are directories.
Every .md.go file in a traversed directory will be parsed.

See [test.go.md](https://github.com/ungerik/mdgo/blob/master/test.go.md)
