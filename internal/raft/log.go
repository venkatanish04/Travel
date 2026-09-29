package raft

type Entry struct {
	Index   uint64 `json:"index"`
	Term    uint64 `json:"term"`
	Command []byte `json:"command"`
}

type LogEntry = Entry

type Log struct {
	Entries []Entry
}

func NewRaftLog() *Log { return &Log{Entries: make([]Entry, 0)} }

func (l *Log) Append(args ...interface{}) Entry {
	var term uint64
	var command []byte
	if len(args) == 1 {
		entry := args[0].(Entry)
		term = entry.Term
		command = append([]byte(nil), entry.Command...)
	} else {
		term = args[0].(uint64)
		command = append([]byte(nil), args[1].([]byte)...)
	}
	entry := Entry{Index: uint64(len(l.Entries) + 1), Term: term, Command: command}
	l.Entries = append(l.Entries, entry)
	return entry
}

func (l *Log) LastIndex() int {
	return len(l.Entries) - 1
}

func (l *Log) LastTerm() int {
	if len(l.Entries) == 0 {
		return -1
	}
	return int(l.Entries[len(l.Entries)-1].Term)
}

func (l *Log) Get(index int) Entry { return l.Entries[index] }

func (l *Log) DeleteFrom(index int) {
	if index < 0 || index >= len(l.Entries) {
		return
	}
	l.Entries = l.Entries[:index]
}
