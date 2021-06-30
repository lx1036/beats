package tests

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(test *testing.T) {
	source := "/var/lib/docker/containers/e6ef92c9e7e54e5f8850c357149e8734da7cdb96f5246643c593be3c99df1178/e6ef92c9e7e54e5f8850c357149e8734da7cdb96f5246643c593be3c99df1178-json.log"
	sourceLen := len(source)
	const containerIdLen = 64

	LogsPath := "/var/lib/docker/containers/"
	logsPathLen := len(LogsPath)

	if sourceLen >= logsPathLen+containerIdLen {
		cid := source[logsPathLen : logsPathLen+containerIdLen]
		fmt.Println(cid)
	}




	source = "/data/docker/containers/e6ef92c9e7e54e5f8850c357149e8734da7cdb96f5246643c593be3c99df1178/e6ef92c9e7e54e5f8850c357149e8734da7cdb96f5246643c593be3c99df1178-json.log"
	sourceLen = len(source)
	LogsPath = "/data/docker/containers/"
	logsPathLen = len(LogsPath)
	if sourceLen >= logsPathLen+containerIdLen {
		cid := source[logsPathLen : logsPathLen+containerIdLen]
		fmt.Println(cid)
	}
	/*if  sourceLen >= containerIdLen+4 {
		containerIDEnd := sourceLen - 4
		fmt.Println(containerIDEnd-containerIdLen, containerIDEnd)
		cid := source[containerIDEnd-containerIdLen : containerIDEnd]
		fmt.Println(cid)
	}*/
}

func TestName2(test *testing.T) {
	fmt.Println(strings.Fields("a b c"))
}
