<script setup>

import { ref } from 'vue';
import {Setting} from "@/js/util.js";

const settings = ref()

async function getSetting() {
  settings.value = await Setting("fetchSettings", "web")
}

async function sendSetting() {
  await Setting("editSettings", "web", settings.value)
}

function find(list, status, value, operation) {
  let index
  for (let i = 0; i < list.length; i++) {
    if (list[i] === value) {
      index = i
      if (operation === "find") {
        return true
      }
    }
  }
  if (operation === "find") {
    return false
  }
  if (operation === "edit" && status === true && index === undefined) {
    list.push(value)
  }
  if (operation === "edit" && status === false) {
    list.splice(index, 1)
  }
}

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="web" rounded @click="getSetting">
        网页
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>网页缩略图开关</p>
          <mdui-radio-group :value="String(settings?.enabled ?? false)"
                            @change="settings.enabled=$event.target.value === 'true'"
                            style="margin-top: 0">
            <mdui-radio value="true">开启</mdui-radio>
            <mdui-radio value="false">关闭</mdui-radio>
          </mdui-radio-group>
          <p>允许生成缩略图的网站</p>
          <div class="mdui-checkbox-group">
            <mdui-checkbox :checked="find(settings?.allowed_hosts || [], false, 'github.com', 'find')"
                           @change="find(settings?.allowed_hosts || [], $event.target.checked, 'github.com', 'edit')">
              Github（需要网络支持）</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?.allowed_hosts || [], false, 'gitee.com', 'find')"
                           @change="find(settings?.allowed_hosts || [], $event.target.checked, 'gitee.com', 'edit')">
              Gitee</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?.allowed_hosts || [], false, 'www.youtube.com', 'find')"
                           @change="find(settings?.allowed_hosts || [], $event.target.checked, 'www.youtube.com', 'edit')">
              Youtube（需要网络支持）</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?.allowed_hosts || [], false, 'www.bilibili.com', 'find')"
                           @change="find(settings?.allowed_hosts || [], $event.target.checked, 'www.bilibili.com', 'edit')">
              B站</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?.allowed_hosts || [], false, 'arxiv.org', 'find')"
                           @change="find(settings?.allowed_hosts || [], $event.target.checked, 'arxiv.org', 'edit')">
              Arxiv</mdui-checkbox>
            <mdui-checkbox :checked="find(settings?.allowed_hosts || [], false, 'www.ithome.com', 'find')"
                           @change="find(settings?.allowed_hosts || [], $event.target.checked, 'www.ithome.com', 'edit')">
              IT之家</mdui-checkbox>
          </div>
          <p>过期时间</p>
          <mdui-text-field variant="outlined" label="分钟"
                           :value="settings?.cache_minutes ?? '60'"
                           @change="settings.cache_minutes = Number($event.target.value)"></mdui-text-field>
          <p>生成失败时返回的图片</p>
          <mdui-text-field variant="outlined" label="URL"
                           :value="settings?.fallback_image_url || ''"
                           @change="settings.fallback_image_url = $event.target.value"></mdui-text-field>
          <mdui-button @click="sendSetting">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>
</style>
