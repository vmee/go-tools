// apple 苹果应用内支付SDK（目前仅支持一个校验逻辑--服务端模式）。
// 整体流程简述：
//  1. app从服务端获取待支付订单ID（这个订单ID是自己服务端产生的订单信息）和待付费productId（如果在苹果上配置了多个productId，则需要从服务端拉取自己的商品和苹果商品关联信息，这里的是指在苹果的productId）；
//  2. app根据productId，发起应用内支付；
//  3. app得到支付成功结果，并将支付结果的receipt信息以及第一步的待支付订单ID发回自己的服务端，服务端调用SDK的Verify，得到校验结果（表示该支付信息是否是苹果处理的）；
//  4. 如果校验成功则表示苹果支付完成，根据第三步传回的订单ID来处理自己的后续业务逻辑；
//
// 综上：当使用苹果的应用内支付时，其实是由自己的客户端APP来发起的“支付结果通知”请求来推动业务支付流程的数据流转。
// 参考文档：https://help.apple.com/app-store-connect/#/devb57be10e7
// **注：实际测试时发现校验接口的返回信息和文档介绍的返回信息仅状态码字段可以对应，其他内容基本不匹配，返回信息可以参考client_test.go**
package apple
