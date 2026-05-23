<script setup>

import { ref } from 'vue';
import {Setting} from "@/js/util.js";

const settings = ref()

async function getSetting() {
  settings.value = await Setting("fetchSettings", "contxt")
}

async function sendSetting() {
  await Setting("editSettings", "contxt", settings.value)
}

</script>

<template>
  <mdui-collapse>
    <mdui-collapse-item rounded>
      <mdui-list-item slot="header" icon="texture" rounded @click="getSetting">
        提示词相关
        <mdui-icon slot="end-icon" name="keyboard_arrow_down"></mdui-icon>
      </mdui-list-item>
      <mdui-list-item nonclickable>
        <mdui-card variant="outlined">
          <p>这里设置提示词语境及具体的提示词</p>
          <mdui-text-field variant="outlined" label="生成的语境"
                            :value="settings?.generation_context || ''"
                            @change="settings.generation_context = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="总结的语境"
                            :value="settings?.summary_context || ''"
                            @change="settings.summary_context = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="生成笑话的提示词 (laugh)"
                            :value="settings?.templates?.laugh || ''"
                            @change="settings.templates.laugh = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="生成诗句的提示词 (poem)"
                            :value="settings?.templates?.poem || ''"
                            @change="settings.templates.poem = $event.target.value"></mdui-text-field>
          <mdui-text-field variant="outlined" label="生成鸡汤的提示词 (sentence)"
                            :value="settings?.templates?.sentence || ''"
                            @change="settings.templates.sentence = $event.target.value"></mdui-text-field>
          <mdui-button @click="sendSetting()">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>
