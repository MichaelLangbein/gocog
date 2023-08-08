package selfmade

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func fetchSize(fileUrl string) (int, error) {
	fmt.Printf("Getting size of %s", fileUrl)
	req, err := http.NewRequest(http.MethodHead, fileUrl, nil)
	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	contentLength := res.Header.Get("Content-Length")
	clInt, err := strconv.Atoi(contentLength)

	return clInt, nil
}

func fetchRange(fileUrl string, startByte int64, nrBytes int) ([]byte, error) {
	fmt.Printf("Fetching bytes %d-%d", startByte, startByte+int64(nrBytes))
	req, err := http.NewRequest(http.MethodGet, fileUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", startByte, startByte+int64(nrBytes)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		return nil, err2
	}

	return data, nil
}

type FetchingReader struct {
	fileUrl         string
	fetchBytes      int
	currentLocation int64
	fetchedData     map[int64][]byte
}

func MakeFetchingReader(fileUrl string) *FetchingReader {
	return &FetchingReader{
		fileUrl: fileUrl, fetchBytes: 4000, currentLocation: 0, fetchedData: map[int64][]byte{},
	}
}

func (r *FetchingReader) getKeysFor(start int64, length int) []int64 {
	nearest := int64(r.fetchBytes) * (start / int64(r.fetchBytes))
	keys := []int64{}
	for position := nearest; position < (start + int64(length)); position += int64(r.fetchBytes) {
		keys = append(keys, position)
	}
	return keys
}

func (r *FetchingReader) getDataForKey(key int64) ([]byte, error) {
	data, ok := r.fetchedData[key]
	if !ok {
		data, err := fetchRange(r.fileUrl, key, r.fetchBytes)
		if err != nil {
			return data, err
		}
		r.fetchedData[key] = data
		return data, nil
	}
	return data, nil
}

func (r *FetchingReader) getDataAt(off int64, nrBytes int) ([]byte, error) {
	keys := r.getKeysFor(off, nrBytes)

	outputData := make([]byte, nrBytes)
	outputPos := 0

	for j, key := range keys {
		keyData, err := r.getDataForKey(key)
		if err != nil {
			return outputData, err
		}

		var startIndex int64 = 0
		if j <= 0 {
			startIndex = off - key
		}
		var endIndex int64 = int64(r.fetchBytes)
		if j >= len(keys)-1 {
			endIndex = (off + int64(nrBytes)) - key
		}
		// if (startIndex == 0 && endIndex == int64(r.fetchBytes)) {
		// 	copy(outputData,
		// }
		for i := startIndex; i < endIndex; i++ {
			outputData[outputPos] = keyData[i]
			outputPos += 1
		}
	}

	return outputData, nil
}

/*
* ReadAt reads len(p) bytes into p starting at offset off in the underlying input source. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.
* When ReadAt returns n < len(p), it returns a non-nil error explaining why more bytes were not returned. In this respect, ReadAt is stricter than Read.
* Even if ReadAt returns n < len(p), it may use all of p as scratch space during the call. If some data is available but not len(p) bytes, ReadAt blocks until either all the data is available or an error occurs. In this respect ReadAt is different from Read.
* If the n = len(p) bytes returned by ReadAt are at the end of the input source, ReadAt may return either err == EOF or err == nil.
* If ReadAt is reading from an input source with a seek offset, ReadAt should not affect nor be affected by the underlying seek offset.
* Clients of ReadAt can execute parallel ReadAt calls on the same input source.
* Implementations must not retain p.
 */
func (r *FetchingReader) ReadAt(p []byte, off int64) (n int, err error) {
	nrBytes := len(p)
	data, err := r.getDataAt(off, nrBytes)
	if err != nil {
		return 0, err
	}
	copy(p, data)
	if len(data) < len(p) {
		return len(data), fmt.Errorf("something went wrong ... did you reach the end of the file?")
	}
	return len(data), nil
}

// Size() is used as a probe to determine wether the given key exists, and should return
// an error if no such key exists. The actual file size may or may not be effectively used
// depending on the underlying GDAL driver opening the file
//
// It may also optionally implement KeyMultiReader which will be used (only?) by
// the GTiff driver when reading pixels. If not provided, this
// VSI implementation will concurrently call ReadAt([]byte,int64)
func (r *FetchingReader) Size(key string) (int64, error) {
	size, err := fetchSize(r.fileUrl)
	return int64(size), err
}

/*
* Read reads up to len(p) bytes into p. It returns the number of bytes
* read (0 <= n <= len(p)) and any error encountered. Even if Read
* returns n < len(p), it may use all of p as scratch space during the call.
* If some data is available but not len(p) bytes, Read conventionally
* returns what is available instead of waiting for more.
*
* When Read encounters an error or end-of-file condition after
* successfully reading n > 0 bytes, it returns the number of
* bytes read. It may return the (non-nil) error from the same call
* or return the error (and n == 0) from a subsequent call.
* An instance of this general case is that a Reader returning
* a non-zero number of bytes at the end of the input stream may
* return either err == EOF or err == nil. The next Read should
* return 0, EOF.
*
* Callers should always process the n > 0 bytes returned before
* considering the error err. Doing so correctly handles I/O errors
* that happen after reading some bytes and also both of the
* allowed EOF behaviors.
*
* Implementations of Read are discouraged from returning a
* zero byte count with a nil error, except when len(p) == 0.
* Callers should treat a return of 0 and nil as indicating that
* nothing happened; in particular it does not indicate EOF.
*
* Implementations must not retain p.
 */
func (r *FetchingReader) Read(p []byte) (n int, err error) {
	off := r.currentLocation
	nrBytesRead, err := r.ReadAt(p, off)
	r.currentLocation += int64(nrBytesRead)
	return nrBytesRead, err
}

/*
* Seek sets the offset for the next Read or Write to offset,
* interpreted according to whence:
* SeekStart means relative to the start of the file,
* SeekCurrent means relative to the current offset, and
* SeekEnd means relative to the end
* (for example, offset = -2 specifies the penultimate byte of the file).
* Seek returns the new offset relative to the start of the
* file or an error, if any.
* Seeking to an offset before the start of the file is an error.
* Seeking to any positive offset may be allowed, but if the new offset exceeds
* the size of the underlying object the behavior of subsequent I/O operations
* is implementation-dependent.
 */
func (r *FetchingReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errors.New("Seek: invalid whence")
	case io.SeekStart:
		break
	case io.SeekCurrent:
		offset += r.currentLocation
		// case io.SeekEnd:
		// 	offset += s.limit
	}
	if offset < 0 {
		return 0, errors.New("Seek: invalid offset")
	}
	r.currentLocation = offset
	return offset, nil
}
