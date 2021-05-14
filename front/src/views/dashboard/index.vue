<template>
  <div class="dashboard-container">
    <!-- <div class="dashboard-text">name: {{ name }}</div> -->
    <panel-group />
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import PanelGroup from './PanelGroup'

export default {
  name: 'Dashboard',
  components: {
    PanelGroup,
  },
  computed: {
    ...mapGetters([
      'name'
    ])
  },
  data() {
    return {
    }
  },
  methods: {
   
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
