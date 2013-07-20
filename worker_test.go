package bufferedwriter

import (
  "io"
  "os"
  "bytes"
  "testing"
  "strings"
  "io/ioutil"
)

func TestGeneratesTheCorrectPaths(t *testing.T) {
  // test both trailing slash and not
  for _, path := range []string{"/home/test", "/home/test/"} {
    w := newWorker(12, nil, Configure().Path(path))
    expected := os.TempDir() + "12.tmp"
    if w.fileTemp != expected {
      t.Errorf("expecting fileTemp to be %q, but got %q", expected, w.fileTemp)
    }
    expected = "/home/test/12_"
    if w.fileRoot != expected {
      t.Errorf("expecting filePath to be %q, but got %q", expected, w.fileRoot)
    }
  }
}

func TestGeneratesTheCorrectPathsWithPrefix(t *testing.T) {
  // test both trailing slash and not
  for _, path := range []string{"/home/test", "/home/test/"} {
    w := newWorker(12, nil, Configure().Path(path).Prefix("bw_"))
    expected := os.TempDir() + "bw_12.tmp"
    if w.fileTemp != expected {
      t.Errorf("expecting fileTemp to be %q, but got %q", expected, w.fileTemp)
    }
    expected = "/home/test/bw_12_"
    if w.fileRoot != expected {
      t.Errorf("expecting filePath to be %q, but got %q", expected, w.fileRoot)
    }
  }
}

func TestGeneratesTheCorrectPathWhenTrailingSlashIsMissing(t *testing.T) {
  w := newWorker(12, nil, Configure().Path("/home/test"))
  expected := os.TempDir() + "12.tmp"
  if w.fileTemp != expected {
    t.Errorf("expecting fileTemp to be %q, but got %q", expected, w.fileTemp)
  }
  expected = "/home/test/12_"
  if w.fileRoot != expected {
    t.Errorf("expecting filePath to be %q, but got %q", expected, w.fileRoot)
  }
}

func TestBuffersWritesInMemory(t *testing.T) {
  expected := "There ain't no such thing as a free lunch"
  w := newWorker(1, nil, testConfig(100))
  w.process(message(expected[0:10]))
  w.process(message(expected[10:]))
  if string(w.data[0:w.length]) != expected {
    t.Errorf("Expected buffer to hold %q, but got %q", expected, w.data[0:w.length])
  } 
  assertNoIO(t)
}

func TestWriteExactSize(t *testing.T) {
  defer cleanup()
  expected := "There ain't no such thing as a free lunch"
  w := newWorker(1, nil, testConfig(len([]byte(expected))))
  w.process(message(expected))
  assertFile(t, expected, ".log")
  if w.length != 0 {
    t.Errorf("Expected buffer to have length of 0, got %d", w.length)
  }
}

func TestHandleMultipleFlushes(t *testing.T) {
  defer cleanup()
  w := newWorker(1, nil, testConfig(5))
  w.process(message("aaaa"))
  w.process(message("bbbbbb"))
  w.process(message("ccccc"))
  files := testFiles(".log")
  assertContent(t, files[0], "aaaabbbbbb")
  assertContent(t, files[1], "ccccc")
}

func TestClosesTheMessage(t *testing.T) {
  var closed bool
  w := newWorker(1, nil, testConfig(100))
  w.process(closeTracker(&closed, "a"))
  if closed == false {
    t.Error("message should have been closed")
  }
}

func testConfig(size int) *Configuration {
  return Configure().Size(size).Path("/tmp").Prefix("bufferwriter_")
}

func message(s string) io.ReadCloser {
  return ioutil.NopCloser(bytes.NewBufferString(s))
}

type ct struct {
  closed *bool
  io.Reader
}

func (c ct) Close() error {
  *c.closed = true
  return nil
}

func closeTracker(closed *bool, s string) io.ReadCloser {
  return ct{closed, bytes.NewBufferString(s)}
}

func assertNoIO(t *testing.T) {
  files := testFiles("*")
  if len(files) != 0 {
    t.Errorf("These files should not exist", files)
  }
}

func assertFile(t *testing.T, expected string, extension string) {
  tmp := testFiles(extension)
  if len(tmp) != 1 {
    t.Errorf("Expecting 1 %v file, got %d", extension, len(tmp))
  } else {
    assertContent(t, tmp[0], expected)
  }
}

func assertContent(t *testing.T, file string, expected string) {
  data, _ := ioutil.ReadFile("/tmp/" + file)
  if bytes.Compare(data, []byte(expected)) != 0 {
    t.Errorf("%v should container %v, but got %v", file, expected, string(data))
  }
}

func cleanup() {
  for _, name := range testFiles("*") { 
    os.Remove("/tmp/" + name) 
  }
}

func testFiles(extension string) []string {
  var matches []string
  files, _ := ioutil.ReadDir("/tmp")
  for _, file := range files {
    name := file.Name()
    if strings.HasPrefix(name, "bufferwriter_") && (extension == "*" || strings.HasSuffix(name, extension)) {
      matches = append(matches, name)
    }
  }
  return matches
}
