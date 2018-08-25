package cmap

func hash(str string)uint64{
	seed:=uint64(13131)
	var hash uint64
	for i:=0;i<len(str);i++{
		hash=hash*seed+uint64(str[i])
	}
	return (hash & 0x7FFFFFFFFFFFFFFF)

	//h:=md5.Sum([]byte(str))
	//var num uint64
	//binary.Read(bytes.NewReader(h[:]),binary.LittleEndian,&num)
	//return num
}
