package command

const (
	//AABBCCCC AA表示设备类型 BB表示该设备的消息类型 CCCC表示具体消息内容
	//->phone
	//心跳反馈
	PHONE_IN_HEARTBEAT uint32 = 0x01010000
	//请求开门
	PHONE_IN_OPENGARAGE_REQ uint32 = 0x01010001
	//询问添加设备
	//向车库主人发送开门提醒
	//<-phone
	//心跳
	PHONE_OUT_HEARTBEAT uint32 = 0x01020000
	//同意开门
	PHONE_OUT_OPEN_OK uint32 = 0x01020001
	//取消开门
	PHONE_OUT_OPEN_CANCEL uint32 = 0x01020002

	//->garage
	//心跳反馈
	GARAGE_IN_HEARTBEAT uint32 = 0x02010000
	//开门命令
	GARAGE_IN_OPEN uint32 = 0x02010001
	//关门命令
	GARAGE_IN_CLOSE uint32 = 0x02010002
	//查询状态
	GARAGE_IN_STATE_QUERY uint32 = 0x02010003
	//<-garage
	//心跳（包含状态信息）
	GARAGE_OUT_HEARTBEAT uint32 = 0x02020000
	//汇报状态
	GARAGE_OUT_STATE uint32 = 0x02020001
	//开门状态超时提醒

	//->camera
	//心跳反馈
	CAMERA_IN_HEARTBEAT uint32 = 0x03010000
	//截取图片命令
	CAMERA_IN_PICKUP_IMG uint32 = 0x03010001
	//拍摄视频
	CAMERA_IN_RECORD_VIDEO uint32 = 0x03010002
	//转动角度
	CAMERA_IN_TURN_AROUND uint32 = 0x03010003

	//<-camera
	//心跳
	CAMERA_OUT_HEARTBEAT uint32 = 0x03020000
	//上传图片
	CAMERA_OUT_IMG uint32 = 0x03020001
	//上传视频
	CAMERA_OUT_VIDEO uint32 = 0x03020002

	//->imgprocessor
	//接收图片
	IMGPROCESSOR_IN_IMG uint32 = 0x04010001
	//<-imgpeocessor
	//返回车牌号
	IMGPROCESSOR_OUT_CAR_NO uint32 = 0x04020001

	//->finduser
	//查找用户（手机）
	FINDUSER_IN_QUERY uint32 = 0x05010001

	//<-finduser
	//返回用户（手机）列表
	FINDUSER_OUT_USERS uint32 = 0x05020001

	//->procengine
	//启动流程
	ENGINE_IN_NEW_INSTANCE uint32 = 0x06010001

	ENGINE_IN_RESULT uint32 = 0x06010002

	EVENT_STRING   uint32 = 0x10010001
	COMMAND_STRING uint32 = 0x10010002

	DEVSIM_IN_MESSAGE uint32 = 0x11010001
)
