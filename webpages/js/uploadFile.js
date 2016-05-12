/**
 * Created by smh on 16-1-30.
 */

//function  uploadFileCall(){
//    var file = document.getElementById("file").files[0];
//    uploadFile(file,"120.76.41.114", "hellojs", function(){});
//}

function uploadFile(file, ipAddr, topic, callback){
    alert("ok")

    var fSize = file.size;

    var shardSize = 64*1024;    //以64KB为一个分片
    var shardCount = Math.ceil(fSize / shardSize);  //总片数

    //var url = "http://120.76.41.114:4151/put?topic=hellojs";
    var url = "http://" + ipAddr + ":4151/put?topic=" + topic;
    var i = 0;
    continueRead();
    function onloadCall(resultArray){
        console.log("onload:"+i);
        var fShard = {
            "fileId":file.name,
            "fileName":file.name,
            "fileSize":fSize,
            "shardCount":shardCount,
            "shardSize": shardSize,
            "shardIndex": i
        };

        fShard['content'] = Array.apply(null, resultArray);

        //console.log(fShard);
        var jObj = JSON.stringify(fShard);

        $.post(url, jObj, function(){});

        console.log("onload complete");
        i++;
        continueRead();
    }
    function continueRead(){
        console.log("---for---");
        //计算每一片的起始与结束位置
        var start = i * shardSize,
            end = Math.min(fSize, start + shardSize),
            data = file.slice(start, end);  //slice方法用于切出文件的一部分
        if(i >= shardCount){
            //在这写文件传输完成后的处理
            console.log("i > shardCount, Close!");
            callback();
            return;
        }

        var reader = new FileReader();
        reader.readAsArrayBuffer(data);
        reader.onload = function (e) {
            onloadCall(new Uint8Array(this.result));
        };
    }

}
