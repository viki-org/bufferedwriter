package bufferedwriter

import (
  "io"
  "os"
  "log"
  "time"
  "strconv"
)

type Worker struct {
  length int
  data []byte
  capacity int
  fileRoot string
  fileTemp string
  logger *log.Logger
  permission os.FileMode
  channel chan io.ReadCloser
}

func newWorker(id int, channel chan io.ReadCloser, configuration *Configuration) (*Worker) {
  idString := strconv.Itoa(id)
  w := &Worker{
    channel: channel,
    logger: configuration.logger,
    capacity: configuration.size,
    permission: configuration.permission,
    data: make([]byte, configuration.size),
    fileRoot: configuration.path,
    fileTemp: configuration.temp,
  }

  if w.fileRoot[len(w.fileRoot)-1:] != "/" {
    w.fileRoot += "/"
  }

  if w.fileTemp[len(w.fileTemp)-1:] != "/" {
    w.fileTemp += "/"
  }
  w.fileRoot += configuration.prefix + idString + "_"
  w.fileTemp += configuration.prefix + idString + ".tmp"
  return w
}

func (w *Worker) work() {
  os.Remove(w.fileTemp)
  for {
    message := <- w.channel
    w.process(message)
  }
}

func (w *Worker) process(message io.ReadCloser) {
  defer message.Close()
  var swapped bool
  for {
    read, err := message.Read(w.data[w.length:])
    w.length += read
    if err == io.EOF {
      if w.length == w.capacity || swapped {
        w.swap()
        w.save()
      }
      break
    }
    if err != nil {
      w.length -= read
      if w.logger != nil {
        w.logger.Print("Failed to read reader", err)
      }
    }
    if w.length == w.capacity { swapped = w.swap() }
  }
}

func (w *Worker) swap() bool {
  if w.length == 0 { return false }
  defer func() { w.length = 0 }()
  f, err := os.OpenFile(w.fileTemp, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0600)
  if err != nil {
    if w.logger != nil {
      w.logger.Print("Failed to create temporary buffered file", err)
    }
    return false
  }
  defer f.Close()
  f.Write(w.data[0:w.length])
  stat, _ := f.Stat()
  return stat.Size() >= int64(w.capacity)
}

func (w *Worker) save() {
  target := w.fileRoot + strconv.FormatInt(time.Now().UnixNano(), 10) + ".log"
  if err := os.Rename(w.fileTemp, target); err != nil {
    if w.logger != nil {
      w.logger.Printf("Failed to rename %v to %v. Error: %v", w.fileTemp, target, err)
    }
    return
  }
  os.Chmod(target, w.permission)
}
