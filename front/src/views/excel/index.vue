<template>
  <div id="excel">
    <!-- <img alt="Vue logo" src="../../assets/logo.png"> -->
    <LuckySheet msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
import LuckySheet from '../LuckySheet/LuckySheet.vue'

export default {
  name: 'excel',
  components: {
    LuckySheet
  }
}
</script>

<style>
#excel {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: calc(100vh - 50px);
  /* margin-top: 60px; */
}
</style>
