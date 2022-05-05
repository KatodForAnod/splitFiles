package main

type FileSystem struct {
	allFiles map[string]FileInfo
}

type FileInfo struct {
	fileStart   int64
	fileName    string
	fileFormat  string
	BytesRemove int64
}

/*func (s *FileSystem) SaveFile(info FileInfo, fileBody []byte) error {
	i := len(fileBody) / block_size

}*/
