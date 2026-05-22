<script setup>

import { ref } from 'vue';
import {Setting} from "@/js/util.js";

const settings = ref()

async function getSetting() {
  settings.value = await Setting("fetchSettings", "img")
}

async function sendSetting() {
  await Setting("editSettings", "img", settings.value)
}

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="image" rounded @click="getSetting">
        图像
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>总开关</p>
          <mdui-radio-group :value="String(settings?.enabled ?? false)"
                            @change="settings.enabled=$event.target.value === 'true'"
                            style="margin-top: 0">
            <mdui-radio value="true">开启</mdui-radio>
            <mdui-radio value="false">关闭</mdui-radio>
          </mdui-radio-group>
          <p>图像生成使用的API</p>
          <mdui-radio-group :value="settings?.api || 'openai'"
                            @change="settings.api=$event.target.value"
                            style="margin-top: 0">
            <mdui-radio value="openai">OpenAI</mdui-radio>
            <mdui-radio value="alibaba">Alibaba</mdui-radio>
          </mdui-radio-group>
          <mdui-divider></mdui-divider>
          <p>过期时间</p>
          <mdui-text-field variant="outlined" label="分钟"
                           :value="settings?.cache_minutes ?? '60'"
                           @change="settings.cache_minutes = Number($event.target.value)"></mdui-text-field>
          <mdui-divider></mdui-divider>
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
