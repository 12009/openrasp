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

package com.baidu.openrasp.hook.server.weblogic;

import com.baidu.openrasp.HookHandler;
import com.baidu.openrasp.hook.server.ServerStartupHook;
import com.baidu.openrasp.tool.Reflection;
import com.baidu.openrasp.tool.model.ApplicationModel;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

import java.io.IOException;

/**
 * @author anyang
 * @Description: weblogic启动hook点
 * @date 2018/8/27 11:50
 */
public class WeblogicStartupHook extends ServerStartupHook {
    @Override
    public boolean isClassMatched(String className) {
        return "weblogic/t3/srvr/T3Srvr".equals(className);
    }

    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(WeblogicStartupHook.class, "handleWeblogicStartup", "$0", Object.class);
        insertBefore(ctClass, "startup", "()V", src);
    }

    public static void handleWeblogicStartup(Object sevrer) {
        try {
            Class clazz = sevrer.getClass().getClassLoader().loadClass("weblogic.version");
            String version = (String)Reflection.invokeStaticMethod(clazz.getName(),"getVersions",new Class[]{});
            ApplicationModel.init("weblogic",version);
        } catch (Exception e) {
            HookHandler.LOGGER.warn("handle weblogic startup failed", e);
        }
    }
}
