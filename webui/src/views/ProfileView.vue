<script>
import UserProfile from '@/components/UserProfile.vue'

export default {
  name: 'ProfileView',
  components: {
    UserProfile
  }
}
</script>

<template>
  <div class="profile-view">
    <UserProfile />
  </div>
</template>

<style scoped>
.profile-view {
  padding: 20px;
}
</style>