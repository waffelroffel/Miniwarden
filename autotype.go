package main

import (
	"time"

	"github.com/lxn/walk"
)

type KeyPress struct {
	ctrl  bool
	alt   bool
	shift bool
	key   rune
}

func (kp *KeyPress) Press() {
	kp.keyDown()
	kp.keyUp()
}

func (kp *KeyPress) keyDown() {
	if kp.ctrl {
		keyDown(_VK_CTRL)
	}
	if kp.alt {
		keyDown(_VK_ALT)
	}
	if kp.shift {
		keyDown(_VK_SHIFT)
	}
	keyDown(kp.key)
}

func (kp *KeyPress) keyUp() {
	keyUp(kp.key)
	if kp.ctrl {
		keyUp(_VK_CTRL)
	}
	if kp.alt {
		keyUp(_VK_ALT)
	}
	if kp.shift {
		keyUp(_VK_SHIFT)
	}
}

const keyDelay = 100 * time.Millisecond

var (
	pasteKey = KeyPress{ctrl: true, key: offset('v')}
	tabKey   = KeyPress{key: offset(VK_TAB)}
	enterKey = KeyPress{key: offset(VK_RETURN)}
)

func AutoType(un, pw string) {
	Type(un)
	tabKey.Press()
	Type(session.Decrypt(pw))
	// enterKey.Press()
}

func Type(str string) {
	fatal(walk.Clipboard().SetText(str))
	pasteKey.Press()
	time.Sleep(keyDelay)
}

func offset(key rune) rune {
	if key >= 'a' && key <= 'z' {
		key -= 'a' - 'A'
	}
	return key
}
