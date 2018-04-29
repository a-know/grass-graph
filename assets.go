package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets823731ccafd08ec013392b546649056853d4385b = "<html>\n    test\n</html>\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"public"}, "/public": []string{"index.html"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1524971911, 1524971911000000000),
		Data:     nil,
	}, "/public": &assets.File{
		Path:     "/public",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1524971380, 1524971380000000000),
		Data:     nil,
	}, "/public/index.html": &assets.File{
		Path:     "/public/index.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1524971394, 1524971394000000000),
		Data:     []byte(_Assets823731ccafd08ec013392b546649056853d4385b),
	}}, "")
