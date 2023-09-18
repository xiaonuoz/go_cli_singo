package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Order 订单
type Order struct {
	FromStoreID uint `json:"from_store_id" gorm:"size:10"` // 下单店铺ID
	SellManID   uint `json:"sell_man_id" gorm:"size:10"`   // 销售人员ID，扫销售二维码则要附带此ID用以分销

	UserID     uint            `json:"userID" gorm:"size:10"`              // 下单用户ID
	Mobile     string          `json:"mobile" gorm:"size:11"`              // 联系人电话
	Nick       string          `json:"nick" gorm:"size:20"`                // 联系人名称
	OrderNo    string          `json:"orderNo" gorm:"size:32"`             // 订单号
	GoodsFee   decimal.Decimal `json:"goodsFee" gorm:"type:decimal(15,5)"` // 总商品金额
	TotalFee   decimal.Decimal `json:"totalFee" gorm:"type:decimal(15,5)"` // 总支付金额
	ExpireTime time.Time       `json:"expireTime" gorm:"type:datetime"`    // 订单过期时间
	Status     uint            `json:"status" gorm:"size:10"`              // 订单状态
	Remark     string          `json:"remark" gorm:"size:200"`             // 订单备注

	// UserCouponID uint            `json:"userCouponID" gorm:"size:10"`        // 用户使用的优惠券ID
	// UserCouponForOrderArr []UserCouponForOrder `json:"userCouponForOrderArr,omitempty" gorm:"ForeignKey:OrderID;AssociationForeignKey:ID"`// 用户使用的优惠劵信息
}
