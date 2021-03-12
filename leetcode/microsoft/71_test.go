package microsoft

import (
	"bytes"
	"container/list"
	path2 "path"
	"testing"
)

func simplifyPath(path string) string {
	path2.Clean(path)
	stack := list.New()
	var i int
	var length int
	for i < len(path) {
		c := path[i]
		j := i
		var part string
		switch c {
		case '.':
			for ; i < len(path) && path[i] == '.'; i++ {
			}
			part = path[j:i]
		case '/':
			for ; i < len(path) && path[i] == '/'; i++ {
			}
			part = "/"
		default:
			for ; i < len(path) && path[i] != '/' && path[i] != '.'; i++ {

			}
			part = path[j:i]
		}
		switch part {
		case "/":
			stack.PushBack(part)
			length += len(part)
		case ".":
		case "..":
			length -= len(stack.Back().Value.(string))
			stack.Remove(stack.Back())
		default:
			stack.PushBack(part)
			length += len(part)
		}
	}
	buf := bytes.NewBuffer(make([]byte, length))
	for p := stack.Front(); p != nil; p = p.Next() {
		v := p.Value.(string)
		buf.WriteString(v)
	}

	return buf.String()
}

func Test_simplifyPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			path: "/a/./b/../../c/",
			want: "/c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := simplifyPath(tt.path); got != tt.want {
				t.Errorf("simplifyPath(%v) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
