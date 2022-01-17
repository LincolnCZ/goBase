package config

// ISP Type
type ISPType = int

// 默认IP的选择依赖下面的数字大小排序
const (
	AUTO_DETECT ISPType = 0
	CTL         ISPType = 1     //电信
	CNC         ISPType = 2     //网通
	CNII        ISPType = 4     //铁通
	EDU         ISPType = 8     //教育网
	WBN         ISPType = 16    //长城宽带
	MOB         ISPType = 32    //移动
	BGP         ISPType = 64    //BGP
	ASIA        ISPType = 128   //亚洲
	SA          ISPType = 256   //南美
	EU          ISPType = 512   //欧洲
	NA          ISPType = 1024  //北美
	INTRANET    ISPType = 32768 //内部网
	MAX_ISP     ISPType = 65536 //最大值
)
