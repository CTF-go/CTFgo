package logs

//日志功能函数

import (
	c "CTFgo/configs"
	"fmt"
	"io"
	"os"
)

func Save_log(Log_path *os.File) {
	save_log_path, _ := os.Create(c.Save_log_path)
	defer save_log_path.Close()
	nbytes, err := io.Copy(Log_path, save_log_path)
	if err != nil {
		fmt.Printf("The copy operation failed %q\n", err)
	} else {
		fmt.Printf("Save %d bytes logs!\n", nbytes)
	}
}
