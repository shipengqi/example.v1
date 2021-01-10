package setting

import "time"

type server struct {
	runMode      string
	httpPort     int
	httpsPort    int
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (s *server) WriteTimeout() time.Duration {
	return s.writeTimeout
}

func (s *server) SetWriteTimeout(writeTimeout time.Duration) {
	s.writeTimeout = writeTimeout
}

func (s *server) ReadTimeout() time.Duration {
	return s.readTimeout
}

func (s *server) SetReadTimeout(readTimeout time.Duration) {
	s.readTimeout = readTimeout
}

func (s *server) HttpsPort() int {
	return s.httpsPort
}

func (s *server) SetHttpsPort(httpsPort int) {
	s.httpsPort = httpsPort
}

func (s *server) RunMode() string {
	return s.runMode
}

func (s *server) SetRunMode(runMode string) {
	s.runMode = runMode
}

func (s *server) HttpPort() int {
	return s.httpPort
}

func (s *server) SetHttpPort(httpPort int) {
	s.httpPort = httpPort
}