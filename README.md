# vk-sign

Пример проверки подписи и получения параметров запуска

```golang
	data, err := Parse("<your-url>", "<your-secret>")
```
или
```golang
	data, err := vksign.ParseWithUrlValues(<url.Values>, "<your-secret>")
```

Пример проверки подписи с [gin](https://github.com/gin-gonic/gin)
```golang
func Foo(c *gin.Context) {
	data, err := vksign.ParseWithUrlValues(c.Request.URL.Query(), "<your-secret>")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	fmt.Println(data)
}
```