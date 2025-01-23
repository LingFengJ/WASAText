<script>
import NewGroup from '@/components/NewGroup.vue'

export default {
  name: 'NewGroupView',
  components: {
    NewGroup
  }
}
</script>

<template>
  <div class="new-group-view">
    <NewGroup />
  </div>
</template>

<style scoped>
.new-group-view {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}
</style>