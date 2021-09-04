package condition

import (
	"bytes"
	"fmt"
	"io"
)

func Stringify(root *Node) string {
	buf := bytes.NewBuffer(nil)
	stringify(buf, root, "")
	return buf.String()
}

func getNodeName(n *Node) string {
	if n == nil {
		return "NULL"
	}

	switch n.Type {
	case NodeTypeLiteral:
		return fmt.Sprintf("%s", n.Token.String())
	case NodeTypeArray:
		return "ARRAY"
	case NodeTypeFunction:
		return fmt.Sprintf("FUNCTION<%s>", n.Token.Value.(string))
	}

	return "UNKNOWN"
}

func stringify(w io.Writer, node *Node, prefix string) {
	fmt.Fprintf(w, "%s%s\n", prefix, getNodeName(node))
	if len(prefix) >= 4 {
		pos := len(prefix) - 4
		if prefix[pos] == '|' {
			prefix = prefix[:pos] + "|   "
		} else {
			prefix = prefix[:pos] + "    "
		}
	}

	if len(node.Children) > 0 {
		for i, child := range node.Children {
			if i == len(node.Children)-1 {
				stringify(w, child, prefix+" \\_ ")
			} else {
				stringify(w, child, prefix+"|__ ")
			}
		}
	}
}
