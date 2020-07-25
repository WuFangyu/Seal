package main

const (
	msgDelimiter = "@"
	defaultChunkSize = 307200
	maxChunk = 80
	bufSize = 4096
	portStart = 9000
	portEnd = 65535
	okSignal = "ok"
	errSignal = "err"
	encExt = ".enc"
	tmpDir = "/tmp"
	keylenECC = 66
	keylenAES = 32
	levelDBpath = "/tmp/level_db"
	hostPubKeyIdx = "host"
	hostPrivKeyIdx = "host_private"
	infoColor    = "\033[1;34m%s\033[0m"
	noticeColor  = "\033[1;36m%s\033[0m"
	warningColor = "\033[1;33m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
	debugColor   = "\033[0;36m%s\033[0m"
	sendCmdErrMsg = "Use `seal send -help` for more info."
	recvCmdErrMsg = "Use `seal recv -help` for more info."
	keyCmdErrMsg = "Use `seal key -help` for more info."
)
