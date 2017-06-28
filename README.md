# paycode
演示支付宝和微信支付的付款码生成原理


### HowToRun

    go run paycode.go `go run key.go mysecret`

### Design

利用two factor auth+uid生成与时间有关的加密18位数字组成的付款码。

### References

https://tools.ietf.org/html/rfc6238
https://en.wikipedia.org/wiki/Time-based_One-time_Password_Algorithm
