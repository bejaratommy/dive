package view

import (
	"errors"

	"github.com/awesome-gocui/gocui"
)

// CursorDown moves the cursor down in the currently selected gocui pane, scrolling the screen as needed.
func CursorDown(g *gocui.Gui, v *gocui.View) error {
	return CursorStep(g, v, 1)
}

// CursorUp moves the cursor up in the currently selected gocui pane, scrolling the screen as needed.
func CursorUp(g *gocui.Gui, v *gocui.View) error {
	return CursorStep(g, v, -1)
}

// CursorStep moves the cursor the given step distance within the view's buffer,
// scrolling the origin so the cursor stays within the visible area.
func CursorStep(g *gocui.Gui, v *gocui.View, step int) error {
	cx, cy := v.Cursor()
	targetY := cy + step

	// don't move onto a nonexistent or empty line
	line, err := v.Line(targetY)
	if err != nil {
		return err
	}
	if len(line) == 0 {
		return errors.New("unable to move the cursor, empty line")
	}

	if err := v.SetCursor(cx, targetY); err != nil {
		return err
	}

	// gocui reports the cursor position relative to the view's buffer, so keep
	// the origin in sync to ensure the cursor remains visible while scrolling.
	ox, oy := v.Origin()
	_, height := v.Size()
	if newOy := scrollOrigin(oy, height, targetY); newOy != oy {
		return v.SetOrigin(ox, newOy)
	}
	return nil
}

// scrollOrigin returns the vertical origin required to keep the buffer line
// target within a viewport of the given height. It returns origin unchanged
// when target is already visible.
func scrollOrigin(origin, height, target int) int {
	if target < origin {
		return target
	}
	if height > 0 && target >= origin+height {
		return target - height + 1
	}
	return origin
}
