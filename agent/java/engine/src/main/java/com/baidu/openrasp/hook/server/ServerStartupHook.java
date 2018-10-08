/*
 * Copyright 2017-2018 Baidu Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.baidu.openrasp.hook.server;

import com.baidu.openrasp.cloud.Register;
import com.baidu.openrasp.cloud.Utils.CloudUtils;
import com.baidu.openrasp.hook.AbstractClassHook;

/**
 * Created by tyy on 18-8-10.
 *
 * 用于 hook 服务器的启动函数，用于记录服务器的基本信息，同时可以用作基线检测
 */
public abstract class ServerStartupHook extends AbstractClassHook {

    @Override
    public String getType() {
        return "server_start";
    }

    protected static void sendRegister() throws Exception {
        if (CloudUtils.checkCloudControlEnter()){
            Register.register();
        }
    }
}
