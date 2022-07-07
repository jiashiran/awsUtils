package check

import (
	"log"
	"spider/pkg"
	sshd "spider/pkg/ssh"
	"spider/pkg/store"
	"time"
)

func main() {
	//StoreCompnayInfo()
	BatchCheckCount()
}

func StoreCompnayInfo() {
	for {
		s := store.Pop()
		if s == "" {
			time.Sleep(5 * time.Second)
		} else {
			log.Println(s)
			pkg.ToTxt(s, "result.csv")
		}
	}
}

func BatchCheckCount() map[string]string {
	ips := []string{"39.107.243.195", "39.96.50.197"}
	sshd.InitServer(ips)
	return sshd.ExecuteBatch("cat result.csv | wc -l")
}

func BatchDownload() {
	ips := []string{"39.107.243.195", "39.96.50.197"}
	sshd.InitServer(ips)
	for _, ip := range ips {
		//os.MkdirAll(ip, os.ModePerm)
		sshd.DownloadFile(sshd.GetSftp(sshd.ServerMap[ip]), "result.csv", ip+"/")
	}
}

func BatchClean() {
	ips := []string{"39.107.243.195", "39.96.50.197"}
	sshd.InitServer(ips)
	sshd.ExecuteBatch("rm -rf result.csv")
}

func BatchSaveResult() {
	ips := []string{"39.107.243.195", "39.96.50.197"}
	sshd.InitServer(ips)
	pids := sshd.ExecuteBatch(`ps -ef|grep SaveResult | grep -v grep | awk -F " " '{print $2}' `)
	for ip, _ := range sshd.ServerMap {
		sshd.Execute(ip, "kill -9 "+pids[ip])
	}
	sshd.ExecuteBatch("rm -rf /var/tmp/SaveResult")
	sshd.ExecuteBatch("rm -rf startSaveResult.sh")
	sshd.UploadBatch("./pkg/chrome/SaveResult", "/var/tmp/")
	sshd.ExecuteBatch("chmod +x /var/tmp/SaveResult")
	sshd.ExecuteBatch(`echo "/var/tmp/SaveResult 1 >> /dev/null 2>&1 &" > startSaveResult.sh`)
	sshd.ExecuteBatch(`chmod +x startSaveResult.sh`)
	sshd.ExecuteBatch(`./startSaveResult.sh`)
	sshd.ExecuteBatch(`ps -ef|grep SaveResult | grep -v grep | awk -F " " '{print $2}' `)
}
