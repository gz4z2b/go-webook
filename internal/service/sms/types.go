/*
 * @Author: p_hanxichen
 * @Date: 2023-09-07 18:23:50
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/service/sms/types.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package sms

import "context"

type Service interface {
	Send(ctx context.Context, numbers []string, tpl string, args []string) error
}
