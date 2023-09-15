/*
 * @Author: p_hanxichen
 * @Date: 2023-08-16 20:24:53
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/main.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package main

func main() {

	server := InitWebService()
	server.Run(":8080")
}
