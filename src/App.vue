<template>
  <div id="app">
    <el-tabs v-model="active_name" @tab-click="$router.push({path:'/'+active_name})">
      <el-tab-pane label="SM2" name="SM2"></el-tab-pane>
      <el-tab-pane label="SM3" name="SM3"></el-tab-pane>
      <el-tab-pane label="SM4" name="SM4"></el-tab-pane>
    </el-tabs>
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      active_name: 'SM3',
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
