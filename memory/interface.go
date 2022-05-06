package memory

type Mem interface {
	LoadFile(info FileInfo) ([]byte, error)
	SaveFile(fileName string, fileBody []byte) (FileInfo, error)
}
