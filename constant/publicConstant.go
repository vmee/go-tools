package constant

import "time"

/**
public constant
*/

// status值
var IsEnable uint64 = 1  // 启用
var IsDisable uint64 = 2 // 禁用

// bool status值
var IsYes uint64 = 1 // 是
var IsNo uint64 = 2  // 否

//软删除
var IsDelNo int64 = 0  //未删除
var IsDelYes int64 = 1 //已删除

//时间格式化模版
var DateTime = "2006-01-02 15:04:05"
var DateHour = "2006-01-02 15:00:00"
var DateOnly = "2006-01-02"
var TimeOnly = "15:04:05"

var CstZone = time.FixedZone("CST", 8*3600) // 东八区
