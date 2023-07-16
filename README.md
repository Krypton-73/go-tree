# go-tree
Go implementation of the unix "tree" command to list contents of directories in a tree-like format.

## Install
```bash
cd <PATH_WHERE_YOU_WANT_REPO>
git clone https://github.com/Krypton-73/go-tree.git
cd go-tree
go build main.go
```

## Usage
```bash
./main [flags]
```

## Flags

```bash
-d, --dir           Flag to only list directories
-h, --help          help for ./main
-i, --indent        Prints tree without indentation lines
-J, --json          Prints tree in JSON format
-L, --level int     Max level of tree depth
-f, --path          Flag to show fullpaths
-p, --permission    Flag to show permission modes
-r, --root string   Root path of the tree (default ".")
-t, --time          Flag to sort output by modified time
-X, --xml           Prints tree in XML format
```

