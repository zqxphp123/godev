// license that can be found in the LICENSE file.

// Code generated by "codegen -type=int"; DO NOT EDIT.

package code

// init register error codes defines in this source code to `imooc/mydev/pkg/errors`
func init() {
	register(ErrConnectDB, 500, "Init db error")
	register(ErrConnectGRPC, 500, "Connect to grpc error")
	register(ErrGoodsNotFound, 404, "Goods not found")
	register(ErrCategoryNotFound, 404, "Category not found")
	register(ErrEsUnmarshal, 500, "Es unmarshal error")
	register(ErrInventoryNotFound, 404, "Inventory not found")
	register(ErrInvSellDetailNotFound, 400, "Inventory sell detail not found")
	register(ErrInvNotEnough, 400, "Inventory not enough")
	register(ErrShopCartItemNotFound, 404, "ShopCart item not found")
	register(ErrSubmitOrder, 400, "Submit order error")
	register(ErrNoGoodsSelect, 404, "No Goods selected")
	register(ErrUserNotFound, 404, "User not found")
	register(ErrUserAlreadyExists, 400, "User already exists")
	register(ErrUserPasswordIncorrect, 400, "User password incorrect")
	register(ErrSmsSend, 400, "Send sms error")
	register(ErrCodeNotExist, 400, "Sms code incorrect or expired")
	register(ErrCodeInCorrect, 400, "Sms code incorrect")
}
