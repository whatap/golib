package cmdutil

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"

	//"log"
	"os"
	"os/exec"

	//"syscall"
	//"runtime/debug"
	"runtime"
	"strconv"
	"strings"

	//"gitlab.whatap.io/go/agent/util/logutil"
	"github.com/whatap/golib/util/stringutil"
)

// Pipeline strings together the given exec.Cmd commands in a similar fashion
// to the Unix pipeline.  Each command's standard output is connected to the
// standard input of the next command, and the output of the final command in
// the pipeline is returned, along with the collected standard error of all
// commands and the first error found (if any).
//
// To provide input to the pipeline, assign an io.Reader to the first's Stdin.
func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	defer func() {
		for _, cmd := range cmds {
			//syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)
			cmd.Process.Kill()
		}
	}()
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			//logutil.Infoln("cmd.StdoutPipe() cmd=", cmd.Path, " error", err)
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start and Wait each command
	// 2018.8.21 먼저 Start를 다 시키고 나서, Wait 실행 중 중간 cmd 에서 에러가 나면 나머지 cmd 는 좀비 프로세스로 변함
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			//logutil.Println("WA30200", "PipeLine Start Error, Path=", cmd.Path, ", err=", err)
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			//logutil.Println("WA30201", "PipeLine Wait Error, Path=", cmd.Path, ", err=", err)
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}

func GetPHPInfo() map[string]string {
	defer func() {
		// recover
		if r := recover(); r != nil {
			//
			//log.Println("recover:", r, string(debug.Stack()))
		}
	}()
	m := make(map[string]string)
	phpinfo := cmdPHPInfo()
	phpinfo = strings.Replace(phpinfo, "\r", "", -1)

	// PHP Version
	phpVersion := stringutil.Substring(phpinfo, "PHP Version", "\n\nConfiguration\n\n")
	s1 := strings.Split(phpVersion, "\n")

	for _, tmp := range s1 {
		k, v := stringutil.ToPair(tmp, "=>")
		if k != "" {
			m[k] = v
		}
	}

	return m
}

func GetPHPModuleInfo() map[string]string {
	defer func() {
		// recover
		if r := recover(); r != nil {
			//
			//log.Println("recover:", r, string(debug.Stack()))
		}
	}()
	m := make(map[string]string)
	keysList := list.New()

	//	pos := -1
	//	pos1 := -1
	//	mpos := -1
	//	mpos1 := -1

	//php -m
	moduleinfo := cmdPHPModuleInfo()
	moduleinfo = strings.Replace(moduleinfo, "\r", "", -1)

	phpmodules := stringutil.Substring(moduleinfo, "[PHP Modules]", "[Zend Modules]")
	s1 := strings.Split(phpmodules, "\n")
	// key 등록
	for _, tmp := range s1 {
		if strings.TrimSpace(tmp) != "" {
			m[tmp] = ""
			keysList.PushBack(tmp)
			//log.Println("PHP Module key= ", tmp)
		}
	}

	zendmodules := stringutil.Substring(moduleinfo, "[Zend Modules]", "")

	s2 := strings.Split(zendmodules, "\n")
	// key 등록
	for _, tmp := range s2 {
		if strings.TrimSpace(tmp) != "" {
			m[tmp] = ""
			keysList.PushBack(tmp)
			//log.Println("Zend Module key= ", tmp)
		}
	}

	mLen := keysList.Len()
	keys := make([]string, mLen)
	idx := 0
	for e := keysList.Front(); e != nil; e = e.Next() {
		//log.Println("keys=", e)
		keys[idx] = string(e.Value.(string))
		idx++
	}

	//phpI := exec.Command(php, "-i")
	phpinfo := cmdPHPInfo()
	phpinfo = strings.Replace(phpinfo, "\r", "", -1)

	// Configuration
	str := stringutil.Substring(phpinfo, "\n\nConfiguration\n\n", "\n\nAdditional Modules\n\n")
	//log.Println("Configuration=", str)
	for i := 0; i < mLen; i++ {
		detail := ""
		if i+1 < mLen {
			detail = stringutil.Substring(str, "\n\n"+keys[i], "\n\n"+keys[i+1])
		} else {
			detail = stringutil.Substring(str, "\n\n"+keys[i], "")
		}

		s3 := stringutil.Tokenizer(detail, "\n")
		m[keys[i]] = strings.Join(s3, ", ")
	}

	return m
}

func cmdPHPInfo() string {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	php := os.Getenv("WHATAP_PHP_BIN")
	if strings.TrimSpace(php) != "" {
		cmd := exec.Command(php, "-i")
		out, err := cmd.Output()

		if err != nil {
			//error
			//log.Println("command err", err)
			return ""
		}

		return string(out)
	}
	return ""
}

func cmdPHPModuleInfo() string {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	php := os.Getenv("WHATAP_PHP_BIN")
	if strings.TrimSpace(php) != "" {
		cmd := exec.Command(php, "-m")
		out, err := cmd.Output()

		if err != nil {
			//error
			//log.Println("command err", err)
			return ""
		}
		return string(out)
	}
	return ""
}

func GetPstackInfo(pid int32) map[int64]string {
	rt := make(map[int64]string, 0)

	// parse , linux, freebsd
	// nts 쓰레드 아이디 없이 스택이 1개  -1 쓰레드 아이디로 설정
	// zts 쓰레드 아이디
	// freebasd nts 인 경우 쓰레드 아이디가 정상(freebsd10), 쓰레드 아이디는 -1 (freebsd11), 스택은 1개.
	//linux
	//# pstack 1720
	//Thread 27 (Thread 0x7f3842d3c700 (LWP 1722)):
	//#0  0x00007f38580b568c in pthread_cond_wait@@GLIBC_2.3.2 () from /lib64/libpthread.so.0
	//#1  0x00007f38598215ed in ap_queue_pop ()
	//#2  0x00007f385981fc54 in ?? ()
	//#3  0x00007f38580b1aa1 in start_thread () from /lib64/libpthread.so.0
	//#4  0x00007f3857dfec4d in clone () from /lib64/libc.so.6
	//Thread 26 (Thread 0x7f384233b700 (LWP 1723)):
	//#0  0x00007f38580b568c in pthread_cond_wait@@GLIBC_2.3.2 () from /lib64/libpthread.so.0
	//#1  0x00007f38598215ed in ap_queue_pop ()
	//#2  0x00007f385981fc54 in ?? ()
	//#3  0x00007f38580b1aa1 in start_thread () from /lib64/libpthread.so.0
	//#4  0x00007f3857dfec4d in clone () from /lib64/libc.so.6

	//freebsd
	//	pstack -O 95874
	//95874: /usr/local/sbin/httpd
	//----------------- thread 100967 (running) -----------------
	// 0x801c5a228 __sys_flock (21539d6, 8, ffffffff, 8, 1e, 0) + 8
	// 0x8021539d6 _init (215364c, 8, 0, 0, 0, 0) + 238e
	// 0x80215364c _init (2152720, 8, 2881200, 8, 14, 0) + 2004
	// 0x802152720 _init (43bdbd, 0, 28ed098, 8, 2820118, 8) + 10d8
	//    0x43bdbd ap_run_mpm (43467f, 0, ffffed10, 7fff, ffffed30, 7fff) + 3d
	//    0x43467f main (433ccf, 0, 433b60, 0, 0, 0) + 8bf
	//    0x433ccf _start (6b1000, 8, 0, 0, 0, 0) + 16f

	out, err := cmdPstackInfo(pid)

	if err != nil {
		//logutil.Println("WA30202", "Error GetPstackInfo ", err)
		return nil
	}

	// DEBUG
	//rt[int64(pid)] = string(out)
	//logutil.Infoln("GetPstackInfo", "pid = ", pid, ", out=", string(out))
	r := bufio.NewScanner(strings.NewReader(string(out)))
	tidCount := 0
	tid := int64(0)
	sb := stringutil.NewStringBuffer()
	for r.Scan() {
		line := r.Text()
		if runtime.GOOS == "linux" {
			if strings.HasPrefix(line, "Thread ") {
				tidCount++
				// 저장.
				if sb.ToString() != "" {
					if tid == 0 {
						rt[-1] = sb.ToString()
					} else {
						rt[tid] = sb.ToString()
					}
				}
				v := stringutil.Substring(line, "(LWP", "))")
				tid, _ = strconv.ParseInt(strings.TrimSpace(v), 10, 64)
				sb.Clear()
			} else {
				if strings.Index(line, "whatap_") < 0 {
					//funcName := stringutil.Substring(strings.TrimSpace(line), "in ", "(")
					funcName := line
					if funcName != "" {
						//if funcName != "init" {
						sb.AppendLine(funcName)
						//}
					}
				}
			}
		} else if runtime.GOOS == "freebsd" {
			if strings.HasPrefix(line, "----------------- thread") {
				tidCount++
				// 저장.
				if sb.ToString() != "" {
					if tid == 0 {
						rt[-1] = sb.ToString()
					} else {
						rt[tid] = sb.ToString()
					}
				}
				v := stringutil.Substring(line, "thread", "(running)")
				//logutil.Infoln("GetPstackInfo", "thread=", v, ",line=", line)
				tid, _ = strconv.ParseInt(strings.TrimSpace(v), 10, 64)
				sb.Clear()
			} else {
				if strings.Index(line, "whatap_") < 0 {
					//funcName := stringutil.Substring(strings.TrimSpace(line), " ", "(")
					funcName := line
					if funcName != "" {
						//if funcName != "init" {
						sb.AppendLine(funcName)
						//}
					}
				}
			}
		}
	}

	if sb.ToString() != "" {
		if tid == 0 {
			// 쓰레드 없는 경우. -1 값으로 통일
			rt[-1] = sb.ToString()
			//logutil.Infoln("GetPstakInfo", "tail none tidr stack=", rt[-1])
		} else {
			rt[tid] = sb.ToString()
			//logutil.Infoln("GetPstakInfo", "tail tid=", tid, ",stack=", rt[tid])

		}
	}

	return rt
}

func cmdPstackInfo(pid int32) (out []byte, err error) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("pstack", fmt.Sprintf("%d", int(pid)))
		out, err = cmd.Output()
	} else if runtime.GOOS == "freebsd" {
		cmd := exec.Command("pstack", "-O", fmt.Sprintf("%d", int(pid)))
		out, err = cmd.Output()
	}
	return out, err
}

// Get docker full id from /proc/self/cgroup
func GetDockerFullId() string {

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	// check exists /proc/self/cgroup
	if _, err := os.Stat("/proc/self/cgroup"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return ""
	}
	//cat /proc/self/cgroup | head -n 1 | cut -d '/' -f3
	c1 := exec.Command("cat", "/proc/self/cgroup")
	c2 := exec.Command("head", "-n", "1")
	c3 := exec.Command("cut", "-d", "/", "-f3")

	// Run the pipeline
	out, _, err := Pipeline(c1, c2, c3)
	if err != nil {
		//logutil.Println("WA30203", "GetDockerFullId Error : errors : ", err)
		return ""
	}
	return strings.TrimSuffix(string(out), "\n")
}

func GetLinuxProductUUID() string {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("cat", "/sys/class/dmi/id/product_uuid")
		out, err := cmd.Output()
		if err == nil {
			return string(out)
		}
	}
	return ""
}

func CMDMain() {
	c1 := exec.Command("ps", "aux")
	c2 := exec.Command("grep", "httpd")
	c3 := exec.Command("awk", "{print $3}")
	c4 := exec.Command("awk", "{total = total + $1} END {print total}")

	// Run the pipeline
	//output, stderr, err := Pipeline(c1, c2, c3, c4)
	output, _, err := Pipeline(c1, c2, c3, c4)
	if err != nil {
		//logutil.Printf("Error : %s", err)
	}

	// Print the stdout, if any
	if len(output) > 0 {
		//logutil.Printf("output %s", output)

	}
}
