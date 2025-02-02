package widget

import (
	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/canvas"
	"github.com/ravenlab/fyne/driver/desktop"
	"github.com/ravenlab/fyne/internal/cache"
	"github.com/ravenlab/fyne/internal/widget"
	"github.com/ravenlab/fyne/theme"
)

// TreeNodeID represents the unique id of a tree node.
type TreeNodeID = string

const treeDividerHeight = 1

var _ fyne.Widget = (*Tree)(nil)

// Tree widget displays hierarchical data.
// Each node of the tree must be identified by a Unique TreeNodeID.
//
// Since: 1.4
type Tree struct {
	BaseWidget
	Root TreeNodeID

	ChildUIDs      func(uid TreeNodeID) (c []TreeNodeID)                     // Return a sorted slice of Children TreeNodeIDs for the given Node TreeNodeID
	CreateNode     func(branch bool) (o fyne.CanvasObject)                   // Return a CanvasObject that can represent a Branch (if branch is true), or a Leaf (if branch is false)
	IsBranch       func(uid TreeNodeID) (ok bool)                            // Return true if the given TreeNodeID represents a Branch
	OnBranchClosed func(uid TreeNodeID)                                      // Called when a Branch is closed
	OnBranchOpened func(uid TreeNodeID)                                      // Called when a Branch is opened
	OnSelected     func(uid TreeNodeID)                                      // Called when the Node with the given TreeNodeID is selected.
	OnUnselected   func(uid TreeNodeID)                                      // Called when the Node with the given TreeNodeID is unselected.
	UpdateNode     func(uid TreeNodeID, branch bool, node fyne.CanvasObject) // Called to update the given CanvasObject to represent the data at the given TreeNodeID

	branchMinSize fyne.Size
	leafMinSize   fyne.Size
	offset        fyne.Position
	open          map[TreeNodeID]bool
	scroller      *ScrollContainer
	selected      []TreeNodeID
}

// NewTree returns a new performant tree widget defined by the passed functions.
// childUIDs returns the child TreeNodeIDs of the given node.
// isBranch returns true if the given node is a branch, false if it is a leaf.
// create returns a new template object that can be cached.
// update is used to apply data at specified data location to the passed template CanvasObject.
//
// Since: 1.4
func NewTree(childUIDs func(TreeNodeID) []TreeNodeID, isBranch func(TreeNodeID) bool, create func(bool) fyne.CanvasObject, update func(TreeNodeID, bool, fyne.CanvasObject)) *Tree {
	t := &Tree{ChildUIDs: childUIDs, IsBranch: isBranch, CreateNode: create, UpdateNode: update}
	t.ExtendBaseWidget(t)
	return t
}

// NewTreeWithStrings creates a new tree with the given string map.
// Data must contain a mapping for the root, which defaults to empty string ("").
//
// Since: 1.4
func NewTreeWithStrings(data map[string][]string) (t *Tree) {
	t = &Tree{
		ChildUIDs: func(uid string) (c []string) {
			c = data[uid]
			return
		},
		IsBranch: func(uid string) (b bool) {
			_, b = data[uid]
			return
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return NewLabel("Template Object")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			node.(*Label).SetText(uid)
		},
	}
	t.ExtendBaseWidget(t)
	return
}

// CloseAllBranches closes all branches in the tree.
func (t *Tree) CloseAllBranches() {
	t.propertyLock.Lock()
	t.open = make(map[TreeNodeID]bool)
	t.propertyLock.Unlock()
	t.Refresh()
}

// CloseBranch closes the branch with the given TreeNodeID.
func (t *Tree) CloseBranch(uid TreeNodeID) {
	t.ensureOpenMap()
	t.propertyLock.Lock()
	t.open[uid] = false
	t.propertyLock.Unlock()
	if f := t.OnBranchClosed; f != nil {
		f(uid)
	}
	t.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (t *Tree) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	c := newTreeContent(t)
	s := NewScrollContainer(c)
	t.scroller = s
	r := &treeRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{s}),
		tree:         t,
		content:      c,
		scroller:     s,
	}
	s.onOffsetChanged = func() {
		if t.offset == s.Offset {
			return
		}
		t.offset = s.Offset
		c.Refresh()
	}
	r.updateMinSizes()
	r.content.viewport = r.MinSize()
	return r
}

// IsBranchOpen returns true if the branch with the given TreeNodeID is expanded.
func (t *Tree) IsBranchOpen(uid TreeNodeID) bool {
	if uid == t.Root {
		return true // Root is always open
	}
	t.ensureOpenMap()
	t.propertyLock.RLock()
	defer t.propertyLock.RUnlock()
	return t.open[uid]
}

// MinSize returns the size that this widget should not shrink below.
func (t *Tree) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// OpenAllBranches opens all branches in the tree.
func (t *Tree) OpenAllBranches() {
	t.ensureOpenMap()
	t.walkAll(func(uid string, branch bool, depth int) {
		if branch {
			t.propertyLock.Lock()
			t.open[uid] = true
			t.propertyLock.Unlock()
		}
	})
	t.Refresh()
}

// OpenBranch opens the branch with the given TreeNodeID.
func (t *Tree) OpenBranch(uid TreeNodeID) {
	t.ensureOpenMap()
	t.propertyLock.Lock()
	t.open[uid] = true
	t.propertyLock.Unlock()
	if f := t.OnBranchOpened; f != nil {
		f(uid)
	}
	t.Refresh()
}

// Resize sets a new size for a widget.
func (t *Tree) Resize(size fyne.Size) {
	t.propertyLock.RLock()
	s := t.size
	t.propertyLock.RUnlock()

	if s == size {
		return
	}

	t.propertyLock.Lock()
	t.size = size
	t.propertyLock.Unlock()

	t.Refresh() // trigger a redraw
}

// Select marks the specified node to be selected
func (t *Tree) Select(uid TreeNodeID) {
	if len(t.selected) > 0 {
		if uid == t.selected[0] {
			return // no change
		}
		if f := t.OnUnselected; f != nil {
			f(t.selected[0])
		}
	}
	t.selected = []TreeNodeID{uid}
	if t.scroller != nil {
		var found bool
		var y int
		var size fyne.Size
		t.walkAll(func(id TreeNodeID, branch bool, depth int) {
			m := t.leafMinSize
			if branch {
				m = t.branchMinSize
			}
			if id == uid {
				found = true
				size = m
			} else if !found {
				// Root node is not rendered unless it has been customized
				if t.Root == "" && id == "" {
					// This is root node, skip
					return
				}
				// If this is not the first item, add a divider
				if y > 0 {
					y += treeDividerHeight
				}

				y += m.Height
			}
		})
		if y < t.scroller.Offset.Y {
			t.scroller.Offset.Y = y
		} else if y+size.Height > t.scroller.Offset.Y+t.scroller.Size().Height {
			t.scroller.Offset.Y = y + size.Height - t.scroller.Size().Height
		}
		t.scroller.onOffsetChanged()
		// TODO Setting a node as selected should open all parents if they aren't already
	}
	t.Refresh()
	if f := t.OnSelected; f != nil {
		f(uid)
	}
}

// ToggleBranch flips the state of the branch with the given TreeNodeID.
func (t *Tree) ToggleBranch(uid string) {
	if t.IsBranchOpen(uid) {
		t.CloseBranch(uid)
	} else {
		t.OpenBranch(uid)
	}
}

// Unselect marks the specified node to be not selected
func (t *Tree) Unselect(uid TreeNodeID) {
	if len(t.selected) == 0 {
		return
	}

	t.selected = nil
	t.Refresh()
	if f := t.OnUnselected; f != nil {
		f(uid)
	}
}

func (t *Tree) ensureOpenMap() {
	t.propertyLock.Lock()
	defer t.propertyLock.Unlock()
	if t.open == nil {
		t.open = make(map[string]bool)
	}
}

func (t *Tree) walk(uid string, depth int, onNode func(string, bool, int)) {
	if isBranch := t.IsBranch; isBranch != nil {
		if isBranch(uid) {
			onNode(uid, true, depth)
			if t.IsBranchOpen(uid) {
				if childUIDs := t.ChildUIDs; childUIDs != nil {
					for _, c := range childUIDs(uid) {
						t.walk(c, depth+1, onNode)
					}
				}
			}
		} else {
			onNode(uid, false, depth)
		}
	}
}

// walkAll visits every open node of the tree and calls the given callback with TreeNodeID, whether node is branch, and the depth of node.
func (t *Tree) walkAll(onNode func(TreeNodeID, bool, int)) {
	t.walk(t.Root, 0, onNode)
}

var _ fyne.WidgetRenderer = (*treeRenderer)(nil)

type treeRenderer struct {
	widget.BaseRenderer
	tree     *Tree
	content  *treeContent
	scroller *ScrollContainer
}

func (r *treeRenderer) MinSize() (min fyne.Size) {
	min = r.scroller.MinSize()
	min = min.Max(r.tree.branchMinSize)
	min = min.Max(r.tree.leafMinSize)
	return
}

func (r *treeRenderer) Layout(size fyne.Size) {
	r.content.viewport = size
	r.scroller.Resize(size)
}

func (r *treeRenderer) Refresh() {
	r.updateMinSizes()
	s := r.tree.Size()
	if s.IsZero() {
		r.tree.Resize(r.tree.MinSize())
	} else {
		r.Layout(s)
	}
	r.scroller.Refresh()
	r.content.Refresh()
	canvas.Refresh(r.tree.super())
}

func (r *treeRenderer) updateMinSizes() {
	if f := r.tree.CreateNode; f != nil {
		r.tree.branchMinSize = newBranch(r.tree, f(true)).MinSize()
		r.tree.leafMinSize = newLeaf(r.tree, f(false)).MinSize()
	}
}

var _ fyne.Widget = (*treeContent)(nil)

type treeContent struct {
	BaseWidget
	tree     *Tree
	viewport fyne.Size
}

func newTreeContent(tree *Tree) (c *treeContent) {
	c = &treeContent{
		tree: tree,
	}
	c.ExtendBaseWidget(c)
	return
}

func (c *treeContent) CreateRenderer() fyne.WidgetRenderer {
	return &treeContentRenderer{
		BaseRenderer: widget.BaseRenderer{},
		treeContent:  c,
		branches:     make(map[string]*branch),
		leaves:       make(map[string]*leaf),
		branchPool:   &syncPool{},
		leafPool:     &syncPool{},
	}
}

func (c *treeContent) Resize(size fyne.Size) {
	c.propertyLock.RLock()
	s := c.size
	c.propertyLock.RUnlock()

	if s == size {
		return
	}

	c.propertyLock.Lock()
	c.size = size
	c.propertyLock.Unlock()

	c.Refresh() // trigger a redraw
}

var _ fyne.WidgetRenderer = (*treeContentRenderer)(nil)

type treeContentRenderer struct {
	widget.BaseRenderer
	treeContent *treeContent
	dividers    []fyne.CanvasObject
	branches    map[string]*branch
	leaves      map[string]*leaf
	branchPool  pool
	leafPool    pool
}

func (r *treeContentRenderer) Layout(size fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	branches := make(map[string]*branch)
	leaves := make(map[string]*leaf)

	offsetY := r.treeContent.tree.offset.Y
	viewport := r.treeContent.viewport
	width := fyne.Max(size.Width, viewport.Width)
	y := 0
	numDividers := 0
	// walkAll open branches and obtain nodes to render in scroller's viewport
	r.treeContent.tree.walkAll(func(uid string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}

		// If this is not the first item, add a divider
		if y > 0 {
			var divider fyne.CanvasObject
			if numDividers < len(r.dividers) {
				divider = r.dividers[numDividers]
			} else {
				divider = NewSeparator()
				r.dividers = append(r.dividers, divider)
			}
			divider.Move(fyne.NewPos(theme.Padding(), y))
			s := fyne.NewSize(width-2*theme.Padding(), treeDividerHeight)
			divider.Resize(s)
			divider.Show()
			y += treeDividerHeight
			numDividers++
		}

		m := r.treeContent.tree.leafMinSize
		if isBranch {
			m = r.treeContent.tree.branchMinSize
		}
		if y+m.Height < offsetY {
			// Node is above viewport and not visible
		} else if y > offsetY+viewport.Height {
			// Node is below viewport and not visible
		} else {
			// Node is in viewport
			var n *treeNode
			if isBranch {
				b, ok := r.branches[uid]
				if !ok {
					b = r.getBranch()
					if f := r.treeContent.tree.UpdateNode; f != nil {
						f(uid, true, b.Content())
					}
					b.update(uid, depth)
				}
				branches[uid] = b
				n = b.treeNode
			} else {
				l, ok := r.leaves[uid]
				if !ok {
					l = r.getLeaf()
					if f := r.treeContent.tree.UpdateNode; f != nil {
						f(uid, false, l.Content())
					}
					l.update(uid, depth)
				}
				leaves[uid] = l
				n = l.treeNode
			}
			if n != nil {
				n.Move(fyne.NewPos(0, y))
				n.Resize(fyne.NewSize(width, m.Height))
			}
		}
		y += m.Height
	})

	// Hide any dividers that haven't been reused
	for ; numDividers < len(r.dividers); numDividers++ {
		r.dividers[numDividers].Hide()
	}

	// Release any nodes that haven't been reused
	for uid, b := range r.branches {
		if _, ok := branches[uid]; !ok {
			b.Hide()
			r.branchPool.Release(b)
		}
	}
	for uid, l := range r.leaves {
		if _, ok := leaves[uid]; !ok {
			l.Hide()
			r.leafPool.Release(l)
		}
	}

	r.branches = branches
	r.leaves = leaves
}

func (r *treeContentRenderer) MinSize() (min fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	r.treeContent.tree.walkAll(func(uid string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}

		// If this is not the first item, add a divider
		if min.Height > 0 {
			min.Height += treeDividerHeight
		}

		m := r.treeContent.tree.leafMinSize
		if isBranch {
			m = r.treeContent.tree.branchMinSize
		}
		m.Width += depth * (theme.IconInlineSize() + theme.Padding())
		min.Width = fyne.Max(min.Width, m.Width)
		min.Height += m.Height
	})
	return
}

func (r *treeContentRenderer) Objects() (objects []fyne.CanvasObject) {
	r.treeContent.propertyLock.RLock()
	objects = r.dividers
	for _, b := range r.branches {
		objects = append(objects, b)
	}
	for _, l := range r.leaves {
		objects = append(objects, l)
	}
	r.treeContent.propertyLock.RUnlock()
	return
}

func (r *treeContentRenderer) Refresh() {
	s := r.treeContent.Size()
	if s.IsZero() {
		r.treeContent.Resize(r.treeContent.MinSize().Max(r.treeContent.tree.Size()))
	} else {
		r.Layout(s)
	}
	r.treeContent.propertyLock.RLock()
	for _, b := range r.branches {
		b.Refresh()
	}
	for _, l := range r.leaves {
		l.Refresh()
	}
	r.treeContent.propertyLock.RUnlock()
	canvas.Refresh(r.treeContent.super())
}

func (r *treeContentRenderer) getBranch() (b *branch) {
	o := r.branchPool.Obtain()
	if o != nil {
		b = o.(*branch)
	} else {
		var content fyne.CanvasObject
		if f := r.treeContent.tree.CreateNode; f != nil {
			content = f(true)
		}
		b = newBranch(r.treeContent.tree, content)
	}
	return
}

func (r *treeContentRenderer) getLeaf() (l *leaf) {
	o := r.leafPool.Obtain()
	if o != nil {
		l = o.(*leaf)
	} else {
		var content fyne.CanvasObject
		if f := r.treeContent.tree.CreateNode; f != nil {
			content = f(false)
		}
		l = newLeaf(r.treeContent.tree, content)
	}
	return
}

var _ desktop.Hoverable = (*treeNode)(nil)
var _ fyne.CanvasObject = (*treeNode)(nil)
var _ fyne.Tappable = (*treeNode)(nil)

type treeNode struct {
	BaseWidget
	tree     *Tree
	uid      string
	depth    int
	hovered  bool
	icon     fyne.CanvasObject
	isBranch bool
	content  fyne.CanvasObject
}

func (n *treeNode) Content() fyne.CanvasObject {
	return n.content
}

func (n *treeNode) CreateRenderer() fyne.WidgetRenderer {
	return &treeNodeRenderer{
		BaseRenderer: widget.BaseRenderer{},
		treeNode:     n,
		indicator:    canvas.NewRectangle(theme.BackgroundColor()),
	}
}

func (n *treeNode) Indent() int {
	return n.depth * (theme.IconInlineSize() + theme.Padding())
}

// MouseIn is called when a desktop pointer enters the widget
func (n *treeNode) MouseIn(*desktop.MouseEvent) {
	n.hovered = true
	n.partialRefresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (n *treeNode) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (n *treeNode) MouseOut() {
	n.hovered = false
	n.partialRefresh()
}

func (n *treeNode) Tapped(*fyne.PointEvent) {
	n.tree.Select(n.uid)
}

func (n *treeNode) partialRefresh() {
	if r := cache.Renderer(n.super()); r != nil {
		r.(*treeNodeRenderer).partialRefresh()
	}
}

func (n *treeNode) update(uid string, depth int) {
	n.uid = uid
	n.depth = depth
	n.propertyLock.Lock()
	n.Hidden = false
	n.propertyLock.Unlock()
	n.partialRefresh()
}

var _ fyne.WidgetRenderer = (*treeNodeRenderer)(nil)

type treeNodeRenderer struct {
	widget.BaseRenderer
	treeNode  *treeNode
	indicator *canvas.Rectangle
}

func (r *treeNodeRenderer) Layout(size fyne.Size) {
	x := 0
	y := 0
	r.indicator.Move(fyne.NewPos(x, y))
	s := fyne.NewSize(theme.Padding(), size.Height)
	r.indicator.SetMinSize(s)
	r.indicator.Resize(s)
	h := size.Height - 2*theme.Padding()
	x += theme.Padding() + r.treeNode.Indent()
	y += theme.Padding()
	if r.treeNode.icon != nil {
		r.treeNode.icon.Move(fyne.NewPos(x, y))
		r.treeNode.icon.Resize(fyne.NewSize(theme.IconInlineSize(), h))
	}
	x += theme.IconInlineSize()
	x += theme.Padding()
	if r.treeNode.content != nil {
		r.treeNode.content.Move(fyne.NewPos(x, y))
		r.treeNode.content.Resize(fyne.NewSize(size.Width-x-theme.Padding(), h))
	}
}

func (r *treeNodeRenderer) MinSize() (min fyne.Size) {
	if r.treeNode.content != nil {
		min = r.treeNode.content.MinSize()
	}
	min.Width += theme.Padding() + r.treeNode.Indent() + theme.IconInlineSize()
	min.Width += 2 * theme.Padding()
	min.Height = fyne.Max(min.Height, theme.IconInlineSize())
	min.Height += 2 * theme.Padding()
	return
}

func (r *treeNodeRenderer) Objects() (objects []fyne.CanvasObject) {
	if r.treeNode.content != nil {
		objects = append(objects, r.treeNode.content)
	}
	if r.treeNode.icon != nil {
		objects = append(objects, r.treeNode.icon)
	}
	objects = append(objects, r.indicator)
	return
}

func (r *treeNodeRenderer) Refresh() {
	if c := r.treeNode.content; c != nil {
		if f := r.treeNode.tree.UpdateNode; f != nil {
			f(r.treeNode.uid, r.treeNode.isBranch, c)
		}
	}
	r.partialRefresh()
}

func (r *treeNodeRenderer) partialRefresh() {
	if r.treeNode.icon != nil {
		r.treeNode.icon.Refresh()
	}
	if len(r.treeNode.tree.selected) > 0 && r.treeNode.uid == r.treeNode.tree.selected[0] {
		r.indicator.FillColor = theme.PrimaryColor()
	} else if r.treeNode.hovered {
		r.indicator.FillColor = theme.HoverColor()
	} else {
		r.indicator.FillColor = theme.BackgroundColor()
	}
	r.indicator.Refresh()
	canvas.Refresh(r.treeNode.super())
}

var _ fyne.Widget = (*branch)(nil)

type branch struct {
	*treeNode
}

func newBranch(tree *Tree, content fyne.CanvasObject) (b *branch) {
	b = &branch{
		treeNode: &treeNode{
			tree:     tree,
			icon:     newBranchIcon(tree),
			isBranch: true,
			content:  content,
		},
	}
	b.ExtendBaseWidget(b)
	return
}

func (b *branch) update(uid string, depth int) {
	b.treeNode.update(uid, depth)
	b.icon.(*branchIcon).update(uid, depth)
}

var _ fyne.Tappable = (*branchIcon)(nil)

type branchIcon struct {
	Icon
	tree *Tree
	uid  string
}

func newBranchIcon(tree *Tree) (i *branchIcon) {
	i = &branchIcon{
		tree: tree,
	}
	i.ExtendBaseWidget(i)
	return
}

func (i *branchIcon) Refresh() {
	if i.tree.IsBranchOpen(i.uid) {
		i.Resource = theme.MoveDownIcon()
	} else {
		i.Resource = theme.NavigateNextIcon()
	}
	i.Icon.Refresh()
}

func (i *branchIcon) Tapped(*fyne.PointEvent) {
	i.tree.ToggleBranch(i.uid)
}

func (i *branchIcon) update(uid string, depth int) {
	i.uid = uid
	i.Refresh()
}

var _ fyne.Widget = (*leaf)(nil)

type leaf struct {
	*treeNode
}

func newLeaf(tree *Tree, content fyne.CanvasObject) (l *leaf) {
	l = &leaf{
		&treeNode{
			tree:     tree,
			content:  content,
			isBranch: false,
		},
	}
	l.ExtendBaseWidget(l)
	return
}
