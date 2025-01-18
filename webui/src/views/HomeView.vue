<script>
import ConversationList from '@/components/ConversationList.vue'

export default {
  name: 'HomeView',
  components: {
    ConversationList
  },
  data() {
    return {
      username: sessionStorage.getItem('username')
    }
  },
  methods: {
    startNewChat() {
      this.$router.push('/new-conversation')
    },
    startNewGroup() {
      this.$router.push('/new-group')
    }
  }
}
</script>

<template>
  <div class="home-view">
    <div class="d-flex justify-content-between align-items-center p-4 border-bottom">
      <h1 class="h3 mb-0">Welcome, {{ username }}</h1>
      <div class="btn-group">
        <button class="btn btn-primary me-2" @click="startNewChat">
          <i class="bi bi-chat-text"></i> New Chat
        </button>
        <button class="btn btn-primary" @click="startNewGroup">
          <i class="bi bi-people"></i> New Group
        </button>
      </div>
    </div>

    <div class="p-4">
      <h2 class="h4 mb-4">Recent Conversations</h2>
      <ConversationList />
    </div>
  </div>
</template>

<style scoped>
.home-view {
  height: 100%;
  background-color: #f8f9fa;
}

.btn-group .btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.bi {
  font-size: 1.2rem;
}
</style>