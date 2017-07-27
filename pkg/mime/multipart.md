## multipart是指`multipart/form-data`，这是http协议中的一种`Content-type`，http的body中可以包含很多子类型，可以通过`boundary`的分界线分割。被`boundary`分割的叫做一个`part`。
```txt
Content-type: multipart/form-data;boundary=18bda0b61d5bec874e64cca1bee33ea1d1e53d5059871dda334488023748
```
1. **Writer**
`multipart.Writer`作用是构造和添加`multipart.Part`。
```go
func main(){
	bodybuf:=&bytes.Buffer{}
	body_writer:=multipart.NewWriter(bodybuf)
	body_part,err:=body_writer.CreateFormFile("upload_file","test.txt")
	if err!=nil{
		fmt.Println("Create form file err:",err)
		return 
	}
	body_part.Write([]byte("This is test file from test!"))
	body_writer.WriteField("this is key","this is value")
	body_writer.Close()
	fmt.Println(bodybuf)
}
```
输出:
```shell
--630d669802ff7cbdeafdffca6775db4410ac0328ebddc43ac3f6f7d15e16
Content-Disposition: form-data; name="upload_file"; filename="test.txt"
Content-Type: application/octet-stream

This is test file from test
--630d669802ff7cbdeafdffca6775db4410ac0328ebddc43ac3f6f7d15e16
Content-Disposition: form-data; name="this_is_key"

this is value
--630d669802ff7cbdeafdffca6775db4410ac0328ebddc43ac3f6f7d15e16--
```
2. **Reader**
`multipart.NewReader`有两个参数，第一个是需要解析的`[]byte`，第二个是`boundary`值。
当循环至最后一个`part`，再次`NextPart`的话，`err`会等于`io.EOF`。
```go
body_reader:=multipart.NewReader(bodybuf,body_writer.Boundary())
	for{
		part,err:=body_reader.NextPart()
		if err==io.EOF{
			break
		}
		content,err:=ioutil.ReadAll(part)
		fmt.Println("Form name",part.FormName(),"Content is:",string(content))
	}
```
输出:
```go
Form name upload_file Content is: This is test file from test
Form name this_is_key Content is: this is value
```