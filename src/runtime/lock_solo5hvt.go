// copied from lock_js.go

package runtime

const (
	mutex_unlocked = 0
	mutex_locked   = 1

	note_cleared = 0
	note_woken   = 1
	note_timeout = 2

	active_spin     = 4
	active_spin_cnt = 30
	passive_spin    = 1
)

func lock(l *mutex) {
	if l.key == mutex_locked {
		// js/wasm is single-threaded so we should never
		// observe this.
		throw("self deadlock")
	}
	gp := getg()
	if gp.m.locks < 0 {
		throw("lock count")
	}
	gp.m.locks++
	l.key = mutex_locked
}

func unlock(l *mutex) {
	if l.key == mutex_unlocked {
		throw("unlock of unlocked lock")
	}
	gp := getg()
	gp.m.locks--
	if gp.m.locks < 0 {
		throw("lock count")
	}
	l.key = mutex_unlocked
}

type noteWithTimeout struct {
	gp       *g
	deadline int64
}

var (
	notes            = make(map[*note]*g)
	notesWithTimeout = make(map[*note]noteWithTimeout)
)

func noteclear(n *note) {
	n.key = note_cleared
}

func notewakeup(n *note) {
	// gp := getg()
	if n.key == note_woken {
		throw("notewakeup - double wakeup")
	}
	cleared := n.key == note_cleared
	n.key = note_woken
	if cleared {
		goready(notes[n], 1)
	}
}

func notesleep(n *note) {
	throw("notesleep not supported by js")
}

func notetsleep(n *note, ns int64) bool {
	throw("notetsleep not supported by js")
	return false
}

// XXX
// same as runtime·notetsleep, but called on user g (not g0)
func notetsleepg(n *note, ns int64) bool {
	gp := getg()
	if gp == gp.m.g0 {
		throw("notetsleepg on g0")
	}

	const net0handle = 1
	const net0ready = 1 << net0handle
	packet := make([]byte, 1500)

	for ns > 0 {
		t0 := nanotime()

		length, ret := solo5Netread(net0handle, packet)
		// println("netreader, length=", length, ", ret=", ret)
		switch ret {
		case s5ok:
			println("netreader, packet length=", length)
			// xxx hand off to network stack
		case s5again:
			// println("netreader, err again")
			readySet, events := solo5Poll(uint64(ns))
			println("notetsleepg, readyset=", readySet, "events=", events)
		case s5invalid:
			throw("netreader s5invalid")
		case s5unspec:
			throw("netreader s5unspec")
		default:
			println("netreader, error", ret)
			throw("netreader unknown error")
		}
		ns -= nanotime() - t0

		/*
			t0 := nanotime()
			readySet, events := solo5Poll(uint64(ns))
			println("notetsleepg, readyset=", readySet, "events=", events)
			if events > 0 && readySet != 0 {
				println("readying netreader")
				goready(netreaderG, 1)
				gopark(nil, nil, waitReasonSleep, traceEvNone, 1)
				println("back in notetsleepg")
			}
			ns -= nanotime() - t0
		*/

		/*
			deadline := nanotime() + ns

			id := scheduleTimeoutEvent(delay)
			mp := acquirem()
			notes[n] = gp
			notesWithTimeout[n] = noteWithTimeout{gp: gp, deadline: deadline}
			releasem(mp)

			gopark(nil, nil, waitReasonSleep, traceEvNone, 1)

			clearTimeoutEvent(id) // note might have woken early, clear timeout
			mp = acquirem()
			delete(notes, n)
			delete(notesWithTimeout, n)
			releasem(mp)

			return n.key == note_woken
		*/
	}

	/*
		for n.key != note_woken {
			mp := acquirem()
			notes[n] = gp
			releasem(mp)

			gopark(nil, nil, waitReasonZero, traceEvNone, 1)

			mp = acquirem()
			delete(notes, n)
			releasem(mp)
		}
	*/

	return true
}

// checkTimeouts resumes goroutines that are waiting on a note which has reached its deadline.
func checkTimeouts() {
	now := nanotime()
	for n, nt := range notesWithTimeout {
		if n.key == note_cleared && now >= nt.deadline {
			n.key = note_timeout
			goready(nt.gp, 1)
		}
	}
}

func beforeIdle() bool {
	return false
}

/*
XXX

var netreaderG *g

func init() {
	go netreader()
}

func netreader() {
	netreaderG = getg()

	const net0handle = 1
	const net0ready = 1<<net0handle
	packet := make([]byte, 1500)
	for {
		length, ret := solo5Netread(net0handle, packet)
		// println("netreader, length=", length, ", ret=", ret)
		switch ret {
		case s5ok:
			println("netreader, packet length=", length)
			// xxx hand off to network stack
		case s5again:
			// println("netreader, err again")
			gopark(nil, nil, waitReasonIOWait, traceEvGoBlockNet, 1)
		case s5invalid:
			panic("netreader s5invalid")
		case s5unspec:
			panic("netreader s5unspec")
		default:
			print("netreader, error %x", ret)
			panic("netreader unknown error")
		}
	}
}
*/
