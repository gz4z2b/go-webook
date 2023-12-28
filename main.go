/*
 * @Author: p_hanxichen
 * @Date: 2023-08-16 20:24:53
 * @LastEditors: p_hanxichen
 * @FilePath: /webook/main.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package main

import "github.com/gz4z2b/go-webook/ioc"

func main() {
	server := ioc.InitDownCacheWebService()
	server.Run(":8080")
}
