package dir

import(
	"os"
	"os/user"
	"runtime"
)

// 返回用户家目录
func HomeDir()string{

	if(runtime.GOOS == "windows"){
		return os.Getenv("USERPROFILE")
	}else{
		if v:=os.Getenv("HOME");v!=""{
			return v
		}
	
		if u,err:= user.Current();err==nil{
			return u.HomeDir
		}
	}

	return ""
}

// 判断所给路径文件或文件夹是否存在 
func IsExists(path string) bool {  
	_, err := os.Stat(path)    //os.Stat获取文件信息  
	if err != nil {  
		if os.IsExist(err) {  
			return true  
		}  
		return false  
	}  
	return true  
  }  
  