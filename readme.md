### Buffer Writes To Memory, Flush to Disk
Writes data into a fixed and pre-allocated buffer and flushes to disk as needed.

Multiple workers (default 4) are listening for writes. This results in flushes to disk only blocking the caller when all workers are flushing at the same time (and a buffered channel is used to further reduce the blocking window). Each worker maintains its own buffer. Each flush happens to a unique file. Flushes first happen in a temp file and are renamed to their final location. This allows one to leverage the rename's atomic behavior.

This package is used to collect statistics at our edge locations, and rsync the files to central servers where more intensive processing can take place. The logged data itself contains enough information to stitch the pieces together (for example, each message contains a `session_id` and `timestamp`, so that even if messages from the same `session_id` are distributed across multiple files (either from different workers, or different flushes) the flow of a session can be stiched back up).

We use this package with our [BytePool](https://github.com/viki-org/bytepool), but any type that implements `io.reader` works.

### Installation
Install using the "go get" command:

    go get github.com/viki-org/bufferedwriter

### Usage
Create a new `Buffer` instance and `Write` to it:

    var buffer = bufferedwriter.New(Configure())
    ...
    buffer.Write(myReadCloser)

### Configuration
The buffer is configured via a fluent interface:

    buffer := bufferedwriter.New(Configure().Size(65536).Prefix("hits_"))

Possible configuration options:

* `Path` (os.TempDir()) The path to store the flushed files
* `Permission` (0400) Permissions to apply to the file
* `Prefix` ("") A prefix to append to file names
* `Size` (65536) The size of the buffer in bytes. The resulting files might be larger than this, to accomodate for overflows
* `Temp` (os.TempDir()) The path to store temp files (either caused by a buffer overflow, or immediately proceeding a rename)
* `Timeout` (100 milliseconds) Time to wait before discarding a message should no worker pick it up
* `Workers` (4) The number of workers
