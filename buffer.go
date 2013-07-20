package bufferedwriter
// Package bufferedwriter buffers message into memory 
// before flushing them to disk

import (
  "io"
  "time"
)

type Buffer struct {
  configuration *Configuration
  channel chan io.ReadCloser
  Workers []*Worker
}

func New(configuration *Configuration) *Buffer {
  channel := make(chan io.ReadCloser, 512)
  buffer := &Buffer{
    channel: channel,
    configuration: configuration,
    Workers: make([]*Worker, configuration.workers),
  }

  for i := 0; i < configuration.workers; i++ {
    buffer.Workers[i] = newWorker(i, channel, configuration)
    go buffer.Workers[i].work()
  }
  return buffer
}

func (b *Buffer) Write(message io.ReadCloser) bool {
  select {
    case b.channel <- message:
      return true
    case <- time.After(b.configuration.timeout):
      message.Close()
      return false
  }
}

func (b *Buffer) Flush() {
  for _, w := range b.Workers {
    w.save()
  }
}
