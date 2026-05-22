<script setup>

import { ref } from 'vue';
import {Setting} from "@/js/util.js";

const settings = ref()
const input1 = ref('')

async function getSetting() {
  settings.value = await Setting("fetchSettings", "imgSecurity")
}

async function sendSetting() {
  await Setting("editSettings", "imgSecurity", settings.value)
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
          <p>允许的Prompt</p>
          <mdui-text-field variant="outlined" label="适配通配符" clearable
                            @input="input1 = $event.target.value" :value="input1">
            <mdui-button-icon slot="end-icon" icon="add" @click="()=>{if (input1!=='') settings.accepted_prompt_glob.push(input1);input1=''}"></mdui-button-icon>
          </mdui-text-field>
          <div class="list">
            <mdui-list>
              <mdui-list-item v-for="(item, index) in settings?.accepted_prompt_glob || []" nonclickable>
                {{ item }}
                <mdui-button-icon slot="end-icon" icon="delete" @click="()=>{if (settings.accepted_prompt_glob.length>1) settings.accepted_prompt_glob.splice(index, 1)}"></mdui-button-icon>
              </mdui-list-item>
            </mdui-list>
          </div>
          <mdui-button @click="sendSetting">确认</mdui-button>
        </mdui-card>
      </mdui-list-item>
    </mdui-collapse-item>
  </mdui-collapse>
</template>

<style scoped>

</style>
