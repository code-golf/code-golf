package terminal

var emptyNode = node{blob: ' ', style: &emptyStyle}

type node struct {
	blob  rune
	style *style
	elem  *element
}

func (n *node) hasSameStyle(o node) bool {
	return n.style.isEqual(o.style)
}
