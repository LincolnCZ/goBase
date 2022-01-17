#ifndef __S2S_CLIENT_IF_HEADER__
#define __S2S_CLIENT_IF_HEADER__
/***************************************************
Version Record:
current version: 1.0.0
version               modification list                                   date
0.0.0                 initial version                                   	-
1.0.0                 add API:                                              2014/03/19
			onSessionLost
			onSessionRecovery
			getAllInfos
			add S2sDataEncoder/S2sDataDecoder
1.1.0                Add API to config console port 
1.2.0                Add API readIpList and WriteIpList
1.3.0                define  __STDC_FORMAT_MACROS
1.4.0                Add API delMine
2.0.0                re-design the APIs, to make it more simple
2.1.0:  		add configRecoverTime api
2.2.0: 		封装bson库，防止冲突
2.2.1:      内部log, 代码优化；
2.3.0: 		a.修复重连时m_condition的sql语法错误的bug b.重连尝试次数 RECONN_THRESHOLD_SC改为2 c.修复mini version比较的bug
3.0.0:		增加新接口: a.是否拉取完全量数据(可用于实现拉取到所有数据再注册的逻辑)
						b.是否拉取到某个服务的所有数据
						c.阻塞订阅接口
3.0.1:		修复应用线程在收到事件通知时isPullAllSub接口可能返回false的bug.(多线程时序问题)
****************************************************/
#ifndef __STDC_FORMAT_MACROS
#define __STDC_FORMAT_MACROS
#endif

#include <inttypes.h>
#include <map>
#include <set>
#include <vector>
#include <string>


#define S2S_VERSION "3.0.1"

namespace S2S
{
enum ISPType{
	AUTO_DETECT = 0,
	CTL  = 1,	 //电信
	CNC  = 2,	 //网通
	MULTI= 3,  //双线
	CNII = 4,	 //铁通
	EDU  = 8,	 //教育网
	WBN  = 16, //长城宽带
	MOB  = 32, //移动
	BGP  = 64,  //BGP
	HK  = 128, // 香港
	BRA = 256, // 巴西
};
}

enum LostCheckType
{
   NoLostCheck_C =0,
//   SinglePointTcpCheck_C,
//   SinglePointUdpCheck_C,   // not supported now.
   MulPointTcpCheck_C =3,
//   MulPointUdpCheck_C,      //  not supported now
   DaemonCheck =5,              // compatible with old daemon
   ClientServerDoubleCheck_C =6,
   ClientCheckOnly_C = 7,
   ServerCheckOnly_C = MulPointTcpCheck_C,
};


enum MetaType
{
	ANY_TYPE = 0,
//	DAEMON_TYPE = 2,		  	// 从旧daemon倒过来的数据
//  1~127: reserved by s2s

	// 以下可以用来区分data的编解码协议
	TEXTPLAIN = 128, 
	S2SDECODER = 129,
	YYPROTOCOL = 130,
	TEXTJSON = 131,
	MUSIC_PROC = 4096
};

struct SubFilter // and relation
{
  	std::string interestedName; // 模糊匹配"prefix"+"%" , like "serviceapp%"
  	int32_t interestedGroup;   // 0表示所有机房; 
  	int32_t s2sType;            // ANY_TYPE表示关注所有type

	SubFilter():interestedGroup(0),s2sType(ANY_TYPE)
	{
	}
};

enum S2sMetaStatus
{
	S2SMETA_OK_C,
	S2SMETA_DIED_C      // be killed; 
};

struct S2sMeta
{
	int64_t serverId;    // server分配的惟一标识id;
	MetaType type;
	std::string name;    // 服务名称；
	int32_t groupId;    // 机房id
	std::string data;   // 服务信息
	int64_t timestamp;	 // set by s2s server
	S2sMetaStatus status;

	S2sMeta():
		serverId(-1),
		type(ANY_TYPE),
		groupId(0),
		timestamp(0),
		status(S2SMETA_OK_C)
	{
	}
};

enum S2sSessionStatus
{
	S2S_SESSIONOFF_C = 0,
	S2S_SESSIONON_C,
	S2S_SESSIONBIND_C,
	S2S_DNSERROR_C,
	S2S_AUTHFAILURE_C,
	S2S_ERROR_C
};


class IMetaServer
{
public:
	virtual ~IMetaServer(){}

	/*
	   初始化
	   参数:
	   myName:  服务名字，初始上线需要先申请；
	   s2sKey:    申请服务名字时，生成的一个字符串，类似于ticket; 
	   myType:   注册时data字段的编解码协议；
	   返回: 成功返回一个fd, 使用者需要监听这个fd的读事件，以获取MetaServer的状态更新和订阅更新； 失败返回-1; 
	*/
	virtual int initialize(const std::string& myName,const std::string& s2sKey, MetaType myType) = 0;

	/*
	  向服务端订阅:
	  参数:
	  filters: 订阅条件, SubFilter间是or的关系，SubFilter里面的成员是and的关系；
	  返回: 成功0, 失败-1; 
	*/
	virtual int subscribe(const std::vector<SubFilter>& filters) = 0;

	/*
	  阻塞订阅接口(会阻塞直到获取订阅节点信息, 预留接口, 3.0.0版本未实现)
	  参数:
	  filters: 订阅条件, SubFilter间是or的关系，SubFilter里面的成员是and的关系；
	  metas: 订阅节点信息
	  返回: 成功0, 失败-1;
	*/
	virtual int syncSubscribe(const std::vector<SubFilter>& filters, std::vector<S2sMeta> &metas) = 0;

	/*
	获取MetaServer的状态更新或者订阅更新；
	参数:
	metas:   输出参数，返回订阅更新；
	返回:  返回MetaServer的当前状态；
	 */
	virtual S2sSessionStatus pollNotify(std::vector<S2sMeta> &metas) = 0;

	/*
	向s2s服务端注册服务信息；
	参数:
	binData:   服务信息；建议可用S2sDataEncoder/S2sDataDecoder进行编解码.
	返回:  成功返回0,  失败返回-1; 
	*/
	virtual int setMine(const std::string& binData) = 0;

	/*
	向s2s服务端注释自己；
	返回:  成功返回0,  失败返回-1; 
	*/
	virtual int delMine()=0;

	/*
	  获取自己的meta, 主要获取serverId, 通过参数mine输出返回；
	  返回: 成功返回0,  失败返回-1; 
	  */
	virtual int getMine(S2sMeta & mine) = 0;

	/*
	是否拉取到所有订阅数据
	注: 1. 一旦拉取到所有数据, 该接口永远返回true;
		2. 若是多次订阅, 只要有一次订阅结果未返回, 则返回false
		3. 订阅结果为空，也返回true;
	*/
	virtual bool isPullAllSub() = 0;

	/*
	是否拉取到服务名=name的所有数据
	注: 1. 不支持模式匹配
		2. 订阅服务结果为空时，返回false
	*/
	virtual bool isSubscribePulled(const std::string& name) = 0;

	/*
	以下三个接口，只在调用initialize前调用才生效；成功返回0，失败返回-1;
	*/

	/*设置lostCheck模式,不设的话，默认为NoLostCheck_C*/
	virtual int setLostCheckType(LostCheckType checkType) = 0;

	/*设置机房id, 不设的话，默认从机器的hostinfo.ini里读， 可用于伪装groupId*/
	virtual int setGroupId(uint32_t groupId) = 0;

	/*设置连接目标，默认连中心点，一些小众机房，部有localDaemon的情况下，可以设置成连localDaemon, 参数置为false*/
	virtual int setTarget(bool isToDaemon) = 0;
};



/*namespace mongo
{
class BSONObj;
class BSONObjBuilder;
};
*/
class S2sDataDecoder
{
public:
	S2sDataDecoder(const std::string & data);
	~S2sDataDecoder();

	/*@return value-- the same for all the following methods
	* 0: ok
	* -1: error
	*/
	int readIpList(std::map<S2S::ISPType, uint32_t>& ips);  // <ispType, ipValue>
	int readTcpPort(uint32_t& port);
	int readUdpPort(uint32_t& port);
	int select(const std::string& key, bool& value);
	int select(const std::string& key, int32_t& value);
	int select(const std::string& key, int64_t& value);
	int select(const std::string& key, uint32_t& value);
	int select(const std::string& key, uint64_t& value);	
	int select(const std::string& key, std::string& value);

	int select(const std::string& key, std::map<uint32_t, uint32_t>& pairs);
	int select(const std::string& key, std::vector<int32_t>& value);
	int select(const std::string& key, std::vector<int64_t>& value);
	int select(const std::string& key, std::vector<std::string>& value);
private:
	void* bDecoder;

};

class S2sDataEncoder
{
public:
	S2sDataEncoder();
	~S2sDataEncoder();

	/*@return value--the same for all the following methods
	* 0: ok
	* -1: error
	*/
	int writeIpList(const std::map<S2S::ISPType, uint32_t>&ips);  // <ispType, ipvalue>
	int writeTcpPort(uint32_t port);
	int writeUdpPort(uint32_t port);
	int insert(const std::string& key, bool value);
	int insert(const std::string& key, int32_t value);
	int insert(const std::string& key, int64_t value);
	int insert(const std::string& key, uint32_t value);
	int insert(const std::string& key, uint64_t value);
	int insert(const std::string& key, const std::string& value);
	int insert(const std::string& key, const char* value);

	int insert(const std::string& key, const std::map<uint32_t, uint32_t>& pairs);
	int insert(const std::string& key,const std::vector<int32_t>& value);
	int insert(const std::string& key,const std::vector<int64_t>& value);
	int insert(const std::string& key,const std::vector<std::string>& value);
	const std::string endToString();
private:
	void* builder;
	std::set<std::string> existedKeys;
};




IMetaServer* newMetaServer(const std::string version = S2S_VERSION);


/*
*  please invoke before doing anything.
*/
void configConsolePort(uint32_t port);
void configRecoverTime(uint32_t interval);   // unit: second

#endif
