package bufferedwriter

import (
  "os"
  "log"
  "time"
)

type Configuration struct {
  size int
  workers int
  temp string
  path string
  prefix string
  logger *log.Logger
  timeout time.Duration
  permission os.FileMode
}

func Configure() *Configuration {
  return &Configuration{
    workers: 4,
    size: 65536,
    permission: 0400,
    path: os.TempDir(),
    temp: os.TempDir(),
    timeout: time.Millisecond * 100,
  }
}

func (c *Configuration) Size(size int) (*Configuration) {
  c.size = size
  return c
}

func (c *Configuration) Workers(count int) (*Configuration) {
  c.workers = count
  return c
}

func (c *Configuration) Path(path string) (*Configuration) {
  c.path = path
  return c
}

func (c *Configuration) Temp(temp string) (*Configuration) {
  c.temp = temp
  return c
}

func (c *Configuration) Prefix(prefix string) (*Configuration) {
  c.prefix = prefix
  return c
}

func (c *Configuration) Logger(logger *log.Logger) (*Configuration) {
  c.logger = logger
  return c
}

func (c *Configuration) Timeout(timeout time.Duration) (*Configuration) {
  c.timeout = timeout
  return c
}

func (c *Configuration) Permission(permission os.FileMode) (*Configuration) {
  c.permission = permission
  return c
}
