package raft

type Entry struct {
	Index   uint64 `json:"index"`
	Term    uint64 `json:"term"`
	Command []byte `json:"command"`
}

type Log struct {
	entries []Entry
}

func (l *Log) Append(term uint64, command []byte) Entry {
	entry := Entry{Index: uint64(len(l.entries) + 1), Term: term, Command: append([]byte(nil), command...)}
	l.entries = append(l.entries, entry)
	return entry
}

func (l *Log) Entries() []Entry {
	result := make([]Entry, len(l.entries))
	copy(result, l.entries)
	return result
}
