package main
import "os"
import "fmt"
import "net/http"
import "net/url"
import "io/ioutil"
import "io"
import "encoding/json"
import "text/template"
import "strings"
import "flag"

type Service_list struct{
	Status string
	Data []string
}
type Instance_list struct{
	Status string
	Data []string
}
type Promconf struct{
	Job_name string
	Instance []string
}
//get service list from api
func Get_service_list()[]string{
	var sl Service_list
	url:="http://jms.miaopai.com/deploy_prometheus/v1/service"
	resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
 
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    
    err = json.Unmarshal(body, &sl)  
	if err != nil {  
    	panic(err) 
	}  
    return sl.Data

}
//get instance list via service from api
func Get_Instance_via_service(service string)[]string{
	var il Instance_list
	post_url:="http://jms.miaopai.com/deploy_prometheus/v1/ip"
	resp, err := http.PostForm(post_url,url.Values{"service": {service}})
 
    if err != nil {
        panic(err) 
    }
 
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err) 
    }
    
    err = json.Unmarshal(body, &il)  
	if err != nil {  
    	panic(err) 
	}  
    return il.Data

}
//service_type: service|base
func Create_Prom_Config_file(service string,service_type string,tmpl_filename string)bool{
	var (
		pc Promconf
		job_name string
		conf_filename string
	)
    job_name = strings.Replace(service,".","_",-1)
	//create prometheus config file
	if service_type == "base"{
		conf_filename="promconfigs/prometheus_"+job_name+"_base.yml"
	}else{
		conf_filename="promconfigs/prometheus_"+job_name+".yml"
	}
	conf_f, err1 := os.OpenFile(conf_filename, os.O_WRONLY|os.O_CREATE, 0640) 
    if err1 != nil { 
    	fmt.Println(err1)
    	return false
    }
	pc.Job_name = job_name
	pc.Instance = Get_Instance_via_service(service)
	tmpl, err := template.ParseFiles(tmpl_filename)
	err = tmpl.Execute(io.Writer(conf_f), pc)
	if err != nil{
		fmt.Println(err)
    	return false
	}
	return true
}

func Show_all_service(){
	fmt.Println("ALL support service list:\n")
	for _,v:= range Get_service_list() {
		fmt.Println(v)
	}
}
func check_file_exists(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true}
    if os.IsNotExist(err) { return false}
    return true
}

func main(){
	 var (
	 	service = flag.String("service","all","service name which you want to create,defualt all ,show all support list")
        service_type = flag.String("service.type", "", "which type you want to create,defaut null|base")
        tmpl_name = flag.String("tmpl.file", "", "must,the path prometheus template filename,promtemplate.yml for service,promtemplate_base.yml for service's base")
     )
    flag.Parse()
    if *service == "all"{
    	Show_all_service()
    	return
    }
    /*
    if *service_type != "" {
    	fmt.Println("Not support this service type:"+*service_type)
    	return
    }
    */

    if !check_file_exists(*tmpl_name){
    	fmt.Println(*tmpl_name+":not exist")
    	return
    }
    if !Create_Prom_Config_file(*service,*service_type,*tmpl_name){
			fmt.Println(*service+" "+*service_type+":create config fail")
	}else{
			fmt.Println(*service+" "+*service_type+":create config success")
	}
}



