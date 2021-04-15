package logs

//日志功能函数

import (
	c "CTFgo/configs"
	"fmt"
	"os"
)

//Save_log在CTFgo退出时将run.log重命名为当前时间.log，保存在CTFgo/logs文件夹。
func Save_log() {
	err := os.Rename(c.Current_log_path, c.Save_log_path)
	if err != nil {
		fmt.Printf("the rename operation failed %q\n", err)
	} else {
		fmt.Println("save logs success!")
	}
}
