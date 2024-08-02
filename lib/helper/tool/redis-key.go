package tool

func BillRedisKey(orderMainID string) string {
	return "bill_" + orderMainID
}
