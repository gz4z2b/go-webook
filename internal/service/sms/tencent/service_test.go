/*
 * @Author: p_hanxichen
 * @Date: 2023-09-07 19:16:08
 * @LastEditors: p_hanxichen
 * @FilePath: /go/src/webook/internal/service/sms/tencent/service_test.go
 * @Description:
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */
package tencent

import (
	"fmt"
	"os"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func TestSend(t *testing.T) {
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)
	client, _ := cvm.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())

	request := cvm.NewDescribeInstancesRequest()
	response, err := client.DescribeInstances(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", response.ToJsonString())
}
