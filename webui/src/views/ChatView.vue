<script>
import ChatComponent from '@/components/ChatView.vue'

export default {
  name: 'ChatView',
  components: {
    ChatComponent
  },
  computed: {
    conversationId() {
      return this.$route.params.id
    }
  }
}
</script>

<template>
  <div class="chat-view">
    <ChatComponent :id="conversationId" />
  </div>
</template>