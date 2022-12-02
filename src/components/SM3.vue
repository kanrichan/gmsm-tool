<template>
  <div style="margin:40px">
    <el-row :gutter="20">
      <el-col :span="11">
        <el-input type="textarea" :rows="20" v-model="input_data" placeholder="请输入内容"></el-input>
      </el-col>
      <el-col :span="2">
        <el-button type="primary" v-on:click="hash">hash</el-button>
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
  name: "SM3",
  data() {
    return {
      input_data: "",
      output_data: ""
    };
  },
  mounted() {
    let script = document.createElement("script");
    script.type = "text/javascript";
    script.src = "./static/wasm/wasm_exec.js";
    document.body.appendChild(script);
  },
  methods: {
    hash: function() {
      if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
          const source = await (await resp).arrayBuffer();
          return await WebAssembly.instantiate(source, importObject);
        };
      }
      function b2h(bytes) {
        return bytes.split("").map(function (byte){
          if (byte < 10) {
            return '0'+(byte & 0xFF).toString(16)
          } else {
            return '' + (byte & 0xFF).toString(16)
          }
        }).join('')
      }

      function loadWasm(path) {
        const go = new Go();

        return new Promise((resolve, reject) => {
          WebAssembly.instantiateStreaming(fetch(path), go.importObject)
            .then(result => {
              go.run(result.instance);
              resolve(result.instance);
            })
            .catch(error => {
              reject(error);
            });
        });
      }

      // Load the wasm file
      loadWasm("./static/wasm/gmsm.wasm")
        .then(wasm => {
          let input = new Uint8Array(Buffer.from(this.input_data));
          let sm3 = new Sm3();
          let output = sm3.write(input).sum()
          this.output_data = Array.prototype.map
            .call(output, (x) => ('00' + x.toString(16)).slice(-2))
            .join('');
        })
        .catch(error => {
          console.log("ouch", error);
        });
    }
  }
};
</script>
