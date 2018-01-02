## runw

simple minimalistic watch for go project

## usage
`runw <folder-containing-go-files-to-watch> <go-file-to-run-on-change>` // non recursive

`runw -r <folder-containing-go-files-to-watch> <go-file-to-run-on-change>` // recursive watcher

## example
`runw ./src main.go`

or

`runw blockchain.go main.go` // will run main.go on change to blockchain.go

or

`runw -r ./src ./src/main.go` // for recursive
