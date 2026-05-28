<script setup>

import { ref } from 'vue';
import {sha256} from "js-sha256";
import {Setting} from "@/js/util.js";

const settings = ref({})
const input1 = ref('')
const input2 = ref('')
const ip = ref('')

async function getSetting() {
  settings.value = await Setting("fetchSettings", "security")
}

async function sendSetting() {
  await Setting("editSettings", "security", settings.value)
}

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item>
      <mdui-list-item slot="header" icon="security" rounded @click="getSetting">
        后台安全
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>后台登录密码</p>
          <mdui-text-field type="password" variant="outlined"
                           @change="settings.password_hash = sha256($event.target.value)"
                           toggle-password label="密码"></mdui-text-field>
          <p>允许登录后台的IP</p>
          <mdui-text-field variant="outlined" label="支持通配符或正则表达式(re:开头)" clearable
                            @input="input1 = $event.target.value" :value="input1">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input1!=='') settings.dashboard_allowed_ips.push(input1);input1=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?.dashboard_allowed_ips || []" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings.dashboard_allowed_ips.length>1) settings.dashboard_allowed_ips.splice(index, 1)}"></mdui-button-icon>
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-button @click="settings.dashboard_allowed_ips.push(ip)">添加本机IP</mdui-button>
          <mdui-divider></mdui-divider>
          <p>可以使用urlAPI的网站（防盗）</p>
          <mdui-text-field variant="outlined" label="支持通配符或正则表达式(re:开头)" clearable
                            @input="input2 = $event.target.value" :value="input2">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input2!=='') settings.allowed_referers.push(input2);input2=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?.allowed_referers || []" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings.allowed_referers.length>1) settings.allowed_referers.splice(index, 1)}"></mdui-button-icon>
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>
