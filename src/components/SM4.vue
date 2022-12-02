<template>
  <div style="margin:40px">
    <el-row :gutter="20">
      <el-col :span="11">
        <el-input v-model="key_data" placeholder="请输入密钥"></el-input>
        <el-input type="textarea" :rows="18" v-model="input_data" placeholder="请输入内容"></el-input>
      </el-col>
      <el-col :span="2">
        <el-button type="primary" v-on:click="hash">加密</el-button>
      </el-col>
      <el-col :span="11">
        <el-input type="textarea" :rows="20" :readonly="true" v-model="output_data"></el-input>
      </el-col>
    </el-row>
  </div>
</template>

<script>
//import ''
export default {
  name: 'SM3',
  data() {
    return {
      input_data: '',
      key_data: '',
      output_data: '',
    }
  },
  mounted() {
    let script = document.createElement('script');
    script.type = 'text/javascript';
    script.src = './static/wasm/wasm_exec.js';
    document.body.appendChild(script);
  },
  methods: {
    hash: function() {
      if (!WebAssembly.instantiateStreaming) { 
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
          const source = await (await resp).arrayBuffer()
          return await WebAssembly.instantiate(source, importObject)
        }
      }
      function loadWasm (path) {
        const go = new Go()

        return new Promise((resolve, reject) => {
          WebAssembly.instantiateStreaming(fetch(path), go.importObject)
            .then((result) => {
              go.run(result.instance)
              resolve(result.instance)
            })
            .catch((error) => {
              reject(error)
            })
        })
      }

      // Load the wasm file
      loadWasm('./static/wasm/gmsm.wasm')
        .then((wasm) => {
          this.output_data = SM4.encode(this.key_data, this.input_data)
        })
        .catch((error) => {
          console.log('ouch', error)
        })
    }
  }
}
</script>
